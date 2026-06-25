package info

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"slices"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/fatih/color"
	"go.foxforensics.eu/entropy/entropy"
	"go.foxforensics.eu/fox/v4/internal/cmd"
	"go.foxforensics.eu/fox/v4/internal/pkg/adapters/formats"
	"go.foxforensics.eu/fox/v4/internal/sys"
	"go.foxforensics.eu/fox/v4/internal/sys/terminal"
)

// Threshold for high entropy files
const Threshold = 7.2

// NoOffset is used for analysis
const NoOffset = math.MaxInt64

var Usage = strings.TrimSpace(`
Usage: fox info [FLAGS...] <PATHS...>

Flags:
  -s, --sort               Sort files by path (slower)
  -j, --json               Show infos as JSON objects
  -J, --jsonl              Show infos as JSON lines

Block flags:
  -B, --block=SIZE         Block size for analysis

Filter flags:
  -N, --min=VALUE          Minimum entropy value (default: 0.0)
  -X, --max=VALUE          Maximal entropy value (default: 8.0)

Example: List only high entropy files
  $ fox info -sN6.0 ./

Example: List blocks by one megabyte
  $ fox info -B1m backup.mdf

Report bugs at: foxforensics.eu/issues
`)

type FileInfo struct {
	File    string  `json:"file,omitempty"`
	Bytes   uint64  `json:"bytes,omitempty"`
	Lines   uint64  `json:"lines,omitempty"`
	Offset  uint64  `json:"offset,omitempty"`
	Entropy float64 `json:"entropy,omitempty"`
}

func (fi *FileInfo) String() string {
	var sb strings.Builder

	e := strings.Repeat("#", int(math.Round(fi.Entropy*2)))

	sb.WriteString(fmt.Sprintf("%7dl ", fi.Lines))
	sb.WriteString(fmt.Sprintf("%7s ", sys.Humanize(fi.Bytes)))

	if fi.Entropy > Threshold {
		sb.WriteString(terminal.AsBold(fmt.Sprintf(" %.1fe ", fi.Entropy)))
		sb.WriteString(terminal.AsBold(fmt.Sprintf("[%-16s] ", e)))
	} else {
		sb.WriteString(fmt.Sprintf(" %.1fe ", fi.Entropy))
		sb.WriteString(fmt.Sprintf("[%-16s] ", e))
	}

	if fi.Offset != NoOffset {
		sb.WriteString(fmt.Sprintf("%.08x ", fi.Offset))
	}

	sb.WriteString(fi.File)

	if fi.Bytes == 0 {
		return terminal.AsGray(sb.String())
	}

	return sb.String()
}

type Info struct {
	Sort  bool `short:"s"`
	Json  bool `short:"j" xor:"json,jsonl"`
	Jsonl bool `short:"J" xor:"json,jsonl"`

	// block flags
	Block string `short:"B" xor:"block"`

	// filter flags
	Min float64 `short:"N" default:"0.0"`
	Max float64 `short:"X" default:"8.0"`

	// paths
	Paths []string `arg:"" optional:""`

	// internal
	block int64 `kong:"-"`
}

func (cmd *Info) Validate() error {
	if cmd.Min > cmd.Max {
		return errors.New("invalid range")
	}

	return nil
}

func (cmd *Info) AfterApply(_ *kong.Kong, _ kong.Vars) error {
	var ok bool

	if len(cmd.Block) > 0 {
		if cmd.block, ok = sys.Mechanize(cmd.Block); !ok {
			return errors.New("invalid block syntax")
		}
	}

	return nil
}

func (cmd *Info) Run(fox *cmd.Globals) error {
	cmd.Paths = append(cmd.Paths, fox.Input...)

	if len(cmd.Paths) == 0 {
		return sys.Usage(Usage)
	}

	if cmd.Sort {
		fox.Threads = 1 // single threaded
	}

	// turn off for calculations while loading
	v := color.NoColor

	fox.NoConvert = true
	fox.NoPretty = true

	ch, err := fox.Init(cmd.Paths, true)

	color.NoColor = v

	if err != nil {
		return err
	}

	defer fox.Discard()

	for h := range ch {
		fi := &FileInfo{File: h.String()}

		n := int64(h.Size)

		if cmd.block > 0 {
			n = cmd.block
		} else {
			fi.Offset = NoOffset
		}

		// because empty files will cause errors
		if h.Size == 0 {
			if cmd.Min == 0 {
				sys.Stdout.Match(formats.Auto(fi, cmd.Json, cmd.Jsonl), fox.Regexp)
			}
			h.Discard()
			continue
		}

		for block := range slices.Chunk(h.Bytes(), int(n)) {
			fi.Bytes = uint64(len(block))
			fi.Lines = uint64(bytes.Count(block, []byte{'\n'}))

			// add possibly remaining end
			if block[len(block)-1] != '\n' {
				fi.Lines++
			}

			fi.Entropy = entropy.Calculate(block)

			if fi.Entropy >= cmd.Min && fi.Entropy <= cmd.Max {
				sys.Stdout.Match(formats.Auto(fi, cmd.Json, cmd.Jsonl), fox.Regexp)
			}

			fi.Offset += uint64(n)
		}

		h.Discard()
	}

	return nil
}
