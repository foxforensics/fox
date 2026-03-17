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
	"github.com/cuhsat/fox/v4/internal/cmd/cat"
	"github.com/cuhsat/fox/v4/internal/cmd/check"
	"github.com/cuhsat/fox/v4/internal/cmd/dump"
	"github.com/cuhsat/fox/v4/internal/cmd/hash"
	"github.com/cuhsat/fox/v4/internal/cmd/help"
	"github.com/cuhsat/fox/v4/internal/cmd/hex"
	"github.com/cuhsat/fox/v4/internal/cmd/hunt"
	"github.com/cuhsat/fox/v4/internal/cmd/stat"
	"github.com/cuhsat/fox/v4/internal/cmd/str"
	"github.com/cuhsat/fox/v4/internal/pkg/text"
)

var Usage = strings.TrimSpace(`
The Forensic Examiners Swiss Army Knife (%s)

Usage:
  fox [COMMAND] [FLAGS...] <PATHS...>

File commands:
   c, cat                  Show file contents (default)
   x, hex                  Show file contents in hex format
   s, str                  Show file contained strings
   l, stat                 Show file stats and entropy
   h, hash                 Show file hashes and checksums

Misc commands:
   v, check                Check suspicious items online
   d, dump                 Dump Active Directory secrets
   e, hunt                 Hunt critical events

File flags:
  -i, --in=FILE            Read paths from file
  -o, --out=FILE           Write output to file (receipted)

Limit flags:
  -h, --head               Limit head of file by...
  -t, --tail               Limit tail of file by...
  -c, --bytes=NUMBER       Number of bytes
  -l, --lines=NUMBER       Number of lines

Filter flags:
  -e, --regexp=PATTERN     Filter output by pattern

Archive flags: 
  -p, --password=TEXT      Use archive password (7Z, RAR, ZIP)

Profile flags:
  -T, --threads=CORES      Use parallel threads

Disable flags:
  -r, --raw                Don't process files at all
  -q, --quiet              Don't print anything
  -y, --no-pretty          Don't prettify the output
      --no-syntax          Don't colorize the syntax
      --no-strict          Don't stop on parser errors
      --no-deflate         Don't deflate automatically
      --no-extract         Don't extract automatically
      --no-convert         Don't convert automatically
      --no-receipt         Don't write the receipt
      --no-warnings        Don't print any warnings

Standard flags:
  -d, --dry-run            Print only the found files
  -v, --verbose[=LEVEL]    Print more details (v/vv/vvv)
      --version            Print the version number
      --help               Print this help message

Positional arguments:
  Globbing paths to open or '-' to read from STDIN

Example: Find occurrences in event logs
  $ fox -eWinlogon ./**/*.evtx

Example: List only high entropy files
  $ fox stat -n0.8 ./**/*

Example: Hunt down critical events
  $ fox hunt -u *.dd

For more information please visit: https://foxhunt.dev
Use "fox help <COMMAND>" to see help on a sub command.
`)

type fox struct {
	// file commands
	Cat  cat.Cat   `cmd:"" aliases:"c,less,more" default:"withargs"`
	Hex  hex.Hex   `cmd:"" aliases:"x,xxd,hexdump"`
	Str  str.Str   `cmd:"" aliases:"s,strings"`
	Stat stat.Stat `cmd:"" aliases:"l,ls,wc,dir"`
	Hash hash.Hash `cmd:"" aliases:"h"`

	// misc commands
	Check check.Check `cmd:"" aliases:"v,vt"`
	Dump  dump.Dump   `cmd:"" aliases:"d"`
	Hunt  hunt.Hunt   `cmd:"" aliases:"e,carve"`

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
