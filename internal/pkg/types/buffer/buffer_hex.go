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

	for j := range 16 {
		if ctx.Index+j >= len(ctx.Data) {
			break
		}

		hex.WriteString(fmt.Sprintf("%02x", ctx.Data[ctx.Index+j]))
		str.WriteString(fmt.Sprintf("%c", ctx.Data[ctx.Index+j]))

		if j+1%8 == 0 {
			hex.WriteString("  ")
		} else {
			hex.WriteString(" ")
		}
	}

	if ctx.Decimal {
		adr = fmt.Sprintf("%08d", ctx.Delta+ctx.Index)
	} else {
		adr = fmt.Sprintf("%08x", ctx.Delta+ctx.Index)
	}

	return HexLine{
		adr,
		fmt.Sprintf("%-*s", 50, hex.String()),
		fmt.Sprintf("|%-16s|", text.ToAscii(str.String())),
	}
}

func fmtHexdump(ctx *HexContext) HexLine {
	var adr string
	var hex strings.Builder

	for j := range 16 {
		if ctx.Index+j >= len(ctx.Data) {
			break
		}

		hex.WriteString(fmt.Sprintf("%02x", ctx.Data[ctx.Index+j]))

		if j+1%2 == 0 {
			hex.WriteString(" ")
		}
	}

	if ctx.Decimal {
		adr = fmt.Sprintf("%07d", ctx.Delta+ctx.Index)
	} else {
		adr = fmt.Sprintf("%07x", ctx.Delta+ctx.Index)
	}

	return HexLine{adr, fmt.Sprintf("%-*s", 50, hex.String()), ""}
}

func fmtXxd(ctx *HexContext) HexLine {
	var adr string
	var hex strings.Builder
	var str strings.Builder

	for j := range 16 {
		if ctx.Index+j >= len(ctx.Data) {
			break
		}

		hex.WriteString(fmt.Sprintf("%02x", ctx.Data[ctx.Index+j]))
		str.WriteString(fmt.Sprintf("%c", ctx.Data[ctx.Index+j]))

		if j+1%2 == 0 {
			hex.WriteString(" ")
		}
	}

	if ctx.Decimal {
		adr = fmt.Sprintf("%08d:", ctx.Delta+ctx.Index)
	} else {
		adr = fmt.Sprintf("%08x:", ctx.Delta+ctx.Index)
	}

	return HexLine{
		adr,
		fmt.Sprintf("%-*s", 40, hex.String()),
		text.ToAscii(str.String()),
	}
}

func fmtRaw(ctx *HexContext) HexLine {
	var hex strings.Builder

	for j := range 16 {
		if ctx.Index+j >= len(ctx.Data) {
			break
		}

		hex.WriteString(fmt.Sprintf("%02x ", ctx.Data[ctx.Index+j]))
	}

	return HexLine{"", hex.String(), ""}
}
