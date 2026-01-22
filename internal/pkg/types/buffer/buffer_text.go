package buffer

import (
	"fmt"
	"math"

	cli "github.com/cuhsat/fox/v4/internal/cmd"

	"github.com/cuhsat/fox/v4/internal/pkg/text"
	"github.com/cuhsat/fox/v4/internal/pkg/types/heap"
	"github.com/cuhsat/fox/v4/internal/pkg/types/smap"
)

type TextLine struct {
	Line   string
	Group  uint
	String string
}

type TextBuffer struct {
	Lines chan *TextLine
	Pad   uint
}

func Text(h *heap.Heap, cli *cli.Globals) *TextBuffer {
	s := cli.Filter.Filter(smap.Map(h.Bytes()))

	var buf = &TextBuffer{
		make(chan *TextLine, cli.Profile*1024),
		uint(math.Log10(float64(len(s)))) + 1,
	}
	var delta int

	if cli.Tail {
		delta = cli.Limit.Offset.Lines
	}

	go streamText(buf, s.Render(), delta)

	return buf
}

func streamText(buf *TextBuffer, s smap.SMap, d int) {
	defer close(buf.Lines)

	var numSep uint = 0
	var numGrp uint = 1
	var tmpGrp uint = 0

	for _, str := range s {

		// insert context separator
		if tmpGrp != str.Group && numGrp > 1 {
			buf.Lines <- &TextLine{Sep, str.Group, ""}
			numGrp = 1
			numSep++
		}

		// build line
		buf.Lines <- &TextLine{
			fmt.Sprintf("%0*d ", buf.Pad, uint(d)+str.Line),
			str.Group,
			text.Sanitize(string(str.Bytes)),
		}

		tmpGrp = str.Group
		numGrp++
	}
}
