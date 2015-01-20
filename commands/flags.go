package commands

import (
	"fmt"
	"strings"
)

type stringValues []string

func (s *stringValues) String() string {
	return fmt.Sprint(*s)
}

func (s *stringValues) Set(value string) error {
	for _, v := range strings.Split(value, ",") {
		*s = append(*s, v)
	}
	return nil
}

func (s *stringValues) Type() string {
	return "[]string"
}
