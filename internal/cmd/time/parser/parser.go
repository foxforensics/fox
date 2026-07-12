package parser

import (
	"context"
	"log/slog"
	"slices"

	"go.foxforensics.eu/fox/v5/internal/cmd/time/entry"
	"go.foxforensics.eu/fox/v5/library/binaries/bin/lnk"
	"go.foxforensics.eu/fox/v5/library/binaries/bin/mft"
	"go.foxforensics.eu/fox/v5/library/binaries/bin/pf"
)

type Options struct {
	Sort    bool
	Threads int
}

type Parser struct {
	opts *Options
}

func New(opts *Options) *Parser {
	return &Parser{
		opts: opts,
	}
}

func (prs *Parser) Parse(ctx context.Context, block []byte) <-chan *entry.Entry {
	ch := make(chan *entry.Entry, prs.opts.Threads*64)

	go func() {
		defer close(ch)

		var fn func([]byte) []entry.Entry

		switch {
		case mft.Detect(block):
			slog.Debug("file detected as mft")
			fn = mft.Parse

		case lnk.Detect(block):
			slog.Debug("file detected as lnk")
			fn = lnk.Parse

		case pf.Detect(block):
			slog.Debug("file detected as pf")
			fn = pf.Parse

		default:
			slog.Debug("format not detected")
			return
		}

		for _, e := range fn(block) {
			select {
			case ch <- &e:
			case <-ctx.Done():
				return
			}
		}
	}()

	if prs.opts.Sort {
		return prs.sort(ctx, ch)
	}

	return ch
}

func (prs *Parser) sort(ctx context.Context, ch <-chan *entry.Entry) <-chan *entry.Entry {
	sorted := make(chan *entry.Entry, cap(ch))

	go func() {
		defer close(sorted)

		v := make([]*entry.Entry, 0)

		for e := range ch {
			v = append(v, e)
		}

		slices.SortStableFunc(v, func(a, b *entry.Entry) int {
			return a.SortKey().Compare(b.SortKey())
		})

		for _, e := range v {
			select {
			case sorted <- e:
			case <-ctx.Done():
				return
			}
		}
	}()

	return sorted
}
