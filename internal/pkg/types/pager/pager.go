package pager

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mattn/go-runewidth"
	"golang.org/x/term"

	"github.com/cuhsat/fox/v4/internal/pkg/text"
)

const help = "Press space to continue, 'q' to quit."

type Pager struct {
	quiet bool
	width int
	limit int
	level int
}

func New(quiet bool) (*Pager, error) {
	w, h, err := term.GetSize(0)

	if err != nil {
		return nil, err
	}

	return &Pager{
		quiet: quiet,
		width: w,
		limit: w * h,
	}, nil
}

func (p *Pager) Close() error {
	return nil
}

func (p *Pager) Write(b []byte) (int, error) {
	var page strings.Builder

	s := string(b)

	for _, r := range []rune(s) {
		p.level += runewidth.RuneWidth(r)
		page.WriteRune(r)

		if p.level >= (p.limit - p.width) {
			_, _ = fmt.Fprintf(os.Stdout, page.String())

			p.level = 0
			page.Reset()

			p.pause()
		}
	}

	// add remaining line
	if b[len(b)-1] == '\n' {
		p.level += p.width - (runewidth.StringWidth(s) % p.width)
	}

	_, _ = fmt.Fprintf(os.Stdout, page.String())

	return len(b), nil
}

func (p *Pager) pause() {
	if !p.quiet {
		print(text.Hide(help))
		defer print(strings.Repeat("\b", len(help)) + "\033")
	}

	fd := int(os.Stdin.Fd())

	fs, err := term.MakeRaw(fd)

	if err != nil {
		log.Println(err)
		return
	}

	key := make([]byte, 1)

	for {
		_, _ = os.Stdin.Read(key)

		switch key[0] {
		case 'q':
			_ = term.Restore(fd, fs)
			os.Exit(0)

		case ' ':
			_ = term.Restore(fd, fs)
			return
		}
	}
}
