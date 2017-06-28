package grep

import (
	"bufio"
	"io"
)

type grep struct {
	r    *bufio.Reader
	err  error
	pred Pred
	i    int
	line []byte
}

type Pred func([]byte) bool

func NewReader(r io.Reader, p Pred) io.Reader {
	return &grep{r: bufio.NewReader(r), pred: p}
}

func (f *grep) Read(p []byte) (n int, err error) {
	n, err = f.copyBuf(p)
	if n > 0 {
		return
	}
	for {
		f.i = 0
		f.line, f.err = f.r.ReadBytes('\n')
		if !f.pred(f.line) {
			if f.err != nil {
				return 0, f.err
			}
			continue
		}
		n, err = f.copyBuf(p)
		return
	}
}

func (f *grep) copyBuf(p []byte) (n int, err error) {
	n = copy(p, f.line[f.i:])
	f.i += n
	if f.i == len(f.line) {
		return n, f.err
	}
	return
}
