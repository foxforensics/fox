// SPDX-License-Identifier: GPL-3.0-or-later
//
//go:generate goversioninfo -arm -64 .goversion.json
package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	_time "time"

	"github.com/alecthomas/kong"
	_ "github.com/josephspurrier/goversioninfo"
	"go.foxforensics.eu/fox/v5/internal/cmd"
	"go.foxforensics.eu/fox/v5/internal/cmd/ad"
	"go.foxforensics.eu/fox/v5/internal/cmd/cat"
	"go.foxforensics.eu/fox/v5/internal/cmd/hash"
	"go.foxforensics.eu/fox/v5/internal/cmd/help"
	"go.foxforensics.eu/fox/v5/internal/cmd/hunt"
	"go.foxforensics.eu/fox/v5/internal/cmd/info"
	"go.foxforensics.eu/fox/v5/internal/cmd/str"
	"go.foxforensics.eu/fox/v5/internal/cmd/time"
	"go.foxforensics.eu/fox/v5/internal/sys"
)

var About = strings.TrimSpace(`
© 2026 Fox Forensics
`)

var Usage = strings.TrimSpace(`
Usage: fox {ad,cat,str,info,time,hash,hunt} <...>

Example: Find occurrences in event logs
  $ fox -FWinlogon ./**/*.evtx

Example: Hunt down critical events
  $ fox hunt -u *.dd

Use 'fox help' to show further information.
`)

type Fox struct {
	Ad   ad.Ad     `cmd:"" aliases:"a"`
	Cat  cat.Cat   `cmd:"" aliases:"c" default:"withargs"`
	Hash hash.Hash `cmd:"" aliases:"h"`
	Hunt hunt.Hunt `cmd:"" aliases:"x"`
	Info info.Info `cmd:"" aliases:"i"`
	Str  str.Str   `cmd:"" aliases:"s"`
	Time time.Time `cmd:"" aliases:"t"`

	// hidden commands
	Help help.Help `cmd:"" hidden:""`

	// support
	Version bool

	// global
	cmd.Globals
}

func main() {
	defer timer(_time.Now())
	defer trace()

	log.SetFlags(0)
	log.SetPrefix("FOX: ")
	slog.SetLogLoggerLevel(slog.LevelWarn)

	fox := new(Fox)
	ctx := kong.Parse(fox,
		kong.NoDefaultHelp(),
		kong.Name("FOX"),
		kong.DefaultEnvars("FOX"),
		kong.Vars{
			"cores": strconv.Itoa(runtime.NumCPU()),
		})

	switch {
	case len(ctx.Args) == 0:
		fallthrough // show usage

	case fox.Globals.Help:
		sys.Usage(Usage)
		os.Exit(0)

	case fox.Version:
		sys.About(About)
		os.Exit(0)

	case ctx.Command() == "help":
		sys.Usage(cmd.Usage)
		os.Exit(0)

	case ctx.Error != nil:
		slog.Error(ctx.Error.Error())
		os.Exit(1)

	default:
		defer fox.Discard()
		kong.Exit(fox.Exit)

		ctx.FatalIfErrorf(ctx.Run(&fox.Globals))
	}
}

func timer(t _time.Time) {
	slog.Info(fmt.Sprintf("total time %v", _time.Since(t)))
}

func trace() {
	if err := recover(); err != nil {
		slog.Error(fmt.Sprintf("%+v\n%s", err, debug.Stack()))
		os.Exit(1)
	}
}
