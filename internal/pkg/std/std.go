package std

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/cuhsat/fox/v4/internal/cmd"
	"github.com/cuhsat/fox/v4/internal/pkg/text"
	"github.com/cuhsat/fox/v4/internal/pkg/types/receipt"
)

var Out io.Writer = os.Stdout

var opts *cmd.Globals

func Init(cli *cmd.Globals) {
	var err error

	switch {
	case len(cli.File) > 0:
		cli.NoPretty = true

		Out, err = os.OpenFile(cli.File, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)

		if err != nil {
			log.Fatalln(err)
		}

	case cli.Quiet:
		log.SetOutput(io.Discard)
		Out, _ = os.Open(os.DevNull)
	}

	opts = cli
}

func Title(s string) {
	_, _ = fmt.Fprintln(Out, text.Title(s))
}

func Write(f string, a ...any) {
	if !opts.NoPretty {
		Writebc(f, a...)
	} else {
		Writeln(f, a...)
	}
}

func Writeln(f string, a ...any) {
	_, _ = fmt.Fprintf(Out, f+"\n", a...)
}

func Writebc(f string, a ...any) {
	_, _ = fmt.Fprintf(Out, text.Border+" "+f+"\n", a...)
}

func Close() {
	if v, is := Out.(io.Closer); is {
		_ = v.Close()
	}

	if !opts.NoReceipt && len(opts.File) > 0 {
		err := receipt.Generate(opts.File)

		if err != nil {
			log.Println(err)
		}
	}
}
