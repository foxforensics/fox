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
	t := strings.SplitN(path, sep, 1)

	if len(t) == 1 {
		return path, ""
	}

	return t[0], strings.Join(t[1:], sep)
}
