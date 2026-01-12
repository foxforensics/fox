package hex

import (
	"fmt"
	"strings"

	cli "github.com/cuhsat/fox/v4/internal/cmd"

	"github.com/cuhsat/fox/v4/internal/pkg/text"
	"github.com/cuhsat/fox/v4/internal/pkg/types"
	"github.com/cuhsat/fox/v4/internal/pkg/types/buffer"
)

var Usage = strings.TrimSpace(`
Prints file in hex format.

fox hex [FLAGS...] <PATHS...>

Flags:
  -m, --mode=<hd|xxd|raw>    use compatible mode for output

Example:
  $ fox hex -hc512 disk.bin
`)

type Hex struct {
	Mode  string   `short:"m" enum:"c,hd,xxd,raw" default:"c"`
	Paths []string `arg:"" type:"path" optional:""`
}

func (cmd *Hex) Run(cli *cli.Globals) error {
	if cli.Help || len(cmd.Paths) == 0 {
		fmt.Print(Usage)
		return nil
	}

	ch := cli.Load(cmd.Paths)
	defer cli.Discard()

	var tail uint

	if cli.Tail {
		tail = cli.Bytes
	}

	for h := range ch {
		if !cli.NoFile {
			_, _ = fmt.Fprintf(cli.Stdout, "%s\n", text.Hide(text.Header(h.String())))
		}

		lastHex, wasCut := "", false

		for l := range buffer.Hex(h, tail, cmd.Mode, cli.Profile).Lines {
			if cli.Filter != nil && !cli.Filter.MatchString(l.Hex) {
				continue // not matched afterward
			}

			l.Hex = text.MarkMatch(l.Hex, cli.Filter)

			// cut similar lines for better readability
			if l.Hex == lastHex && cmd.Mode != types.Raw {
				if !wasCut {
					wasCut = true
					_, _ = fmt.Fprintln(cli.Stdout, text.Hide("*"))
				}
				continue
			}

			switch cmd.Mode {
			case types.Canonical:
				_, _ = fmt.Fprintf(cli.Stdout, "%s  %s%s\n", text.Hide(l.Nr), l.Hex, text.Hide(l.Str))
			case types.Hexdump:
				_, _ = fmt.Fprintf(cli.Stdout, "%s %s\n", text.Hide(l.Nr), l.Hex)
			case types.Xxd:
				_, _ = fmt.Fprintf(cli.Stdout, "%s %s %-16s\n", text.Hide(l.Nr), l.Hex, text.Hide(l.Str))
			case types.Raw:
				_, _ = fmt.Fprintf(cli.Stdout, "%s\n", l.Hex)
			}

			lastHex, wasCut = l.Hex, false
		}

		h.Discard()
	}

	return nil
}
