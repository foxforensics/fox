// SPDX-License-Identifier: GPL-3.0-or-later
package main

import (
	"fmt"
	"log"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/alecthomas/kong"

	"github.com/cuhsat/fox/v4/internal"
	"github.com/cuhsat/fox/v4/internal/cmd"
	"github.com/cuhsat/fox/v4/internal/cmd/cat"
	"github.com/cuhsat/fox/v4/internal/cmd/hash"
	"github.com/cuhsat/fox/v4/internal/cmd/help"
	"github.com/cuhsat/fox/v4/internal/cmd/hex"
	"github.com/cuhsat/fox/v4/internal/cmd/hunt"
	"github.com/cuhsat/fox/v4/internal/cmd/info"
	"github.com/cuhsat/fox/v4/internal/cmd/test"
	"github.com/cuhsat/fox/v4/internal/cmd/text"
)

var Usage = strings.TrimSpace(`
.-------.----.--.  .--.   .--. .--.--. .--.-. .--.-----.
|   ___/ .__. \  \/  /    |  |_|  |  | |  |  \|  |   _/
|   __|  |  |  >    <     |   _   |  | |  |   '  |  |
|  |   \ '--' /  /\  \    |  | |  |  '-'  |  |\  |  |
'--'    '----'--'  '--'   '--' '--'-------'--' '-'--'
The Forensic Examiners Swiss Army Knife (%s)

Usage:
  fox [MODE] [FLAGS...] <PATHS...>

Modes:
  cat    prints contents (default mode)
  hex    prints contents in hex format
  info   prints infos and entropy
  test   prints test results
  text   prints text contents
  hash   prints hashes and checksums
  hunt   hunt suspicious activities

File flags:
  -f, --file=FILE          read extra paths from file
  -i, --input=TEXT         read input instead of file
  -o, --output=FILE        write output to receipted file

Limit flags:
  -h, --head               limit head of file by...
  -t, --tail               limit tail of file by...
  -c, --bytes=NUMBER       number of bytes
  -l, --lines=NUMBER       number of lines

Filter flags:
  -e, --regexp=PATTERN     filter lines by pattern

Crypto flags: 
  -p, --password=TEXT      archive password (7Z, RAR, ZIP)

Profile flags:
  -P, --profile=CPUS       parallel processing profile

Disable flags:
  -r, --raw                don't process files at all
  -q, --quiet              don't print anything
      --no-file            don't print filenames
      --no-line            don't print line numbers
      --no-color           don't colorize the output
      --no-pretty          don't prettify the output
      --no-deflate         don't deflate automatically
      --no-extract         don't extract automatically
      --no-convert         don't convert automatically
      --no-receipt         don't write the receipt
      --no-warnings        don't show any warnings

Standard flags:
  -m, --more               prints pagewise (press SPACE or Q)
  -d, --dry-run            prints only the found filenames
  -v, --verbose[=LEVEL]    prints more details (v/vv/vvv)
      --version            prints the version number
      --help               prints this help message

Positional arguments:
  Globbing paths to open or '-' to also read from STDIN

Example: Find occurrences in event logs
  $ fox -eWinlogon ./**/*.evtx

Example: Show strings in binary
  $ fox text -w sample.exe

Example: Hunt down suspicious events
  $ fox hunt -sv ./**/*.E01

Use "fox help MODE" to show more help on a specific mode.
`)

type fox struct {
	// command modes
	Cat  cat.Cat   `cmd:"" aliases:"c,less,more" default:"withargs"`
	Hex  hex.Hex   `cmd:"" aliases:"x,xxd,hd"`
	Info info.Info `cmd:"" aliases:"i,wc"`
	Test test.Test `cmd:"" aliases:"t,check"`
	Text text.Text `cmd:"" aliases:"s,strings"`
	Hash hash.Hash `cmd:"" aliases:"h,sum"`
	Hunt hunt.Hunt `cmd:"" aliases:"u"`
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
			"cpus": strconv.Itoa(runtime.NumCPU()),
		})

	switch {
	case cli.Version:
		fmt.Printf("fox %s\n", app.Version)
	case cli.Globals.Help, ctx.Command() == "help":
		fallthrough
	case len(ctx.Args) == 0, ctx.Error != nil:
		fmt.Printf(Usage, app.Version)
	default:
		if cli.Verbose > 0 {
			defer timer(time.Now())
		}

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
