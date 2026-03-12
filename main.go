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

	"github.com/cuhsat/fox/v4/internal/pkg/std"
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
	"github.com/cuhsat/fox/v4/internal/cmd/mcp"
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
  (c) cat     Show file contents (default mode)
  (x) hex     Show file contents in hex format
  (t) text    Show file contained strings
  (a) hash    Show file hashes and checksums
  (s) stat    Show file stats and entropy
  (d) dump    Dump sensitive data
  (e) test    Test suspicious files
  (u) hunt    Hunt suspicious events
  (m) mcp     Init MCP server (blocks)

File flags:
  -i, --in=FILE            Read paths from file
  -o, --out=FILE           Writeln output to file (receipted)

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
      --no-warnings        Don't show any warnings

Standard flags:
  -d, --dry-run            Print only the found files
  -v, --verbose[=LEVEL]    Print more details (v/vv/vvv)
      --version            Print the version number
      --help               Print this help message

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

		std.Init(&cli.Globals)
		defer std.Close()

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
