package pager

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/cuhsat/fox/v4/internal/pkg/text"
	"github.com/mattn/go-runewidth"
	"golang.org/x/term"
)

type Pager struct {
	stdout io.WriteCloser
	limit  int
	level  int
}

func New() (*Pager, error) {
	w, h, err := term.GetSize(0)

	if err != nil {
		return nil, err
	}

	return &Pager{
		stdout: os.Stdout,
		limit:  w * (h - 1),
		level:  0,
	}, nil
}

func (p *Pager) Close() error {
	return p.stdout.Close()
}

func (p *Pager) Write(b []byte) (int, error) {
	p.ingest(string(b))
	return len(b), nil
}

func (p *Pager) ingest(s string) {
	var sb strings.Builder

	for _, r := range []rune(s) {
		p.level += runewidth.RuneWidth(r)

		if p.level >= p.limit {
			_, _ = fmt.Fprintf(p.stdout, sb.String())

			sb.Reset()
			p.level = 0

			p.pause()
		}

		sb.WriteRune(r)
	}

	_, _ = fmt.Fprintf(p.stdout, sb.String())
}

func (p *Pager) pause() {
	fd := int(os.Stdin.Fd())

	st, err := term.MakeRaw(fd)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer term.Restore(fd, st)

	_, _ = fmt.Fprintf(p.stdout, text.Hide("\nPress ESC to exit or SPACE to continue.\n"))

	b := make([]byte, 16)

	for {
		_, _ = os.Stdin.Read(b)

		if key > 0 {
			break
		}
	}
}
