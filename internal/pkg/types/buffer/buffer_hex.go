package buffer

import (
	"fmt"
	"strings"

	cli "go.foxforensics.dev/fox/v4/internal/cmd"

	"go.foxforensics.dev/fox/v4/internal/pkg/text"
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
	Data   []byte
	Delta  int
	Index  int
	Pretty bool
}

func Hex(cli *cli.Globals, ctx *HexContext) *HexBuffer {
	var buf = &HexBuffer{make(chan HexLine, cli.Parallel*1024)}

	if cli.Tail {
		ctx.Delta = cli.Limit.Offset.Bytes
	}

	go streamHex(buf, ctx)

	return buf
}

func streamHex(buf *HexBuffer, ctx *HexContext) {
	defer close(buf.Lines)

	for ; ctx.Index < len(ctx.Data); ctx.Index += 16 {
		if ctx.Pretty {
			buf.Lines <- formatStd(ctx)
		} else {
			buf.Lines <- formatRaw(ctx)
		}
	}
}

func formatStd(ctx *HexContext) HexLine {
	var adr = fmt.Sprintf("%08x", ctx.Delta+ctx.Index)
	var val strings.Builder
	var str strings.Builder

	for i := range 16 {
		if ctx.Index+i >= len(ctx.Data) {
			break
		}

		val.WriteString(fmt.Sprintf("%02x ", ctx.Data[ctx.Index+i]))

		if (i+1)%8 == 0 {
			val.WriteString(" ") // middle separator
		}

		str.WriteString(fmt.Sprintf("%c", ctx.Data[ctx.Index+i]))
	}

	val.WriteString(strings.Repeat(" ", max(0, 50-val.Len())))

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
