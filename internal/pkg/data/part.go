package data

import (
	"fmt"
	"strings"
)

const sep = ":"

func JoinPart(path, part string) string {
	return fmt.Sprintf("%s%s%s", path, sep, part)
}

func SplitPart(path string) (string, string) {
	tokens := strings.Split(path, sep) // TODO: take smb:// into account

	if len(tokens) == 1 {
		return path, ""
	}

	return tokens[0], strings.Join(tokens[1:], sep)
}
