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
	t := strings.SplitN(path, sep, 2)

	if len(t) > 1 {
		return t[0], t[1]
	}

	return path, ""
}
