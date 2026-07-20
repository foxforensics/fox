// Package receipt provides a chain of custody file based receipt
package receipt

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"go.foxforensics.eu/fox/v5/internal/pkg/version"
)

var header = strings.TrimSpace(`
FOX CHAIN OF CUSTODY
====================
v%s (%s-%s)

EVIDENCE
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
%x SHA-256

COMMAND
-------
%s
`)

func Generate(path string) error {
	f, err := os.Open(path)

	if err != nil {
		return err
	}

	defer func() {
		if err := f.Close(); err != nil {
			slog.Error(err.Error())
		}
	}()

	fi, err := f.Stat()

	if err != nil {
		return err
	}

	h := sha256.New()

	_, err = io.Copy(h, f)

	if err != nil {
		return err
	}

	hn, err := os.Hostname()

	if err != nil {
		return err
	}

	un, err := user.Current()

	if err != nil {
		return err
	}

	p, err := filepath.Abs(path)

	if err != nil {
		return err
	}

	cc := fmt.Sprintf("%s.cc", filepath.Clean(p))

	//nolint:gosec // G703: path is not externally tainted
	return os.WriteFile(cc, []byte(fmt.Sprintf(header,
		version.Number, runtime.GOOS, runtime.GOARCH,
		p,
		time.Now().UTC().Format(time.RFC3339Nano),
		un.Name, un.Username,
		hn,
		fi.Size(),
		h.Sum(nil),
		strings.Join(os.Args, " "),
	)), 0600)
}
