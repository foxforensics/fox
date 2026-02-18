package data

import (
	"fmt"
	"path/filepath"
	"strings"
)

const sep = ":"

func JoinPart(path, part string) string {
	return fmt.Sprintf("%s%s%s", path, sep, part)
}

func SplitPart(path string) (string, string) {
	base := filepath.Base(path)

	i := strings.LastIndex(base, sep)

	if i < 0 {
		return path, ""
	}

	return path[:(len(path)-len(base))+i], base[i+1:]
}
