package commands

import (
	"fmt"
	"strings"
)

type stringArray []string

func (s *stringArray) String() string {
	return fmt.Sprint(*s)
}

func (s *stringArray) Set(value string) error {
	for _, v := range strings.Split(value, ",") {
		*s = append(*s, v)
	}
	return nil
}

func (s *stringArray) Type() string {
	return "[]string"
}
