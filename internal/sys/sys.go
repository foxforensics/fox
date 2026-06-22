package sys

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"go.foxforensics.eu/fox/v4/internal/sys/terminal"
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

// Stdout default output.
var Stdout = terminal.NewWriter(os.Stdout)

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

func SetOutput(wc io.WriteCloser, err error) {
	if err == nil {
		Stdout = terminal.NewWriter(wc)
	} else {
		slog.Error(err.Error())
	}
}

func JoinPart(path, part string) string {
	return fmt.Sprintf("%s%s%s", path, ":", part)
}

func SplitPart(path string) (string, string) {
	t := strings.SplitN(path, ":", 2)

	if len(t) < 2 {
		return path, ""
	}

	return t[0], t[1]
}
