package pager

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/term"
)

const enotty = "inappropriate ioctl for device"
const escape = 0x1B

type Pager struct {
	limit int
	count int
}

func New(limit int) (*Pager, error) {
	_, h, err := term.GetSize(0)

	if err != nil {
		if err.Error() == enotty {
			err = errors.New("can't use - with pause")
		}

		return nil, err
	}

	if limit == 0 {
		limit = h
	}

	return &Pager{
		limit: limit,
	}, nil
}

func (p *Pager) Close() error {
	return nil
}

func (p *Pager) Write(b []byte) (int, error) {
	for _, s := range strings.SplitAfter(string(b), "\n") {
		if p.count == p.limit {
			p.count = 0
			pause()
		}

		p.count++

		_, _ = fmt.Fprint(os.Stdout, s)
	}

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
		case escape, 'c', 'q':
			_ = term.Restore(fd, fs)
			os.Exit(0)

		default:
			_ = term.Restore(fd, fs)
			return
		}
	}
}
