// SPDX-License-Identifier: GPL-3.0-or-later
//
//go:generate goversioninfo -arm -64 .goversion.json
package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	_ "github.com/josephspurrier/goversioninfo"

	"github.com/alecthomas/kong"

	ver "github.com/cuhsat/fox/v4/internal"

	"github.com/cuhsat/fox/v4/internal/cmd"
	"github.com/cuhsat/fox/v4/internal/cmd/dump"
	"github.com/cuhsat/fox/v4/internal/cmd/hash"
	"github.com/cuhsat/fox/v4/internal/cmd/help"
	"github.com/cuhsat/fox/v4/internal/cmd/hunt"
	"github.com/cuhsat/fox/v4/internal/cmd/info"
	"github.com/cuhsat/fox/v4/internal/cmd/show"
	"github.com/cuhsat/fox/v4/internal/cmd/str"
	"github.com/cuhsat/fox/v4/internal/pkg/text"
)

var Usage = strings.TrimSpace(`
The Forensic Examiners Swiss Army Knife (%s)

Usage:
  fox [COMMAND] [FLAGS...] <PATHS...>

Commands:
   s, str                  Show file contained strings
   i, info                 Show file infos and entropy
   h, hash                 Show file hashes and checksums
   d, dump                 Dump Active Directory secrets
   e, hunt                 Hunt critical system events

File flags:
  -i, --in=FILE            Read paths from file
  -o, --out=FILE           Write output to file (receipted)

Limit flags:
  -h, --head               Limit head of file by...
  -t, --tail               Limit tail of file by...
  -c, --bytes=NUMBER       Number of bytes
  -l, --lines=NUMBER       Number of lines

Unique flags:
  -u, --uniq               Filter using unique hash
  -D, --dist=LENGTH        Filter using Levenshtein distance

Filter flags:
  -e, --regexp=PATTERN     Filter output by pattern
  -C, --context=LINES      Lines surrounding a match
  -B, --before=LINES       Lines leading before a match
  -A, --after=LINES        Lines trailing after a match

Special flags:
  -p, --password=TEXT      Use archive password (7Z, RAR, ZIP)
  -P, --parallel=CORES     Use parallel processing
  -T, --force-text         Force output as text
  -X, --force-hex          Force output as hex

Disable flags:
  -r, --raw                Don't process files at all
  -q, --quiet              Don't print anything
  -N, --no-pretty          Don't prettify the output
      --no-strict          Don't stop on parser errors
      --no-deflate         Don't deflate automatically
      --no-extract         Don't extract automatically
      --no-convert         Don't convert automatically
      --no-receipt         Don't write the receipt

Standard flags:
  -v, --verbose[=LEVEL]    Print more details (v/vv/vvv)
  -d, --dry-run            Print only the found files
      --version            Print the version number
      --help               Print this help message

Positional arguments:
  Globbing paths to open or '-' to read from STDIN

Example: Find occurrences in event logs
  $ fox -eWinlogon ./**/*.evtx

Example: List only high entropy files
  $ fox info -n6.0 ./**/*

Example: Hunt down critical events
  $ fox hunt -u *.dd

For more information please visit: https://foxhunt.dev
Use "fox help <COMMAND>" to see help on a sub command.
`)

type fox struct {
	Str  str.Str   `cmd:"" aliases:"s"`
	Hash hash.Hash `cmd:"" aliases:"h"`
	Info info.Info `cmd:"" aliases:"i"`
	Dump dump.Dump `cmd:"" aliases:"d"`
	Hunt hunt.Hunt `cmd:"" aliases:"e"`
	Show show.Show `cmd:"" default:"withargs"`

	// hidden commands
	Help help.Help `cmd:"" hidden:""`

	// support flags
	Version bool

	// global flags
	cmd.Globals
}

func main() {
	defer trace()

	log.SetFlags(0)
	log.SetPrefix("fox: ")

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
		_ = text.Usage(fmt.Sprintf(Usage, ver.Version))
	case cli.Version:
		fmt.Printf("fox %s\n", ver.Version)
	default:
		if cli.Verbose > 0 {
			defer timer(time.Now())
		}

		// redirect output
		if len(cli.File) > 0 {
			text.Setup(os.OpenFile(cli.File, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600))
		} else if cli.Quiet {
			text.Setup(os.Open(os.DevNull))
			log.SetOutput(io.Discard)
		}

		defer text.Close(cli.File, !cli.NoReceipt)

		ctx.FatalIfErrorf(ctx.Run(&cli.Globals))
	}
}

func timer(t time.Time) {
	log.Printf("time %v\n", time.Since(t))
}

func trace() {
	if err := recover(); err != nil {
		log.Printf("%+v\n\n%s\n", err, debug.Stack())
	}
}
