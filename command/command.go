package command

import (
	"github.com/go-zoox/core-utils/regexp"
	"github.com/go-zoox/core-utils/strings"
)

var isCommandRe, _ = regexp.New("^/\\w+\\s{0,}")

func IsCommand(text string) bool {
	return isCommandRe.Match(text)
}

func ParseCommandWithArg(text string) (command string, args string, err error) {
	parts := strings.SplitN(text[1:], " ", 2)
	if len(parts) <= 1 {
		return parts[0], "", nil
	}

	return parts[0], parts[1], nil
}
