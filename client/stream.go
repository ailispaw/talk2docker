package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

	api "github.com/yungsang/dockerclient"
)

var (
	terminalWidth, terminalHeight int
)

func init() {
	terminalHeight, terminalWidth = getTerminalSize()
}

func getTerminalSize() (int, int) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return 0, 0
	}
	arr := strings.Split(strings.TrimSpace(string(out)), " ")
	height, err := strconv.Atoi(arr[0])
	if err != nil {
		height = 0
	}
	width, err := strconv.Atoi(arr[1])
	if err != nil {
		width = 0
	}
	return height, width
}

func DoStreamRequest(client *api.DockerClient, method string, path string, body []byte, headers map[string]string) error {
	b := bytes.NewBuffer(body)
	req, err := http.NewRequest(method, client.URL.String()+path, b)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	if headers != nil {
		for header, value := range headers {
			req.Header.Add(header, value)
		}
	}
	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		if !strings.Contains(err.Error(), "connection refused") && client.TLSConfig == nil {
			return fmt.Errorf("%v. Are you trying to connect to a TLS-enabled daemon without TLS?", err)
		}
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		if len(body) == 0 {
			return fmt.Errorf("Error :%s", http.StatusText(resp.StatusCode))
		}
		return fmt.Errorf("Error: %s", bytes.TrimSpace(body))
	}

	mimetype, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err == nil && mimetype == "application/json" {
		return displayJSONMessagesStream(resp.Body, os.Stdout, uintptr(syscall.Stdin), true)
	}

	_, err = io.Copy(os.Stdout, resp.Body)
	return err
}

type JSONError struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func (e *JSONError) Error() string {
	return e.Message
}

type JSONProgress struct {
	terminalFd uintptr
	Current    int   `json:"current,omitempty"`
	Total      int   `json:"total,omitempty"`
	Start      int64 `json:"start,omitempty"`
}

func (p *JSONProgress) String() string {
	var (
		width       = terminalWidth
		pbBox       string
		numbersBox  string
		timeLeftBox string
	)

	if p.Current <= 0 && p.Total <= 0 {
		return ""
	}
	current := fmt.Sprintf("%.3f MB", float64(p.Current)/1000000)
	if p.Total <= 0 {
		return fmt.Sprintf("%8v", current)
	}
	total := fmt.Sprintf("%.3f MB", float64(p.Total)/1000000)
	percentage := int(float64(p.Current)/float64(p.Total)*100) / 2
	if width > 110 {
		// this number can't be negetive gh#7136
		numSpaces := 0
		if 50-percentage > 0 {
			numSpaces = 50 - percentage
		}
		pbBox = fmt.Sprintf("[%s>%s] ", strings.Repeat("=", percentage), strings.Repeat(" ", numSpaces))
	}
	numbersBox = fmt.Sprintf("%8v/%v", current, total)

	if p.Current > 0 && p.Start > 0 && percentage < 50 {
		fromStart := time.Now().UTC().Sub(time.Unix(int64(p.Start), 0))
		perEntry := fromStart / time.Duration(p.Current)
		left := time.Duration(p.Total-p.Current) * perEntry
		left = (left / time.Second) * time.Second

		if width > 50 {
			timeLeftBox = " " + left.String()
		}
	}
	return pbBox + numbersBox + timeLeftBox
}

type JSONMessage struct {
	Stream          string        `json:"stream,omitempty"`
	Status          string        `json:"status,omitempty"`
	Progress        *JSONProgress `json:"progressDetail,omitempty"`
	ProgressMessage string        `json:"progress,omitempty"` //deprecated
	ID              string        `json:"id,omitempty"`
	From            string        `json:"from,omitempty"`
	Time            int64         `json:"time,omitempty"`
	Error           *JSONError    `json:"errorDetail,omitempty"`
	ErrorMessage    string        `json:"error,omitempty"` //deprecated
}

func (jm *JSONMessage) Display(out io.Writer, isTerminal bool) error {
	if jm.Error != nil {
		if jm.Error.Code == 401 {
			return fmt.Errorf("Authentication is required.")
		}
		return jm.Error
	}
	var endl string
	if isTerminal && jm.Stream == "" && jm.Progress != nil {
		// <ESC>[2K = erase entire current line
		fmt.Fprintf(out, "%c[2K\r", 27)
		endl = "\r"
	} else if jm.Progress != nil { //disable progressbar in non-terminal
		return nil
	}
	if jm.Time != 0 {
		fmt.Fprintf(out, "%s ", time.Unix(jm.Time, 0).Format("2006-01-02T15:04:05.000000000Z07:00"))
	}
	if jm.ID != "" {
		fmt.Fprintf(out, "%s: ", jm.ID)
	}
	if jm.From != "" {
		fmt.Fprintf(out, "(from %s) ", jm.From)
	}
	if jm.Progress != nil {
		fmt.Fprintf(out, "%s %s%s", jm.Status, jm.Progress.String(), endl)
	} else if jm.ProgressMessage != "" { //deprecated
		fmt.Fprintf(out, "%s %s%s", jm.Status, jm.ProgressMessage, endl)
	} else if jm.Stream != "" {
		fmt.Fprintf(out, "%s%s", jm.Stream, endl)
	} else {
		fmt.Fprintf(out, "%s%s\n", jm.Status, endl)
	}
	return nil
}

func displayJSONMessagesStream(in io.Reader, out io.Writer, terminalFd uintptr, isTerminal bool) error {
	var (
		dec  = json.NewDecoder(in)
		ids  = make(map[string]int)
		diff = 0
	)
	for {
		var jm JSONMessage
		if err := dec.Decode(&jm); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		if jm.Progress != nil {
			jm.Progress.terminalFd = terminalFd
		}
		if jm.ID != "" && (jm.Progress != nil || jm.ProgressMessage != "") {
			line, ok := ids[jm.ID]
			if !ok {
				line = len(ids)
				ids[jm.ID] = line
				if isTerminal {
					fmt.Fprintf(out, "\n")
				}
				diff = 0
			} else {
				diff = len(ids) - line
			}
			if jm.ID != "" && isTerminal {
				// <ESC>[{diff}A = move cursor up diff rows
				fmt.Fprintf(out, "%c[%dA", 27, diff)
			}
		}
		err := jm.Display(out, isTerminal)
		if jm.ID != "" && isTerminal {
			// <ESC>[{diff}B = move cursor down diff rows
			fmt.Fprintf(out, "%c[%dB", 27, diff)
		}
		if err != nil {
			return err
		}
	}
	return nil
}
