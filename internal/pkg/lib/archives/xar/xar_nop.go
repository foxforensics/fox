package xar

import "io"

type nopCloser struct {
	io.ReaderAt
}

func nop(r io.ReaderAt) nopCloser {
	return nopCloser{r}
}

func (nopCloser) Close() error {
	return nil
}
