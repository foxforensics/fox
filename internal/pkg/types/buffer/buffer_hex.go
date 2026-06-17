package buffer

import (
	"context"
	"fmt"
	"strings"

	cli "go.foxforensics.eu/fox/v4/internal/cmd"

	"go.foxforensics.eu/fox/v4/internal/pkg/text"
)

type HexLine struct {
	Address string
	Values  string
	String  string
}

type HexBuffer struct {
	Lines chan HexLine
}

type HexContext struct {
	Parent context.Context
	Data   []byte
	Delta  int
	Index  int
	Pretty bool
}

func Hex(ctx *HexContext, cli *cli.Globals) *HexBuffer {
	var buf = &HexBuffer{make(chan HexLine, cli.Threads*1024)}

	if cli.Limits.IsTail {
		ctx.Delta = cli.Limits.Values.Bytes
	}

	go streamHex(ctx, buf)

	return buf
}

func streamHex(ctx *HexContext, buf *HexBuffer) {
	defer close(buf.Lines)

	for ; ctx.Index < len(ctx.Data); ctx.Index += 16 {
		select {
		case <-ctx.Parent.Done():
			return
		default:
			if ctx.Pretty {
				buf.Lines <- formatStd(ctx)
			} else {
				buf.Lines <- formatRaw(ctx)
			}
		}
	}
}

func formatStd(ctx *HexContext) HexLine {
	var adr = fmt.Sprintf("%08x", ctx.Delta+ctx.Index)
	var val strings.Builder
	var str strings.Builder

	j := 0

	for i := range 16 {
		if ctx.Index+i >= len(ctx.Data) {
			break
		}

		switch v := ctx.Data[ctx.Index+i]; {
		case v == 0:
			val.WriteString(text.AsGray("%02x ", v))
		case v >= 1 && v <= 31:
			val.WriteString(text.AsBold("%02x ", v))
		default:
			val.WriteString(fmt.Sprintf("%02x ", v))
		}

		if (i+1)%8 == 0 {
			val.WriteString(" ") // middle separator
			j += 4
		} else {
			j += 3
		}

		str.WriteString(fmt.Sprintf("%c", ctx.Data[ctx.Index+i]))
	}

	val.WriteString(strings.Repeat(" ", max(0, 50-j)))

	return HexLine{adr, val.String(), text.ToAscii(str.String(), "·")}
}

func formatRaw(ctx *HexContext) HexLine {
	var val strings.Builder

	for i := range 16 {
		if ctx.Index+i >= len(ctx.Data) {
			break
		}

		val.WriteString(fmt.Sprintf("%02x ", ctx.Data[ctx.Index+i]))
	}

	return HexLine{"", val.String(), ""}
}
