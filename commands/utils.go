package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func GetStringFlag(ctx *cobra.Command, name string) string {
	flag := ctx.Flag(name)
	if flag == nil {
		return ""
	}
	return flag.Value.String()
}

func GetBoolFlag(ctx *cobra.Command, name string) bool {
	flag := ctx.Flag(name)
	if flag == nil {
		return false
	}
	return flag.Value.String() == "true"
}

func FormatDateTime(t time.Time) string {
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
}

func Truncate(s string, maxlen int) string {
	if len(s) <= maxlen {
		return s
	}
	return s[:maxlen]
}

func FormatNonBreakingString(str string) string {
	return strings.Replace(str, " ", "\u00a0", -1)
}

func FormatBool(b bool, str string) string {
	if b {
		return str
	} else {
		return ""
	}
}
