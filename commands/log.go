package commands

import (
	"bytes"
	"fmt"
	"sort"

	log "github.com/sirupsen/logrus"
)

type LogFormatter struct{}

func (f *LogFormatter) Format(entry *log.Entry) ([]byte, error) {
	if entry.Level < log.InfoLevel {
		return (new(log.TextFormatter)).Format(entry)
	}

	var keys []string
	for k := range entry.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	b := &bytes.Buffer{}
	fmt.Fprintf(b, "%-55s ", entry.Message)
	for _, k := range keys {
		v := entry.Data[k]
		fmt.Fprintf(b, " %s=%v", k, v)
	}
	b.WriteByte('\n')
	return b.Bytes(), nil
}
