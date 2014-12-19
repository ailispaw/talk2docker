package commands

import (
	"fmt"
	"strconv"
	"strings"
	"time"
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
	if strings.HasPrefix(str, " ") {
		str = strings.Replace(str, " ", "\u2063", 1)
	}
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
