package sys

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"go.foxforensics.eu/fox/v4/internal/sys/version"
)

const Version = "Version %s %s"
const Banner = `
  .-------.----.--.  .--. 
  |   ___/ .__. \  \/  /   © 2026 by Fox Forensics
  |   __|  |  |  >    <    https://foxforensics.eu
  |  |   \ '--' /  /\  \   Version %s
  '--'    '----'--'  '--'
`

func About(msg string) {
	_, _ = fmt.Println(msg)
	_, _ = fmt.Println(fmt.Sprintf(Version, version.Number, version.ID()))
}

func Usage(msg string) {
	_, _ = fmt.Println(fmt.Sprintf(Banner, version.Number))
	_, _ = fmt.Println(msg)
}

func CreateFile(path string) io.Writer {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)

	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	return f
}

func ParseList(b []byte) []string {
	v := strings.Split(strings.TrimSpace(string(b)), "\n")

	for i, s := range v {
		v[i] = strings.TrimRight(s, "\r") // Windows line breaks
	}

	return v
}

func JoinPart(path, part string) string {
	return fmt.Sprintf("%s%s%s", path, "::", part)
}

func SplitPart(path string) (string, string) {
	t := strings.SplitN(path, "::", 2)

	if len(t) < 2 {
		return path, ""
	}

	return t[0], t[1]
}
