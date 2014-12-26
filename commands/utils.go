package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/yungsang/tablewriter"
	"gopkg.in/yaml.v2"
)

func Truncate(s string, maxlen int) string {
	if len(s) <= maxlen {
		return s
	}
	return s[:maxlen]
}

func FormatDateTime(t time.Time) string {
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
}

func FormatNonBreakingString(str string) string {
	return strings.Replace(str, " ", "\u00a0", -1)
}

func FormatBool(b bool, strTrue, strFalse string) string {
	if b {
		return strTrue
	} else {
		return strFalse
	}
}

func FormatNumber(n float64, prec int) string {
	sign := ""
	if n < 0 {
		sign = "-"
		n = -n
	}

	if prec > 0 {
		f := math.Pow10(prec)
		x := n * f
		n = math.Floor(x+0.5) / f
	}

	str := strconv.FormatFloat(n, 'f', prec, 64)
	arr := strings.Split(str, ".")

	strInt := arr[0]
	for i := len(strInt); i > 3; {
		i -= 3
		strInt = strInt[:i] + "," + strInt[i:]
	}

	if prec > 0 {
		return sign + strInt + "." + arr[1]
	} else {
		return sign + strInt
	}
}

func FormatInt(n int64) string {
	return FormatNumber(float64(n), 0)
}

func FormatFloat(n float64) string {
	return FormatNumber(n, 3)
}

func PrintInTable(out io.Writer, header []string, items [][]string, width, align int) {
	table := tablewriter.NewWriter(out)
	if !boolNoHeader {
		if header != nil {
			table.SetHeader(header)
		}
	} else {
		table.SetBorder(false)
	}
	if width != 0 {
		table.SetColWidth(width)
	}
	if align != tablewriter.ALIGN_DEFAULT {
		table.SetAlignment(align)
	}
	table.AppendBulk(items)
	table.Render()
}

func FormatPrint(out io.Writer, value interface{}) error {
	switch {
	case boolJSON:
		return PrintInJSON(out, value)
	}
	return PrintInYAML(out, value)
}

func PrintInYAML(out io.Writer, value interface{}) error {
	data, err := yaml.Marshal(value)
	_, err = out.Write(data)
	return err
}

func PrintInJSON(out io.Writer, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	indented := new(bytes.Buffer)
	err = json.Indent(indented, data, "", "  ")
	if err != nil {
		return err
	}
	indented.WriteString("\n")

	_, err = io.Copy(out, indented)
	if err != nil {
		return err
	}

	return nil
}
