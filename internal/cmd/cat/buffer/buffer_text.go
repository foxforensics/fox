package buffer

import (
	"context"
	"fmt"
	"math"

	"go.foxforensics.eu/fox/v5/internal/cmd"
	"go.foxforensics.eu/fox/v5/internal/pkg"
	"go.foxforensics.eu/fox/v5/internal/pkg/smap"
	"go.foxforensics.eu/fox/v5/internal/pkg/writer"
)

const Limit = 1024 * 1024 * 4 // 4mb

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
	Parent context.Context
	SMap   smap.SMap
	Data   []byte
	Delta  int
	Hint   string
}

func Text(ctx *TextContext, fox *cmd.Globals) *TextBuffer {
	var data = ctx.Data
	var last uint

	// turn off color for big files. color must be applied to raw data,
	// so that grep results may still have syntax highlighting.
	if len(data) < Limit {
		data = []byte(writer.ColorizeAs(string(data), ctx.Hint))
	}

	ctx.SMap = smap.Map(data)
	ctx.SMap = fox.Query.Filter(ctx.SMap, fox.Threads).Render(fox.Threads)

	if len(ctx.SMap) > 0 {
		last = ctx.SMap[len(ctx.SMap)-1].Line
	}

	var buf = &TextBuffer{
		make(chan *TextLine, fox.Threads*1024),
		uint(math.Log10(float64(max(1, last)))) + 1,
	}

	if fox.Query.IsTail {
		ctx.Delta = fox.Query.Values.Lines
	}

	go streamText(ctx, buf)

	return buf
}

func streamText(ctx *TextContext, buf *TextBuffer) {
	defer close(buf.Lines)

	var num uint = 1
	var grp uint = 0

	for _, str := range ctx.SMap {
		select {
		case <-ctx.Parent.Done():
			return
		default:
			// insert context separator
			if grp != str.Group && num > 1 {
				buf.Lines <- &TextLine{Sep, str.Group, ""}
				num = 1
			}

			// build line
			buf.Lines <- &TextLine{
				fmt.Sprintf("%0*d ", buf.Pad, uint(ctx.Delta)+str.Line),
				str.Group,
				pkg.Sanitize(string(str.Bytes)),
			}

			grp = str.Group
			num++
		}
	}
}
