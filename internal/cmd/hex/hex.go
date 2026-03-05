package hex

import (
	"fmt"
	"strings"

	cli "github.com/cuhsat/fox/v4/internal/cmd"

	"github.com/cuhsat/fox/v4/internal/pkg/text"
	"github.com/cuhsat/fox/v4/internal/pkg/types/buffer"
)

var Usage = strings.TrimSpace(`
Show file contents in hex format.

fox hex [FLAGS...] <PATHS...>

Flags:
  -H, --hexdump            Format output like hexdump
  -X, --xxd                Format output like xxd
  -R, --raw                Don't format output

Format flags:
  -D, --decimal            Format addresses as decimal

Examples:
  $ fox hex -hc512 disk.dd
`)

type Hex struct {
	Hexdump bool `short:"H" xor:"hexdump,xxd,raw"`
	Xxd     bool `short:"X" xor:"hexdump,xxd,raw"`
	Raw     bool `short:"R" xor:"hexdump,xxd,raw"`

	// format
	Decimal bool `short:"D"`

	// paths
	Paths []string `arg:"" optional:""`
}

func (cmd *Hex) Run(cli *cli.Globals) error {
	if len(cmd.Paths)+len(cli.Paths) == 0 {
		fmt.Println(Usage)
		return nil
	}

	var mode buffer.HexMode

	switch {
	case cmd.Hexdump:
		mode = buffer.Hexdump
	case cmd.Xxd:
		mode = buffer.Xxd
	case cmd.Raw:
		mode = buffer.Raw
	default:
		mode = buffer.Canonical
	}

	cli.NoConvert = true // forced

	ch := cli.Load(cmd.Paths)
	defer cli.Discard()

	for h := range ch {
		if !cli.NoFile {
			_, _ = fmt.Fprintf(cli.Stdout, "%s\n", text.Title(h.String()))
		}

		lastHex, wasCut := "", false

		for l := range buffer.Hex(h, cli, &buffer.HexContext{
			Mode: mode, Decimal: cmd.Decimal,
		}).Lines {
			if cli.Regexp != nil && !cli.Regexp.MatchString(l.Values) {
				continue // not matched afterward
			}

			l.Values = text.MarkMatch(l.Values, cli.Regexp)

			// cut similar lines for better readability
			if l.Values == lastHex && !cmd.Raw {
				if !wasCut {
					wasCut = true
					_, _ = fmt.Fprintln(cli.Stdout, text.Hide("*"))
				}
				continue
			}

			switch mode {
			case buffer.Canonical:
				_, _ = fmt.Fprintf(cli.Stdout, "%s  %s%s\n", text.Hide(l.Address), l.Values, l.String)
			case buffer.Hexdump:
				_, _ = fmt.Fprintf(cli.Stdout, "%s %s\n", text.Hide(l.Address), l.Values)
			case buffer.Xxd:
				_, _ = fmt.Fprintf(cli.Stdout, "%s %s %-16s\n", text.Hide(l.Address), l.Values, l.String)
			case buffer.Raw:
				_, _ = fmt.Fprintf(cli.Stdout, "%s\n", l.Values)
			}

			lastHex, wasCut = l.Values, false
		}

		h.Discard()
	}

	return nil
}
