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
	Offset string
	Values string
	String string
}

type HexBuffer struct {
	Lines chan HexLine
}

func Hex(h *heap.Heap, cli *cli.Globals, mode HexMode) *HexBuffer {
	var buf = &HexBuffer{make(chan HexLine, cli.Threads*1024)}
	var delta int

	if cli.Tail {
		delta = cli.Limit.Offset.Bytes
	}

	go streamHex(buf, mode, h.Bytes(), delta)

	return buf
}

func streamHex(buf *HexBuffer, mode HexMode, b []byte, d int) {
	defer close(buf.Lines)

	for i := 0; i < len(b); i += 16 {
		switch mode {
		case Canonical:
			buf.Lines <- fmtCanonical(b, i, d)
		case Hexdump:
			buf.Lines <- fmtHexdump(b, i, d)
		case Xxd:
			buf.Lines <- fmtXxd(b, i, d)
		case Raw:
			buf.Lines <- fmtRaw(b, i, d)
		}
	}
}

func fmtCanonical(b []byte, i int, d int) HexLine {
	var hex strings.Builder
	var str strings.Builder

	for j := range 16 {
		if i+j >= len(b) {
			break
		}

		hex.WriteString(fmt.Sprintf("%02x", b[i+j]))
		str.WriteString(fmt.Sprintf("%c", b[i+j]))

		if j+1%8 == 0 {
			hex.WriteString("  ")
		} else {
			hex.WriteString(" ")
		}
	}

	return HexLine{
		fmt.Sprintf("%08x", d+i),
		fmt.Sprintf("%-*s", 50, hex.String()),
		fmt.Sprintf("|%-16s|", text.ToAscii(str.String())),
	}
}

func fmtHexdump(b []byte, i int, d int) HexLine {
	var hex strings.Builder

	for j := range 16 {
		if i+j >= len(b) {
			break
		}

		hex.WriteString(fmt.Sprintf("%02x", b[i+j]))

		if j+1%2 == 0 {
			hex.WriteString(" ")
		}
	}

	return HexLine{
		fmt.Sprintf("%07x", d+i),
		fmt.Sprintf("%-*s", 50, hex.String()),
		"",
	}
}

func fmtXxd(b []byte, i int, d int) HexLine {
	var hex strings.Builder
	var str strings.Builder

	for j := range 16 {
		if i+j >= len(b) {
			break
		}

		hex.WriteString(fmt.Sprintf("%02x", b[i+j]))
		str.WriteString(fmt.Sprintf("%c", b[i+j]))

		if j+1%2 == 0 {
			hex.WriteString(" ")
		}
	}

	return HexLine{
		fmt.Sprintf("%08x:", d+i),
		fmt.Sprintf("%-*s", 40, hex.String()),
		text.ToAscii(str.String()),
	}
}

func fmtRaw(b []byte, i int, _ int) HexLine {
	var hex strings.Builder

	for j := range 16 {
		if i+j >= len(b) {
			break
		}

		hex.WriteString(fmt.Sprintf("%02x ", b[i+j]))
	}

	return HexLine{"", hex.String(), ""}
}
