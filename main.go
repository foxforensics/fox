// Visit https://foxhunt.wtf
package main

import (
	"fmt"
	"log"
	"runtime/debug"
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
	"github.com/cuhsat/fox/v4/internal/cmd/text"
)

// Short usage
var short = strings.TrimSpace(`
The Forensic Swiss Army Knife %s

Usage:
  fox [MODE] [FLAGS ...] <PATHS ...>

Modes:
  cat    prints file (default)
  hex    prints file in hex format
  info   prints file infos and entropy
  text   prints file text contents
  hash   prints file hashes and checksums
  hunt   hunt suspicious activities

Type "fox --help" for more help...
`)

// Long usage
var long = strings.TrimSpace(`
.-------.----.--.  .--.   .--. .--.--. .--.-. .--.-----.
|   ___/ .__. \  \/  /    |  |_|  |  | |  |  \|  |   _/
|   __|  |  |  >    <     |   _   |  | |  |   '  |  |
|  |   \ '--' /  /\  \    |  | |  |  '-'  |  |\  |  |
'--'    '----'--'  '--'   '--' '--'-------'--' '-'--'
The Forensic Swiss Army Knife %s

fox [MODE] [FLAGS ...] <PATHS ...>

Modes:
  cat    prints file (default)
  hex    prints file in hex format
  info   prints file infos and entropy
  text   prints file text contents
  hash   prints file hashes and checksums
  hunt   hunt suspicious activities

File limits:
  -h, --head               limit head of file by ...
  -t, --tail               limit tail of file by ...
  -n, --lines=NUMBER       number of lines
  -c, --bytes=NUMBER       number of bytes

File loader:
  -i, --input=STRING       use input in place of file content
  -p, --pass=PASSWORD      use password for decryption (only 7Z, RAR, ZIP)

File writer:
  -f, --file=FILE          write output to file name (with SHA256)

Line filter:
  -e, --regexp=PATTERN     filter for lines that match pattern
  -C, --context=NUMBER     number of lines surrounding context of match
  -B, --before=NUMBER      number of lines leading context before match
  -A, --after=NUMBER       number of lines trailing context after match

Disable:
  -r, --raw                don't process files at all
  -q, --quiet              don't print anything
      --no-file            don't print filenames
      --no-line            don't print line numbers
      --no-color           don't colorize the output
      --no-deflate         don't deflate automatically
      --no-convert         don't convert automatically

Standard:
  -d, --dry-run            prints only the found filenames
  -v, --verbose[=LEVEL]    prints more details (v/vv/vvv)
      --version            prints the version number
      --help               prints this help message

Positional arguments:
  Globbing paths to open or '-' to read from STDIN

Examples: Find occurrences in event logs
  $ fox -eWinlogon ./**/*.evtx

Examples: Show the MBR in canonical hex
  $ fox hex -mc -hc512 disk.bin

Examples: Hunt down suspicious events
  $ fox hunt -sxv ./**/*.dd

Please report bugs to <issue@foxhunt.wtf>
`)

type fox struct {
	// command modes
	Cat  cat.Cat   `cmd:"" aliases:"c,less" default:"withargs"`
	Hex  hex.Hex   `cmd:"" aliases:"x"`
	Info info.Info `cmd:"" aliases:"i,wc"`
	Text text.Text `cmd:"" aliases:"t,strings"`
	Hash hash.Hash `cmd:"" aliases:"h"`
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
		kong.Name("fox"),
		kong.DefaultEnvars("FOX"),
		kong.NoDefaultHelp(),
	)

	switch {
	case cli.Version:
		fmt.Printf("fox %s\n", app.Version)
	case (cli.Help && ctx.Command() == "cat") || ctx.Error != nil:
		fmt.Printf(long, app.Version)
	case len(ctx.Args) == 0:
		fmt.Printf(short, app.Version)
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
