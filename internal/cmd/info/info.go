package info

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"slices"
	"strings"
	"time"

	"github.com/alecthomas/kong"

	cli "github.com/cuhsat/fox/v4/internal/cmd"

	"github.com/cuhsat/fox/v4/internal/pkg/data/api"
	"github.com/cuhsat/fox/v4/internal/pkg/data/api/vt"
	"github.com/cuhsat/fox/v4/internal/pkg/hash"
	"github.com/cuhsat/fox/v4/internal/pkg/text"
	"github.com/cuhsat/fox/v4/internal/pkg/types"
	"github.com/cuhsat/fox/v4/internal/pkg/types/heap"
)

var Usage = strings.TrimSpace(`
Show file infos with verdict.

fox info [FLAGS...] <PATHS...>

Flags:
  -s, --sort               Sort files by path (slower)
  -b, --block=SIZE         Block size for analysis

Filter flags:
  -n, --min=VALUE          Minimum entropy value (default: 0.0)
  -x, --max=VALUE          Maximal entropy value (default: 1.0)

Examples:
  $ fox info -n0.8 ./**/*

Remarks:
  Files hashes will be checked with VirusTotal, if FOX_API_KEY env is set.
`)

type FileInfo struct {
	Path     string      `json:"path,omitempty"`
	Lines    int64       `json:"lines,omitempty"`
	Bytes    int64       `json:"bytes,omitempty"`
	Offset   int64       `json:"offset,omitempty"`
	Entropy  float64     `json:"entropy,omitempty"`
	Modified time.Time   `json:"modified,omitempty"`
	Verdict  *api.Result `json:"verdict,omitempty"`
}

func (fi *FileInfo) String() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("%.10f ", fi.Entropy))
	sb.WriteString(fmt.Sprintf("%10d ", fi.Lines))
	sb.WriteString(fmt.Sprintf("%11s ", text.Humanize(fi.Bytes)))
	sb.WriteString(fi.Modified.Format(time.RFC3339))
	sb.WriteString(fmt.Sprintf(" %08x %s", fi.Offset, fi.Path))

	if fi.Verdict != nil {
		switch fi.Verdict.Verdict {
		case api.Unknown:
			sb.WriteString(fi.Verdict.Verdict)
		case api.Unrated, api.Clean:
			sb.WriteString(fmt.Sprintf(" [%s]", text.AsGray(fi.Verdict.Verdict)))
		default:
			sb.WriteString(fmt.Sprintf(" [%s]", text.AsWarn(fi.Verdict.Verdict)))
		}
	}

	return sb.String()
}

type Info struct {
	Sort  bool    `short:"s"`
	Block string  `short:"b"`
	Min   float64 `short:"n" default:"0.0"`
	Max   float64 `short:"x" default:"1.0"`

	// paths
	Paths []string `arg:"" name:"path" optional:""`

	// hidden
	Key  string `hidden:"" long:"api-key"`
	Jack string `hidden:"" xor:"jack,john"`
	John string `hidden:"" xor:"jack,john"`

	// internal
	block uint64 `kong:"-"`
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
		cmd.Key = api.Decrypt(vt.ReserveKey1, cmd.Jack)

	case len(cmd.John) > 0:
		cmd.Key = api.Decrypt(vt.ReserveKey2, cmd.John)
	}

	if len(cmd.Block) > 0 {
		cmd.block = uint64(text.Mechanize(cmd.Block))
	} else {
		cmd.block = math.MaxUint64
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
		text.Title(fmt.Sprintf("%-12s %10s %11s %s  %17s %6s", "Entropy", "Lines", "Size", "Modified", "Offset", "File"))
	}

	ch := cli.LoadPlain(cmd.Paths)
	defer cli.Discard()

	for h := range ch {
		fi := &FileInfo{
			Path:     h.String(),
			Modified: time.UnixMilli(int64(h.Time)).UTC(),
		}

		if len(cmd.Key) > 0 {
			var err error

			fi.Verdict, err = vt.CheckFile(hash.MustSum(types.SHA256, h.Bytes()), cmd.Key)

			if err != nil {
				log.Println(err)
			}
		}

		var off int64

		for block := range slices.Chunk(h.Bytes(), int(min(cmd.block, h.Size))) {
			fi.Offset = off
			fi.Bytes = int64(len(block))
			fi.Lines = int64(bytes.Count(block, []byte{'\n'}))
			fi.Entropy = heap.Entropy(block)

			// add possibly remaining end
			if block[len(block)-1] != '\n' {
				fi.Lines++
			}

			if fi.Entropy >= cmd.Min && fi.Entropy <= cmd.Max {

				text.Write(fi.String())
			}

			off += int64(len(block))
		}

		h.Discard()
	}

	return nil
}
