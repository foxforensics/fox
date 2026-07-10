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
	"go.foxforensics.eu/fox/v5/internal/cmd"
	"go.foxforensics.eu/fox/v5/internal/sys"
	"go.foxforensics.eu/fox/v5/internal/sys/writer"
	"go.foxforensics.eu/fox/v5/library/formats"
)

// Threshold for high entropy files
const Threshold = 7.2
const Precision = 1e6

var Usage = strings.TrimSpace(`
Usage: fox info [FLAGS...] <PATHS...>

Flags:
  -j, --json               Show infos as JSON objects
  -J, --jsonl              Show infos as JSON lines

Block flags:
  -B, --block=SIZE         Block size for analysis

Filter flags:
  -N, --min=VALUE          Minimum entropy value (default: 0.0)
  -X, --max=VALUE          Maximal entropy value (default: 8.0)

Example: List only high entropy files
  $ fox info -N6.0 ./

Example: List blocks by one megabyte
  $ fox info -B1m backup.mdf

Report bugs at: foxforensics.eu/issues
`)

type FileInfo struct {
	File    string  `json:"file,omitempty"`
	Bytes   uint64  `json:"bytes"`
	Lines   uint64  `json:"lines"`
	Offset  uint64  `json:"offset"`
	Entropy float64 `json:"entropy"`
	IsBlock bool    `json:"is_block,omitempty"`
}

func (fi *FileInfo) String() string {
	var sb strings.Builder

	e := strings.Repeat("■", int(math.Round(fi.Entropy*2)))

	_, _ = fmt.Fprintf(&sb, "%7dl ", fi.Lines)
	_, _ = fmt.Fprintf(&sb, "%7s ", sys.Humanize(fi.Bytes))

	if fi.Entropy > Threshold {
		sb.WriteString(writer.AsBold(fmt.Sprintf(" %.1fe ", fi.Entropy)))
		sb.WriteString(writer.AsBold(fmt.Sprintf("[%-16s] ", e)))
	} else {
		_, _ = fmt.Fprintf(&sb, " %.1fe ", fi.Entropy)
		_, _ = fmt.Fprintf(&sb, "[%-16s] ", e)
	}

	if fi.IsBlock {
		_, _ = fmt.Fprintf(&sb, "%.08x ", fi.Offset)
	}

	sb.WriteString(fi.File)

	if fi.Bytes == 0 {
		return writer.AsGray(sb.String())
	}

	return sb.String()
}

type Info struct {
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
	cmd.Paths = append(cmd.Paths, fox.Paths...)

	if len(cmd.Paths) == 0 {
		sys.Usage(Usage)
		return nil
	}

	// turn off for calculations while loading
	v := color.NoColor

	fox.NoPretty = true

	heaps, err := fox.Init(cmd.Paths, true)

	color.NoColor = v

	if err != nil {
		return err
	}

	for h := range heaps {
		fi := &FileInfo{File: h.String(), IsBlock: cmd.block > 0}

		n := int64(h.Size)

		if fi.IsBlock {
			n = cmd.block
		}

		// because empty files will cause errors
		if h.Size == 0 {
			if cmd.Min == 0 {
				fox.Writer.Match(formats.Auto(fi, cmd.Json, cmd.Jsonl), fox.Regexp)
			}
			h.Free()
			continue
		}

		for block := range slices.Chunk(h.Bytes(), int(n)) {
			fi.Bytes = uint64(len(block))
			fi.Lines = uint64(bytes.Count(block, []byte{'\n'}))
			fi.Entropy = float64(int(entropy.Calculate(block)*Precision)) / Precision

			if fi.Entropy >= cmd.Min && fi.Entropy <= cmd.Max {
				fox.Writer.Match(formats.Auto(fi, cmd.Json, cmd.Jsonl), fox.Regexp)
			}

			fi.Offset += uint64(n)
		}

		h.Free()
	}

	return nil
}
