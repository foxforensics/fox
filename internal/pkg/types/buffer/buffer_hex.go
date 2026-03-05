package buffer

import (
	"fmt"
	"strings"

	cli "github.com/cuhsat/fox/v4/internal/cmd"

	"github.com/cuhsat/fox/v4/internal/pkg/text"
	"github.com/cuhsat/fox/v4/internal/pkg/types/heap"
)

const (
	Canonical HexMode = iota
	Hexdump
	Xxd
	Raw
)

type HexMode int

type HexLine struct {
	Address string
	Values  string
	String  string
}

type HexBuffer struct {
	Lines chan HexLine
}

type HexContext struct {
	Mode    HexMode
	Data    []byte
	Index   int
	Delta   int
	Decimal bool
}

func Hex(h *heap.Heap, cli *cli.Globals, ctx *HexContext) *HexBuffer {
	var buf = &HexBuffer{make(chan HexLine, cli.Threads*1024)}

	if cli.Tail {
		ctx.Delta = cli.Limit.Offset.Bytes
	}

	ctx.Data = h.Bytes()

	go streamHex(buf, ctx)

	return buf
}

func streamHex(buf *HexBuffer, ctx *HexContext) {
	defer close(buf.Lines)

	for ; ctx.Index < len(ctx.Data); ctx.Index += 16 {
		switch ctx.Mode {
		case Canonical:
			buf.Lines <- fmtCanonical(ctx)
		case Hexdump:
			buf.Lines <- fmtHexdump(ctx)
		case Xxd:
			buf.Lines <- fmtXxd(ctx)
		case Raw:
			buf.Lines <- fmtRaw(ctx)
		}
	}
}

func fmtCanonical(ctx *HexContext) HexLine {
	var adr string
	var hex strings.Builder
	var str strings.Builder

	if ctx.Decimal {
		adr = fmt.Sprintf("%08d", ctx.Delta+ctx.Index)
	} else {
		adr = fmt.Sprintf("%08x", ctx.Delta+ctx.Index)
	}

	var l int
	for j := range 16 {
		if ctx.Index+j >= len(ctx.Data) {
			break
		}

		hex.WriteString(fmtHex(ctx.Data[ctx.Index+j]))
		l += 2

		if (j+1)%8 == 0 {
			hex.WriteString("  ")
			l += 2
		} else {
			hex.WriteByte(' ')
			l++
		}

		str.WriteString(fmt.Sprintf("%c", ctx.Data[ctx.Index+j]))
	}

	if l < 50 {
		hex.WriteString(strings.Repeat(" ", 50-l))
	}

	return HexLine{
		adr,
		hex.String(),
		fmt.Sprintf("%-16s", text.ToAscii(str.String(), "·")),
	}
}

func fmtHexdump(ctx *HexContext) HexLine {
	var adr string
	var hex strings.Builder

	if ctx.Decimal {
		adr = fmt.Sprintf("%07d", ctx.Delta+ctx.Index)
	} else {
		adr = fmt.Sprintf("%07x", ctx.Delta+ctx.Index)
	}

	for j := range 16 {
		if ctx.Index+j >= len(ctx.Data) {
			break
		}

		hex.WriteString(fmtHex(ctx.Data[ctx.Index+j]))

		if (j+1)%2 == 0 {
			hex.WriteByte(' ')
		}
	}

	return HexLine{adr, fmt.Sprintf("%-*s", 50, hex.String()), ""}
}

func fmtXxd(ctx *HexContext) HexLine {
	var adr string
	var hex strings.Builder
	var str strings.Builder

	if ctx.Decimal {
		adr = fmt.Sprintf("%08d:", ctx.Delta+ctx.Index)
	} else {
		adr = fmt.Sprintf("%08x:", ctx.Delta+ctx.Index)
	}

	var l int
	for j := range 16 {
		if ctx.Index+j >= len(ctx.Data) {
			break
		}

		hex.WriteString(fmtHex(ctx.Data[ctx.Index+j]))
		l += 2

		if (j+1)%2 == 0 {
			hex.WriteByte(' ')
			l++
		}

		str.WriteString(fmt.Sprintf("%c", ctx.Data[ctx.Index+j]))
	}

	if l < 40 {
		hex.WriteString(strings.Repeat(" ", 40-l))
	}

	return HexLine{
		adr,
		hex.String(),
		text.ToAscii(str.String(), "."),
	}
}

func fmtRaw(ctx *HexContext) HexLine {
	var hex strings.Builder

	for j := range 16 {
		if ctx.Index+j >= len(ctx.Data) {
			break
		}

		hex.WriteString(fmtHex(ctx.Data[ctx.Index+j]))
		hex.WriteByte(' ')
	}

	return HexLine{"", hex.String(), ""}
}

func fmtHex(b byte) string {
	s := fmt.Sprintf("%02x", b)

	switch {
	case b == 0:
		return text.HexZero(s)
	case b < 0x20:
		return text.HexLow(s)
	default:
		return text.HexHigh(s)
	}
}
