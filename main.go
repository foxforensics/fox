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
	"go.foxforensics.eu/fox/v4/internal/sys"
)

type Fox struct {
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
	defer timer(time.Now())
	defer trace()

	log.SetFlags(0)
	log.SetPrefix("FOX: ")

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

	case fox.Globals.Help, ctx.Command() == "help":
		sys.Usage(cmd.Usage)
		os.Exit(0)

	case fox.Version:
		sys.About(cmd.About)
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

func timer(t time.Time) {
	slog.Info(fmt.Sprintf("total time %v", time.Since(t)))
}

func trace() {
	if err := recover(); err != nil {
		slog.Error(fmt.Sprintf("%+v\n%s", err, debug.Stack()))
		os.Exit(1)
	}
}
