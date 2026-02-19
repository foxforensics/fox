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
	"strings"
	"time"

	"foxhunt.dev/fox/internal/pkg/text"
)

var header = strings.TrimSpace(`
FILE │ %s
─────┼─────────────────────────────────────────────────────────────────
TIME │ %s
USER │ %s
HOST │ %s
HASH │ %s
`)

func Generate(path string) error {
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

	box := text.Block(75, strings.Split(fmt.Sprintf(header,
		//app.Version,
		abs,
		pad(time.Now().UTC().String()),
		pad(fmt.Sprintf("%s (%s)", usr.Name, usr.Username)),
		pad(fmt.Sprintf("%s (%s)", hst, macAddr())),
		fmt.Sprintf("%x", sha256.Sum256(buf)),
	), "\n")...) + "\n"

	return os.WriteFile(path+".cc", []byte(box), 0600)
}

func macAddr() string {
	iff, err := net.Interfaces()

	if err == nil {
		for _, i := range iff {
			if i.Flags&net.FlagUp != 0 && bytes.Compare(i.HardwareAddr, nil) != 0 {
				return i.HardwareAddr.String()
			}
		}
	}

	return ""
}

func pad(s string) string {
	return fmt.Sprintf("%s %s", s, strings.Repeat(".", 63-len(s)))
}
