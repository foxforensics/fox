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

	app "github.com/cuhsat/fox/v4/internal"
)

var header = strings.TrimSpace(`
FOX CHAIN OF CUSTODY RECEIPT %s
Time: %s
User: %s (%s)
Host: %s (%s)
Path: %s
Hash: %x SHA256
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

	return os.WriteFile(path+".cc", []byte(fmt.Sprintf(header,
		app.Version,
		time.Now().UTC(),
		usr.Name,
		usr.Username,
		hst,
		macAddr(),
		abs,
		sha256.Sum256(buf),
	)), 0600)
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
