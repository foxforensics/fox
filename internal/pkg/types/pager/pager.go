package pager

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mattn/go-runewidth"
	"golang.org/x/term"
)

type Pager struct {
	limit int
	width int
	level int
}

func New() (*Pager, error) {
	w, h, err := term.GetSize(0)

	if err != nil {
		return nil, err
	}

	return &Pager{
		limit: w * (h - 1),
		width: w,
	}, nil
}

func (p *Pager) Close() error {
	return nil
}

func (p *Pager) Write(b []byte) (int, error) {
	var page strings.Builder

	s := string(b)

	for _, r := range []rune(s) {
		if p.level >= p.limit {
			_, _ = fmt.Fprintf(os.Stdout, page.String())

			p.level = 0
			page.Reset()

			pause()
		}

		p.level += runewidth.RuneWidth(r)
		page.WriteRune(r)
	}

	// add remaining line
	if b[len(b)-1] == '\n' {
		p.level += p.width - (runewidth.StringWidth(s) % p.width)
	}

	_, _ = fmt.Fprintf(os.Stdout, page.String())

	return len(b), nil
}

func pause() {
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
