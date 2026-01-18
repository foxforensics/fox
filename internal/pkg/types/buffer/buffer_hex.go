package buffer

import (
	"fmt"
	"strings"

	cli "github.com/cuhsat/fox/v4/internal/cmd"

	"github.com/cuhsat/fox/v4/internal/pkg/text"
	"github.com/cuhsat/fox/v4/internal/pkg/types"
	"github.com/cuhsat/fox/v4/internal/pkg/types/heap"
)

type HexLine struct {
	Offset string
	Values string
	String string
}

type HexBuffer struct {
	Lines chan HexLine
}

func Hex(h *heap.Heap, cli *cli.Globals, mode string) *HexBuffer {
	var buf = &HexBuffer{make(chan HexLine, cli.Parallel*1024)}
	var off uint

	if cli.Tail {
		off = max(uint(h.Size)-cli.Bytes, 0)
	}

	go streamHex(buf, mode, h.Bytes(), off)

	return buf
}

func streamHex(buf *HexBuffer, mode string, b []byte, off uint) {
	defer close(buf.Lines)

	for i := 0; i < len(b); i += 16 {
		switch mode {
		case types.Canonical:
			buf.Lines <- fmtCanonical(b, i, off)
		case types.Hexdump:
			buf.Lines <- fmtHexdump(b, i, off)
		case types.Xxd:
			buf.Lines <- fmtXxd(b, i, off)
		case types.Raw:
			buf.Lines <- fmtRaw(b, i, off)
		}
	}
}

func fmtCanonical(b []byte, i int, off uint) HexLine {
	var hex strings.Builder
	var str strings.Builder

	for j := range 16 {
		if i+j >= len(b) {
			break
		}

		hex.WriteString(fmt.Sprintf("%02x", b[i+j]))
		str.WriteString(fmt.Sprintf("%c", b[i+j]))

		if uint(j+1)%8 == 0 {
			hex.WriteString("  ")
		} else {
			hex.WriteString(" ")
		}
	}

	return HexLine{
		fmt.Sprintf("%08x", off+uint(i)),
		fmt.Sprintf("%-*s", 50, hex.String()),
		fmt.Sprintf("|%-16s|", text.ToAscii(str.String())),
	}
}

func fmtHexdump(b []byte, i int, off uint) HexLine {
	var hex strings.Builder

	for j := range 16 {
		if i+j >= len(b) {
			break
		}

		hex.WriteString(fmt.Sprintf("%02x", b[i+j]))

		if uint(j+1)%2 == 0 {
			hex.WriteString(" ")
		}
	}

	return HexLine{
		fmt.Sprintf("%07x", off+uint(i)),
		fmt.Sprintf("%-*s", 50, hex.String()),
		"",
	}
}

func fmtXxd(b []byte, i int, off uint) HexLine {
	var hex strings.Builder
	var str strings.Builder

	for j := range 16 {
		if i+j >= len(b) {
			break
		}

		hex.WriteString(fmt.Sprintf("%02x", b[i+j]))
		str.WriteString(fmt.Sprintf("%c", b[i+j]))

		if uint(j+1)%2 == 0 {
			hex.WriteString(" ")
		}
	}

	return HexLine{
		fmt.Sprintf("%08x:", off+uint(i)),
		fmt.Sprintf("%-*s", 40, hex.String()),
		text.ToAscii(str.String()),
	}
}

func fmtRaw(b []byte, i int, _ uint) HexLine {
	var hex strings.Builder

	for j := range 16 {
		if i+j >= len(b) {
			break
		}

		hex.WriteString(fmt.Sprintf("%02x ", b[i+j]))
	}

	return HexLine{"", hex.String(), ""}
}
