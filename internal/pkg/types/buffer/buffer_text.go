package buffer

import (
	"fmt"
	"math"

	cli "github.com/cuhsat/fox/v4/internal/cmd"

	"github.com/cuhsat/fox/v4/internal/pkg/text"
	"github.com/cuhsat/fox/v4/internal/pkg/types/heap"
	"github.com/cuhsat/fox/v4/internal/pkg/types/register"
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

type TextContext struct {
	SMap  smap.SMap
	Delta int
}

func Text(h *heap.Heap, cli *cli.Globals, ctx *TextContext) *TextBuffer {
	b := h.Bytes()

	for _, f := range register.Formats {
		if f.Detect(b) {
			b = f.Format(b)
		}
	}

	ctx.SMap = cli.Filter.Filter(smap.Map(b)).Render()

	var buf = &TextBuffer{
		make(chan *TextLine, cli.Threads*1024),
		uint(math.Log10(float64(len(ctx.SMap)))) + 1,
	}

	if cli.Tail {
		ctx.Delta = cli.Limit.Offset.Lines
	}

	go streamText(buf, ctx)

	return buf
}

func streamText(buf *TextBuffer, ctx *TextContext) {
	defer close(buf.Lines)

	var numSep uint = 0
	var numGrp uint = 1
	var tmpGrp uint = 0

	for _, str := range ctx.SMap {

		// insert context separator
		if tmpGrp != str.Group && numGrp > 1 {
			buf.Lines <- &TextLine{Sep, str.Group, ""}
			numGrp = 1
			numSep++
		}

		// build line
		buf.Lines <- &TextLine{
			fmt.Sprintf("%0*d ", buf.Pad, uint(ctx.Delta)+str.Line),
			str.Group,
			text.Sanitize(string(str.Bytes)),
		}

		tmpGrp = str.Group
		numGrp++
	}
}
