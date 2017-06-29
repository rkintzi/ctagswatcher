package ctags

import (
	"bufio"
	"io"
)

type Reader interface {
	ReadTag() (*TagLine, error)
}

type reader struct {
	r *bufio.Reader
}

func NewReader(r io.Reader) Reader {
	return &reader{r: bufio.NewReader(r)}
}

func (r *reader) ReadTag() (*TagLine, error) {
	l, err := r.r.ReadBytes('\n')
	return ParseTagLine(l), err
}
