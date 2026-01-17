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
	Nr  string
	Grp uint
	Str string
}

type TextBuffer struct {
	Lines chan *TextLine
	Pad   uint
}

func Text(h *heap.Heap, cli *cli.Globals) *TextBuffer {
	s := smap.Map(h.Bytes())
	s = cli.Limit.ReduceSMap(s)
	s = cli.Filter.FilterSMap(s)

	tb := &TextBuffer{
		make(chan *TextLine, cli.Profile*1024),
		uint(math.Log10(float64(len(s)))) + 1,
	}

	go textStream(tb, s.Render())

	return tb
}

func textLine(nr, str string, grp uint) *TextLine {
	return &TextLine{nr, grp, str}
}

func textStream(tb *TextBuffer, s smap.SMap) {
	defer close(tb.Lines)

	var numSep uint = 0
	var numGrp uint = 1
	var tmpGrp uint = 0

	for _, str := range s {

		// insert context separator
		if tmpGrp != str.Grp && numGrp > 1 {
			tb.Lines <- textLine(Sep, "", str.Grp)
			numGrp = 1
			numSep++
		}

		// build line
		tb.Lines <- textLine(
			fmt.Sprintf("%0*d ", tb.Pad, str.Nr),
			text.Sanitize(str.Str),
			str.Grp,
		)

		tmpGrp = str.Grp
		numGrp++
	}
}
