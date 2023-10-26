package form

import (
	"regexp"
	"strings"
)

var spaceSeperate = regexp.MustCompile(`[\s]+`)

func SplitBySpace(value string) []string {
	value = strings.TrimSpace(value)
	if len(value) == 0 {
		return nil
	}
	return spaceSeperate.Split(value, -1)
}
