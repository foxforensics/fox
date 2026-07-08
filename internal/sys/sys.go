package sys

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"go.foxforensics.eu/fox/v4/internal/pkg"
	"go.foxforensics.eu/fox/v4/internal/sys/version"
	"golang.org/x/term"
)

const Version = "Version %s %s"
const Prompt = "Password: "
const Banner = `
  .-------.----.--.  .--. 
  |   ___/ .__. \  \/  /   © 2026 by Fox Forensics
  |   __|  |  |  >    <    https://foxforensics.eu
  |  |   \ '--' /  /\  \   Version %s %s
  '--'    '----'--'  '--'
`

func About(msg string) {
	_, _ = fmt.Println(msg)
	_, _ = fmt.Println(fmt.Sprintf(Version, version.Number, version.ID()))
}

func Usage(msg string) {
	_, _ = fmt.Println(fmt.Sprintf(Banner, version.Number, version.ID()))
	_, _ = fmt.Println(msg)
}

func Password() (string, error) {
	fd := int(os.Stdin.Fd())

	if !term.IsTerminal(fd) {
		return "", errors.New("stdin is not a terminal")
	}

	_, err := fmt.Fprint(os.Stderr, Prompt)

	if err != nil {
		return "", err
	}

	b, err := term.ReadPassword(fd)

	if err != nil {
		return "", err
	}

	_, err = fmt.Fprintln(os.Stderr)

	if err != nil {
		return "", err
	}

	return string(b), nil
}

func CreateFile(path string) (io.Writer, error) {
	return os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
}

func ParseList(b []byte) []string {
	v := strings.Split(strings.TrimSpace(string(b)), "\n")

	for i, s := range v {
		v[i] = strings.TrimRight(s, "\r") // Windows line breaks
	}

	return v
}

func ParseMap(b []byte) map[string]pkg.Nil {
	v := strings.Split(strings.TrimSpace(string(b)), "\n")
	m := make(map[string]pkg.Nil, len(v))

	for _, s := range v {
		m[strings.TrimRight(s, "\r")] = pkg.Nil{} // Windows line breaks
	}

	return m
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
