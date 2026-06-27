// Package receipt provides a chain of custody file based receipt
package receipt

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	"go.foxforensics.eu/fox/v4/internal/sys/version"
)

var header = strings.TrimSpace(`
┏━━━━━━┯━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓
┃ FILE │ %-71s ┃
┠──────┼─────────────────────────────────────────────────────────────────────────┨
┃ TIME │ %s ┃
┃ USER │ %s ┃
┃ HOST │ %s ┃
┃ HASH │ %s ┃
┠──────┼─────────────────────────────────────────────────────────────────────────┨
┃ INFO │ %-71s ┃
┗━━━━━━┴━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛
`)

func Generate(path string) error {
	fi, err := os.Stat(path)

	if err != nil {
		return err
	}

	buf, err := os.ReadFile(path)

	if err != nil {
		return err
	}

	hst, err := os.Hostname()

	if err != nil {
		return err
	}

	usr, err := user.Current()

	if err != nil {
		return err
	}

	abs, err := filepath.Abs(path)

	if err != nil {
		return err
	}

	return os.WriteFile(path+".cc", []byte(fmt.Sprintf(header,
		abs,
		pad(fi.ModTime().UTC().String()),
		pad(fmt.Sprintf("%s (%s)", usr.Name, usr.Username)),
		pad(fmt.Sprintf("%s (%s)", hst, macAddr())),
		fmt.Sprintf("%x SHA256", sha256.Sum256(buf)),
		fmt.Sprintf("Fox Version %s (%s-%s)", version.Number, runtime.GOOS, runtime.GOARCH),
	)), 0600)
}

func macAddr() string {
	iff, err := net.Interfaces()

	if err == nil {
		for _, i := range iff {
			if i.Flags&net.FlagUp != 0 && !bytes.Equal(i.HardwareAddr, nil) {
				return i.HardwareAddr.String()
			}
		}
	}

	return "unknown"
}

func pad(s string) string {
	if len(s) < 71 {
		s = fmt.Sprintf("%s %s", s, strings.Repeat(".", 70-len(s)))
	}

	return s
}
