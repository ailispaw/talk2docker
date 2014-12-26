package api

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh/terminal"
)

var (
	terminalWidth, terminalHeight int
)

func getFd(out io.Writer) int {
	if file, ok := out.(*os.File); ok {
		return int(file.Fd())
	}
	panic("unreachable")
}

type JSONError struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func (e *JSONError) Error() string {
	return e.Message
}

type JSONProgress struct {
	Current int   `json:"current,omitempty"`
	Total   int   `json:"total,omitempty"`
	Start   int64 `json:"start,omitempty"`
}

func (p *JSONProgress) String() string {
	var (
		width       = 200
		pbBox       string
		numbersBox  string
		timeLeftBox string
	)

	if terminalWidth > 0 {
		width = terminalWidth
	}

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

func displayJSONMessagesStream(in io.Reader, out io.Writer) error {
	var (
		dec        = json.NewDecoder(in)
		ids        = map[string]int{}
		diff       = 0
		fd         = getFd(out)
		isTerminal = terminal.IsTerminal(fd)
	)

	//oldState, err := terminal.MakeRaw(fd)
	//if err == nil {
	//	defer terminal.Restore(fd, oldState)
	//}
	terminalWidth, terminalHeight, _ = terminal.GetSize(fd)

	for {
		var jm JSONMessage
		if err := dec.Decode(&jm); err != nil {
			if err == io.EOF {
				break
			}
			return err
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
