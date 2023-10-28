package form

import (
	"regexp"
	"strings"
)

var spaceSeperate = regexp.MustCompile(`[\s]+`)

func SplitBySpace(value string, formatter ...func(string) string) []string {
	value = strings.TrimSpace(value)
	if len(value) == 0 {
		return nil
	}
	values := spaceSeperate.Split(value, -1)
	if len(formatter) == 0 {
		return values
	}
	for index, value := range values {
		for _, format := range formatter {
			value = format(value)
		}
		values[index] = value
	}
	return values
}

func ExplodeCombinedLogFormat(value string) string {
	return strings.Replace(value, `{combined}`, `{remote} - {user} [{when}] "{method} {uri} {proto}" {status} {size}`, 1)
}
