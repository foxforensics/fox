package sys

import (
	"errors"
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

func About(msg string) error {
	_, err1 := fmt.Println(msg)
	_, err2 := fmt.Println(fmt.Sprintf(Version, version.Number, version.ID()))
	return errors.Join(err1, err2)
}

func Usage(msg string) error {
	_, err1 := fmt.Println(fmt.Sprintf(Banner, version.Number))
	_, err2 := fmt.Println(msg)
	return errors.Join(err1, err2)
}

func CreateFile(path string) io.Writer {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)

	if err != nil {
		slog.Error(err.Error())
	}

	return f
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
