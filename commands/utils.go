package commands

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

func GetStringFlag(ctx *cobra.Command, name string) string {
	return ctx.Flag(name).Value.String()
}

func GetBoolFlag(ctx *cobra.Command, name string) bool {
	return ctx.Flag(name).Value.String() == "true"
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
