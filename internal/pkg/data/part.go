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
	t := strings.Split(path, sep)

	if len(t) == 1 {
		return path, ""
	}

	return strings.Join(t[:len(t)-2], sep), t[len(t)-1] // TODO
}
