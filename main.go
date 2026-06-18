// SPDX-License-Identifier: GPL-3.0-or-later
//
//go:generate goversioninfo -arm -64 .goversion.json
package main

import (
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/alecthomas/kong"
	_ "github.com/josephspurrier/goversioninfo"

	"go.foxforensics.eu/fox/v4/internal/cmd"
	"go.foxforensics.eu/fox/v4/internal/cmd/ad"
	"go.foxforensics.eu/fox/v4/internal/cmd/cat"
	"go.foxforensics.eu/fox/v4/internal/cmd/hash"
	"go.foxforensics.eu/fox/v4/internal/cmd/help"
	"go.foxforensics.eu/fox/v4/internal/cmd/hunt"
	"go.foxforensics.eu/fox/v4/internal/cmd/info"
	"go.foxforensics.eu/fox/v4/internal/cmd/str"
	"go.foxforensics.eu/fox/v4/internal/pkg/text"
	"go.foxforensics.eu/fox/v4/internal/pkg/version"
)

var About = strings.TrimSpace(`
© 2026 Fox Forensics
Version %s %s
`)

var Usage = strings.TrimSpace(`
Usage: fox [COMMAND] [FLAGS...] <PATHS...>

Commands:
   a, ad                   Show Active Directory infos
   c, cat                  Show file contents (default)
   s, str                  Show file contained strings
   i, info                 Show file infos and entropy
   h, hash                 Show file hashes and checksums
   x, hunt                 Hunt critical system events

File flags:
  -I, --in=FILE            Read paths from file
  -O, --out=FILE           Write output to file (receipted)

Filter flags:
  -L, --limit=NUMBER       Filter using byte or line count
  -F, --find=PATTERN       Filter using regular expression

Process flags:
  -T, --threads=CORES      Use parallel threads
  -P, --password=TEXT      Use archive password (7z, Rar, Zip)

Disable flags:
  -r, --raw                Don't process files (r/rr/rrr)
  -q, --quiet              Don't print anything
  -n, --no-pretty          Don't prettify the output
      --no-strict          Don't stop on parser errors
      --no-deflate         Don't deflate automatically
      --no-extract         Don't extract automatically
      --no-convert         Don't convert automatically
      --no-receipt         Don't write the receipt

Standard flags:
  -v, --verbose[=LEVEL]    Print more details (v/vv)
  -d, --dry-run            Print only the found files
      --version            Print the version number
      --help               Print this help message

Positional arguments:
  Globbing paths to open or '-' to read from STDIN

Example: Find occurrences in event logs
  $ fox -FWinlogon ./**/*.evtx

Example: Hunt down critical events
  $ fox hunt -u *.dd

Example: Show help on sub commands
  $ fox help info

Report bugs at: foxforensics.eu/issues
`)

type fox struct {
	Ad   ad.Ad     `cmd:"" aliases:"a"`
	Cat  cat.Cat   `cmd:"" aliases:"c" default:"withargs"`
	Hash hash.Hash `cmd:"" aliases:"h"`
	Hunt hunt.Hunt `cmd:"" aliases:"x"`
	Info info.Info `cmd:"" aliases:"i"`
	Str  str.Str   `cmd:"" aliases:"s"`

	// hidden commands
	Help help.Help `cmd:"" hidden:""`

	// support
	Version bool

	// global
	cmd.Globals
}

func main() {
	defer trace()

	log.SetFlags(0)
	log.SetPrefix("FOX ") // 🦊

	cli := new(fox)
	ctx := kong.Parse(cli,
		kong.NoDefaultHelp(),
		kong.Name("fox"),
		kong.DefaultEnvars("FOX"),
		kong.Vars{
			"cores": strconv.Itoa(runtime.NumCPU()),
		})

	switch {
	case len(ctx.Args) == 0, ctx.Error != nil:
		fallthrough // show usage

	case cli.Globals.Help, ctx.Command() == "help":
		_ = text.Usage(Usage)

	case cli.Version:
		fmt.Printf(About, version.Number, version.ID())

	default:
		defer timer(time.Now())

		// parse input
		if len(cli.In) > 0 {
			cli.Input = split(cli.In)
		}

		// redirect output
		if len(cli.Out) > 0 {
			store(cli.Out)
		} else if cli.Quiet {
			quiet()
		}

		defer text.Stdout.Close(cli.Out, !cli.NoReceipt)

		ctx.FatalIfErrorf(ctx.Run(&cli.Globals))
	}
}

func split(b []byte) []string {
	v := strings.Split(strings.TrimSpace(string(b)), "\n")

	// normalize Windows paths
	for i, s := range v {
		v[i] = strings.TrimRight(s, "\r")
	}

	return v
}

func store(f string) {
	text.SetOutput(os.OpenFile(f, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600))
}

func timer(t time.Time) {
	slog.Debug(fmt.Sprintf("time %v", time.Since(t)))
}

func quiet() {
	text.SetOutput(os.Open(os.DevNull))
	log.SetOutput(io.Discard)
}

func trace() {
	if err := recover(); err != nil {
		slog.Error(fmt.Sprintf("%+v", err), err)
		slog.Debug("--")
		slog.Debug(string(debug.Stack()))
		slog.Debug("--")
	}
}
