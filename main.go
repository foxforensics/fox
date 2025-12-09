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
)

var short = strings.TrimSpace(`
Usage: fox [COMMAND] <PATHS>

  hunt   hunt suspicious activities
  hash   prints file hash using algorithm(s)
  info   prints file info and entropy
  text   prints file strings (ASCII)
  hex    prints file in hex format
  cat    prints file (default)

Type "fox --help" for more help...
`)

var long = strings.TrimSpace(`
.-------.----.--.  .--.   .--. .--.--. .--.-. .--.-----.
|   ___/ .__. \  \/  /    |  |_|  |  | |  |  \|  |   _/
|   __|  |  |  >    <     |   _   |  | |  |   '  |  |
|  |   \ '--' /  /\  \    |  | |  |  '-'  |  |\  |  |
'--'    '----'--'  '--'   '--' '--'-------'--' '-'--'
The Forensic Swiss Army Knife %s

Usage:
  fox [COMMAND] [FLAGS] <PATHS>

Commands:
  hunt [FLAGS] <PATHS>     hunt suspicious activities
    -a, --all              show logs with all severities
    -x, --ext              show logs with all extensions (slow)
    -s, --sort             show logs sorted by timestamp (slow)
    -j, --json             show logs as JSON objects
    -J, --jsonl            show logs as JSON lines
    -D, --sqlite           save logs to SQLite3 DB

  hash [FLAGS] <PATHS>     prints file hashes and checksums
    -a, --algo=ALGO[,]     use algorithm(s) (default: SHA256)
    -F, --find=HASH[,]     show only files that match

  info [FLAGS] <PATHS>     prints file infos and entropy
    -a, --min=DECIMAL      minimum entropy value (default: 0.0)
    -b, --max=DECIMAL      maximal entropy value (default: 1.0)

  text [FLAGS] <PATHS>     prints file strings (ASCII)
    -a, --min=NUMBER       minimum string length (default: 3)
    -b, --max=NUMBER       maximal string length (default: 256)
    -w, --wtf=[=LEVEL]     show text classifications (w/ww/www)
    -1, --first            show only first text class

  hex [FLAGS] <PATHS>      prints file in hex format
    -m, --mode=[c|hd|xxd]  use compatible mode for output 

  cat [FLAGS] <PATHS>      prints file (default)

File limits:
  -h, --head               limit head of file by ...
  -t, --tail               limit tail of file by ...
  -n, --lines=NUMBER       number of lines
  -c, --bytes=NUMBER       number of bytes

File loader:
  -p, --pass=PASSWORD      password for decryption (only RAR and ZIP)

Line filter:
  -e, --regexp=PATTERN     filter for lines that match pattern
  -C, --context=NUMBER     number of lines surrounding context of match
  -B, --before=NUMBER      number of lines leading context before match
  -A, --after=NUMBER       number of lines trailing context after match

Data stream:
  -f, --file=FILE          stream data to file name
  -u, --url=SERVER         stream data to server address
  -T, --auth=TOKEN         stream data using auth token
  -E, --ecs                use ECS schema for streaming
  -H, --hec                use HEC schema for streaming

Disable:
  -r, --raw                don't process files at all
  -q, --quiet              don't print anything
      --no-file            don't print filenames
      --no-line            don't print line numbers
      --no-color           don't colorize the output
      --no-deflate         don't deflate automatically
      --no-convert         don't convert automatically

Aliases:
  -L, --logstash           alias for: -Euhttp://localhost:8080
  -S, --splunk             alias for: -Huhttp://localhost:8088/...

Standard:
  -d, --dry-run            prints only the found filenames
  -v, --verbose[=LEVEL]    prints more details (v/vv/vvv)
      --version            prints the version number
      --help               prints this help message

Positional arguments:
  Path(s) to open or '-' for STDIN

Hashes (cryptographic):
  SHA1, SHA256, SHA3, SHA3-224, SHA3-256, SHA3-384, SHA3-512,
  MD4, MD5, BLAKE3-256, BLAKE3-512, RIPEMD-160

Hashes (performance):
  FNV-1, FNV-1A, XXH64, XXH3

Hashes (similarity):
  SDHASH, SSDEEP, TLSH

Checksums:
  ADLER32, CRC32-IEEE, CRC64-ECMA, CRC64-ISO

Example: Find occurrences in event logs
  $ fox cat -eWinlogon ./**/*.evtx

Example: Show the MBR in canonical hex
  $ fox hex -mc -hc512 disk.bin

Example: Hunt down suspicious events
  $ fox hunt -sxv ./**/*.dd

Report bugs at <mail@foxhunt.wtf>
`)

type Fox struct {
	Help, Version bool
	cmd.Cli
}

func main() {
	defer trace()

	log.SetFlags(0)
	log.SetPrefix("fox: ")

	fox := new(Fox)
	ctx := kong.Parse(fox,
		kong.NoDefaultHelp(),
		kong.DefaultEnvars("FOX"),
		kong.ConfigureHelp(kong.HelpOptions{}),
	)

	switch {
	case fox.Version:
		fmt.Printf("fox %s\n", app.Version)
	case fox.Help || ctx.Error != nil:
		fmt.Printf(long, app.Version)
	case len(ctx.Args) == 0:
		fmt.Print(short)
	default:
		if fox.Cli.Verbose > 1 {
			defer timer(time.Now())
		}

		if err := ctx.Run(&fox.Cli); err != nil {
			log.Fatal(err)
		}
	}
}

func timer(t time.Time) {
	log.Printf("took %v\n", time.Since(t))
}

func trace() {
	if err := recover(); err != nil {
		log.Printf("%+v\n\n%s\n", err, debug.Stack())
	}
}
