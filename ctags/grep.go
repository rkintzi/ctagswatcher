package ctags

import (
	"bufio"
	"bytes"
	"io"
)

type filter struct {
	r   *bufio.Reader
	err error
	fn  []byte
	b   []byte
	i   int
}

func NewFilenameFilter(r io.Reader, fn string) io.Reader {
	return &filter{r: bufio.NewReader(r), fn: []byte(fn)}
}

func (f *filter) Read(p []byte) (n int, err error) {
	n, err = f.copyBuf(p)
	if n > 0 {
		return
	}
	for {
		f.i = 0
		f.b, f.err = f.r.ReadBytes('\n')
		var tl tagline
		tl.line = f.b
		tl.ti = bytes.Index(tl.line, []byte("\t"))
		tl.fi = bytes.Index(tl.line[tl.ti+1:], []byte("\t")) + tl.ti + 1
		if !tl.IsEmpty() && !tl.IsComment() && bytes.Compare(tl.Filename(), f.fn) == 0 {
			if f.err != nil {
				return 0, f.err
			}
			continue
		}
		n, err = f.copyBuf(p)
		return
	}
}

func (f *filter) copyBuf(p []byte) (n int, err error) {
	n = copy(p, f.b[f.i:])
	f.i += n
	if f.i == len(f.b) {
		return n, f.err
	}
	return
}
