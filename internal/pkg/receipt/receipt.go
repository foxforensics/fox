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
FOX CHAIN OF CUSTODY v%s [%s-%s]

ACQUIRED : %s
EXAMINER : %s (%s)
HOSTNAME : %s
EVIDENCE : %s

SHA256 %x
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

	sha := sha256.New()

	_, err = io.Copy(sha, f)

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

	cc := fmt.Sprintf("%s.cc", filepath.Base(path))

	//nolint:gosec // G703: path is not externally tainted
	return os.WriteFile(cc, []byte(fmt.Sprintf(header,
		version.Number, runtime.GOOS, runtime.GOARCH,
		time.Now().UTC().Format(time.RFC3339Nano),
		usr.Name, usr.Username,
		hst,
		strings.Join(os.Args, " "),
		sha.Sum(nil),
	)), 0600)
}
