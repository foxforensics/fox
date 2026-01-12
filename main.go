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
	"github.com/cuhsat/fox/v4/internal/cmd/hex"
	"github.com/cuhsat/fox/v4/internal/cmd/hunt"
	"github.com/cuhsat/fox/v4/internal/cmd/info"
	"github.com/cuhsat/fox/v4/internal/cmd/test"
	"github.com/cuhsat/fox/v4/internal/cmd/text"
)

var short = strings.TrimSpace(`
Usage: fox [MODE] [FLAGS...] <PATHS...>

 cat    prints file (default mode)
 hex    prints file in hex format
 info   prints file infos and entropy
 test   prints file test results
 text   prints file text contents
 hash   prints file hashes and checksums
 hunt   hunt suspicious activities

Use "fox --help" to show the full help.
`)

var long = strings.TrimSpace(`
.-------.----.--.  .--.   .--. .--.--. .--.-. .--.-----.
|   ___/ .__. \  \/  /    |  |_|  |  | |  |  \|  |   _/
|   __|  |  |  >    <     |   _   |  | |  |   '  |  |
|  |   \ '--' /  /\  \    |  | |  |  '-'  |  |\  |  |
'--'    '----'--'  '--'   '--' '--'-------'--' '-'--'
The Forensic Examiners Swiss Army Knife (%s)

Usage:
  fox [MODE] [FLAGS...] <PATHS...>

Modes:
  cat    prints file (default)
  hex    prints file in hex format
  info   prints file infos and entropy
  test   prints file test results
  text   prints file text contents
  hash   prints file hashes and checksums
  hunt   hunt suspicious activities

File limits:
  -h, --head               limit head of file by ...
  -t, --tail               limit tail of file by ...
  -n, --lines=NUMBER       number of lines
  -c, --bytes=NUMBER       number of bytes

File loader:
  -p, --pass=PASSWORD      password for decryption (7Z, RAR, ZIP)
  -f, --file=FILENAME      read extra paths from file
  -i, --input=TEXT         read input instead of file

Line output:
  -o, --output=FILE        write all output to receipted file

Line filter:
  -e, --regexp=PATTERN     filter for lines that match pattern
  -C, --context=NUMBER     number of lines surrounding context of match
  -B, --before=NUMBER      number of lines leading context before match
  -A, --after=NUMBER       number of lines trailing context after match

Profile:
  -P, --profile=CORES      parallel profile to use overall

Disable:
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

Standard:
  -d, --dry-run            prints only the found filenames
  -v, --verbose[=LEVEL]    prints more details (v/vv/vvv)
      --version            prints the version number
      --help               prints this help message

Positional arguments:
  Globbing paths to open or '-' to also read from STDIN

Example: Find occurrences in event logs
  $ fox -eWinlogon ./**/*.evtx

Example: Show MBR in canonical hex
  $ fox hex -hc512 disk.bin

Example: List high entropy files
  $ fox info -m0.9 ./**/*

Example: Test suspicious file
  $ fox test sample.exe

Example: Show strings in binary
  $ fox text -w sample.exe

Example: Hash archive contents
  $ fox hash -uTLSH files.7z

Example: Hunt down suspicious events
  $ fox hunt -sv ./**/*.E01

Use "fox MODE --help" to show more help on a specific mode.
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
	case cli.Version:
		fmt.Printf("fox %s\n", app.Version)
	case (cli.Help && ctx.Command() == "cat") || ctx.Error != nil:
		fmt.Printf(long, app.Version)
	case len(ctx.Args) == 0:
		fmt.Println(short)
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
