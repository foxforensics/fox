package hex

import (
	"strings"

	cli "github.com/cuhsat/fox/v4/internal/cmd"

	"github.com/cuhsat/fox/v4/internal/pkg/text"
	"github.com/cuhsat/fox/v4/internal/pkg/types/buffer"
)

var Usage = strings.TrimSpace(`
Show file contents in hex format.

fox hex [FLAGS...] <PATHS...>

Flags:
  -C, --canonical          Format output as canonical
  -H, --hexdump            Format output like hexdump
  -X, --xxd                Format output like xxd

Format flags:
  -D, --decimal            Format addresses as decimal

Disable flags:
  -R, --no-format          Don't format output at all

Examples:
  $ fox hex -hc512 disk.dd
`)

type Hex struct {
	Canonical bool `short:"C" xor:"canonical,hexdump,xxd"`
	Hexdump   bool `short:"H" xor:"canonical,hexdump,xxd"`
	Xxd       bool `short:"X" xor:"canonical,hexdump,xxd"`

	// format
	Decimal bool `short:"D"`

	// disable
	NoFormat bool `short:"R" long:"no-format"`

	// paths
	Paths []string `arg:"" optional:""`
}

func (cmd *Hex) Run(cli *cli.Globals) error {
	if len(cmd.Paths)+len(cli.Paths) == 0 {
		return text.Usage(Usage)
	}

	var mode buffer.HexMode

	switch {
	case cmd.Canonical:
		mode = buffer.Canonical
	case cmd.Hexdump:
		mode = buffer.Hexdump
	case cmd.Xxd:
		mode = buffer.Xxd
	case cmd.NoFormat:
		mode = buffer.Raw
	default:
		mode = buffer.Default
	}

	if mode > 0 {
		cli.NoPretty = !cli.NoPretty
	}

	cli.NoConvert = true // forced

	ch := cli.Load(cmd.Paths)
	defer cli.Discard()

	for h := range ch {
		if !cli.NoPretty {
			text.Title(h.String())
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
			if l.Values == lastHex && !cmd.NoFormat {
				if !wasCut {
					wasCut = true
					text.Write(text.AsGray("*"))
				}
				continue
			}

			if mode == buffer.Default {
				text.Write("%s  %s%s", text.AsGray(l.Address), l.Values, l.String)
			} else if mode == buffer.Canonical {
				text.Write("%s  %s|%s|", l.Address, l.Values, l.String)
			} else if mode == buffer.Hexdump {
				text.Write("%s %s", l.Address, l.Values)
			} else if mode == buffer.Xxd {
				text.Write("%s %s %-16s", l.Address, l.Values, l.String)
			} else {
				text.Write(l.Values)
			}

			lastHex, wasCut = l.Values, false
		}

		h.Discard()
	}

	return nil
}
