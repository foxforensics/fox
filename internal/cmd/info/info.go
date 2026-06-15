package info

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"slices"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/fatih/color"
	"go.foxforensics.eu/entropy/entropy"

	cli "go.foxforensics.eu/fox/v4/internal/cmd"

	"go.foxforensics.eu/fox/v4/internal/pkg/text"
	"go.foxforensics.eu/fox/v4/internal/pkg/types/lookup"
)

// Threshold for high entropy files
const Threshold = 7.2

// NoOffset is used for analysis
const NoOffset = math.MaxInt64

var Usage = strings.TrimSpace(`
Usage: fox info [FLAGS...] <PATHS...>

Flags:
  -l, --lookup             Lookup files via VirusTotal
  -s, --sort               Sort files by path (slower)
  -j, --json               Show infos as JSON objects
  -J, --jsonl              Show infos as JSON lines

Block flags:
  -B, --block=SIZE         Block size for analysis

Filter flags:
  -N, --min=VALUE          Minimum entropy value (default: 0.0)
  -X, --max=VALUE          Maximal entropy value (default: 8.0)

Remarks:
  A VirusTotal API key is required for lookup. 

Example: List only high entropy files
  $ fox info -N6.0 ./**/*

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
	Suspect bool    `json:"suspect,omitempty"`
}

func (fi *FileInfo) String() string {
	var sb strings.Builder

	e := strings.Repeat("#", int(math.Round(fi.Entropy*2)))

	sb.WriteString(fmt.Sprintf("%7dl ", fi.Lines))
	sb.WriteString(fmt.Sprintf("%7s ", text.Humanize(fi.Bytes)))

	if fi.Entropy > Threshold {
		sb.WriteString(text.AsWarn(fmt.Sprintf(" %.1fe ", fi.Entropy)))
		sb.WriteString(text.AsWarn(fmt.Sprintf("[%-16s] ", e)))
	} else {
		sb.WriteString(fmt.Sprintf(" %.1fe ", fi.Entropy))
		sb.WriteString(fmt.Sprintf("[%-16s] ", e))
	}

	if fi.Offset != NoOffset {
		sb.WriteString(fmt.Sprintf("%.08x ", fi.Offset))
	}

	if fi.Suspect {
		sb.WriteString(text.AsWarn(fi.File))
	} else {
		sb.WriteString(fi.File)
	}

	if fi.Bytes == 0 {
		return text.AsGray(sb.String())
	}

	return sb.String()
}

func (fi *FileInfo) ToJSON() string {
	b, _ := json.MarshalIndent(fi, "", "  ")
	return string(b)
}

func (fi *FileInfo) ToJSONL() string {
	b, _ := json.Marshal(fi)
	return string(b)
}

type Info struct {
	Lookup bool `short:"l" xor:"lookup,block"`
	Sort   bool `short:"s"`
	Bars   bool `short:"b"`
	Json   bool `short:"j" xor:"json,jsonl"`
	Jsonl  bool `short:"J" xor:"json,jsonl"`

	// block flags
	Block string `short:"B" xor:"lookup,block"`

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

	if cmd.Lookup {
		log.Println("warning: data will be transmitted to a third-party service!")
	}

	return nil
}

func (cmd *Info) AfterApply(_ *kong.Kong, _ kong.Vars) error {
	if len(cmd.Block) > 0 {
		cmd.block = text.Mechanize(cmd.Block)
	}

	return nil
}

func (cmd *Info) Run(cli *cli.Globals) error {
	cmd.Paths = append(cmd.Paths, cli.Input...)

	if len(cmd.Paths) == 0 {
		return text.Usage(Usage)
	}

	if cmd.Sort {
		cli.Threads = 1 // single threaded
	}

	// turn off for calculations
	v := cli.NoPretty

	cli.NoConvert = true
	cli.NoPretty = true

	ch := cli.Load(cmd.Paths, true)
	defer cli.Discard()

	color.NoColor = v

	for h := range ch {
		fi := &FileInfo{File: h.String()}

		n := int64(h.Size)

		if cmd.block > 0 {
			n = cmd.block
		} else {
			fi.Offset = NoOffset
		}

		if cmd.Lookup {
			fi.Suspect = lookup.Lookup(h.Bytes(), cli.Verbose)
		}

		// because empty files will cause errors
		if h.Size == 0 {
			if cmd.Min == 0 {
				text.Stdout.Match(cmd.format(fi), cli.Regexp)
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
				text.Stdout.Match(cmd.format(fi), cli.Regexp)
			}

			fi.Offset += uint64(n)
		}

		h.Discard()
	}

	return nil
}

func (cmd *Info) format(fi *FileInfo) string {
	switch {
	case cmd.Jsonl:
		return text.ColorizeAs(fi.ToJSONL(), "json")
	case cmd.Json:
		return text.ColorizeAs(fi.ToJSON(), "json")
	default:
		return fi.String()
	}
}
