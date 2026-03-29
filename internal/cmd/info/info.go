package info

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/alecthomas/kong"

	cli "go.foxforensics.dev/fox/v4/internal/cmd"

	"go.foxforensics.dev/fox/v4/internal/pkg/file"
	"go.foxforensics.dev/fox/v4/internal/pkg/hash"
	"go.foxforensics.dev/fox/v4/internal/pkg/text"
)

// Threshold for high entropy files
const Threshold = 7.2

// NoOffset is used for analysis
const NoOffset = -1

var Usage = strings.TrimSpace(`
fox info [FLAGS...] <PATHS...>

Flags:
  -s, --sort               Sort files by path (slower)
  -j, --json               Show infos as JSON objects
  -J, --jsonl              Show infos as JSON lines

Block flags:
  -b, --block=SIZE         Block size for analysis

Filter flags:
  -n, --min=VALUE          Minimum entropy value (default: 0.0)
  -x, --max=VALUE          Maximal entropy value (default: 8.0)

Examples:
  $ fox info -n6.0 ./**/*

Remarks:
  If FOX_API_KEY is set, then file hashes will be checked via VirusTotal.
`)

type FileInfo struct {
	File    string       `json:"file,omitempty"`
	Bytes   int64        `json:"bytes,omitempty"`
	Lines   int64        `json:"lines,omitempty"`
	Offset  int64        `json:"offset,omitempty"`
	Entropy float64      `json:"entropy,omitempty"`
	Report  *file.Report `json:"report,omitempty"`
}

func (fi *FileInfo) String() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("%7dl ", fi.Lines))
	sb.WriteString(fmt.Sprintf("%6s ", text.Humanize(fi.Bytes)))
	sb.WriteString(fmt.Sprintf("%.1fe ", fi.Entropy))

	if fi.Offset != NoOffset {
		sb.WriteString(fmt.Sprintf("%.08x:", fi.Offset))
	}

	sb.WriteString(fi.File)

	if fi.Report != nil && len(fi.Report.String()) > 0 {
		return fmt.Sprintf("%s [%s]", sb.String(), text.AsWarn(text.AsBold(fi.Report)))
	} else if fi.Entropy > Threshold {
		return text.AsWarn(sb.String())
	} else if fi.Bytes == 0 {
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
	Sort  bool `short:"s"`
	Json  bool `short:"j" xor:"json,jsonl"`
	Jsonl bool `short:"J" xor:"json,jsonl"`

	// block
	Block string `short:"b"`

	// filter
	Min float64 `short:"n" default:"0.0"`
	Max float64 `short:"x" default:"8.0"`

	// paths
	Paths []string `arg:"" optional:""`

	// hidden
	Key  string `hidden:"" long:"api-key"`
	Jack string `hidden:"" xor:"jack,john"`
	John string `hidden:"" xor:"jack,john"`

	// internal
	block int64 `kong:"-"`
}

func (cmd *Info) Validate() error {
	if cmd.Min > cmd.Max {
		log.Fatalln("invalid range")
	}

	return nil
}

func (cmd *Info) AfterApply(_ *kong.Kong, _ kong.Vars) error {
	switch {
	case len(cmd.Jack) > 0:
		cmd.Key = file.ReserveKey(1, cmd.Jack)

	case len(cmd.John) > 0:
		cmd.Key = file.ReserveKey(2, cmd.John)
	}

	if len(cmd.Block) > 0 {
		cmd.block = text.Mechanize(cmd.Block)
	}

	return nil
}

func (cmd *Info) Run(cli *cli.Globals) error {
	if len(cmd.Paths)+len(cli.Paths) == 0 {
		return text.Usage(Usage)
	}

	if cmd.Sort {
		cli.Parallel = 1 // single threaded
	}

	if !cli.NoPretty {
		text.Title(cmd.Paths...)
	}

	ch := cli.Load(cmd.Paths, true)
	defer cli.Discard()

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
				text.Match(cmd.format(fi), cli.Regexp)
			}
			h.Discard()
			continue
		}

		if len(cmd.Key) > 0 {
			fi.Report = file.GetReport(hash.Sha256(h.Bytes()), cmd.Key)
		}

		for block := range slices.Chunk(h.Bytes(), int(n)) {
			fi.Bytes = int64(len(block))
			fi.Lines = int64(bytes.Count(block, []byte{'\n'}))

			// add possibly remaining end
			if block[len(block)-1] != '\n' {
				fi.Lines++
			}

			fi.Entropy = file.Entropy(block)

			if fi.Entropy >= cmd.Min && fi.Entropy <= cmd.Max {
				text.Match(cmd.format(fi), cli.Regexp)
			}

			fi.Offset += n
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
