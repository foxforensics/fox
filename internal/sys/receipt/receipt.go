// Package receipt provides a chain of custody file based receipt
package receipt

import (
	"crypto/sha256"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"go.foxforensics.eu/fox/v4/internal/sys/version"
)

var header = strings.TrimSpace(`
FOX CHAIN OF CUSTODY
====================
v%s (%s-%s)

ARTIFACT
--------
%s

METADATA
--------
Acquired : %s
Examiner : %s (%s)
Host     : %s
Size     : %d

INTEGRITY
---------
SHA-256 : %x

COMMAND
-------
%s
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
		version.Number, runtime.GOOS, runtime.GOARCH,
		abs,
		time.Now().UTC().Format(time.RFC3339Nano),
		usr.Name, usr.Username,
		hst,
		len(buf),
		sha256.Sum256(buf),
		strings.Join(os.Args, " "),
	)), 0600)
}
