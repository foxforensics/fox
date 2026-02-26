// SPDX-License-Identifier: GPL-3.0-or-later
//
//go:generate goversioninfo -arm -64 .goversion.json
package main

import (
	"fmt"
	"log"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/cuhsat/fox/v4/internal/cmd/mcp"
	_ "github.com/josephspurrier/goversioninfo"

	"github.com/alecthomas/kong"

	"github.com/cuhsat/fox/v4/internal"
	"github.com/cuhsat/fox/v4/internal/cmd"
	"github.com/cuhsat/fox/v4/internal/cmd/cat"
	"github.com/cuhsat/fox/v4/internal/cmd/dump"
	"github.com/cuhsat/fox/v4/internal/cmd/hash"
	"github.com/cuhsat/fox/v4/internal/cmd/help"
	"github.com/cuhsat/fox/v4/internal/cmd/hex"
	"github.com/cuhsat/fox/v4/internal/cmd/hunt"
	"github.com/cuhsat/fox/v4/internal/cmd/stat"
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
  (c) cat      shows file contents (default mode)
  (x) hex      shows file contents in hex format
  (t) text     shows file contained strings
  (a) hash     shows file hashes and checksums
  (s) stat     shows file stats and entropy
  (d) dump     dumps sensitive data
  (e) test     tests suspicious files
  (u) hunt     hunts suspicious events
  (m) mcp      loads MCP server (blocking)

File flags:
  -i, --in=FILE            reads paths from file
  -o, --out=FILE           writes output to file (receipted)

Limit flags:
  -h, --head               limits head of file by...
  -t, --tail               limits tail of file by...
  -c, --bytes=NUMBER       number of bytes
  -l, --lines=NUMBER       number of lines

Filter flags:
  -e, --regexp=PATTERN     filters output by pattern

Archive flags: 
  -p, --password=TEXT      uses archive password (7Z, RAR, ZIP)

Profile flags:
  -T, --threads=CORES      uses parallel threads

Disable flags:
  -r, --raw                don't process files at all
  -q, --quiet              don't print anything
      --no-file            don't print filenames
      --no-line            don't print line numbers
      --no-color           don't colorize the output
      --no-syntax          don't colorize the syntax
      --no-pretty          don't prettify the output
      --no-strict          don't stop on parser errors
      --no-deflate         don't deflate automatically
      --no-extract         don't extract automatically
      --no-convert         don't convert automatically
      --no-receipt         don't write the receipt
      --no-warnings        don't show any warnings

Standard flags:
  -d, --dry-run            prints only the found files
  -v, --verbose[=LEVEL]    prints more details (v/vv/vvv)
      --version            prints the version number
      --help               prints this help message

Positional arguments:
  Globbing paths to open or '-' to read from STDIN

Example: Find occurrences in event logs
  $ fox -eWinlogon ./**/*.evtx

Example: List high entropy files
  $ fox stat -n0.9 ./**/*

Example: Hunt down suspicious events
  $ fox hunt -u *.dd

For more information please visit: https://foxhunt.dev
Use "fox help <MODE>" to see help on a specific mode.
`)

type fox struct {
	// command modes
	Cat  cat.Cat   `cmd:"" aliases:"c,less,more" default:"withargs"`
	Hex  hex.Hex   `cmd:"" aliases:"x,xxd,hexdump"`
	Text text.Text `cmd:"" aliases:"t,strings"`
	Hash hash.Hash `cmd:"" aliases:"a,sum"`
	Stat stat.Stat `cmd:"" aliases:"s,ls,wc"`
	Dump dump.Dump `cmd:"" aliases:"d"`
	Test test.Test `cmd:"" aliases:"e,vt,check"`
	Hunt hunt.Hunt `cmd:"" aliases:"u"`
	Help help.Help `cmd:"" aliases:"h" hidden:""`

	// server modes
	Mcp mcp.Mcp `cmd:"" aliases:"m,serve,listen"`

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
		fmt.Println(fmt.Sprintf(Usage, res.Version))
	case cli.Version:
		fmt.Printf("fox %s\n", res.Version)
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
