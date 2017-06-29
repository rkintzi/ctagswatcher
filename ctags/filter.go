package ctags

import (
	"bytes"
	"io"
)

func filenamePred(filenames []string) Predicate {
	names := make([][]byte, len(filenames))
	for i, fn := range filenames {
		names[i] = []byte(fn)
	}
	return func(tl *TagLine) bool {
		for _, name := range names {
			if bytes.Compare(tl.Filename(), name) == 0 {
				return false
			}
		}
		return true
	}
}

type Predicate func(*TagLine) bool

type filter struct {
	r Reader
	p Predicate
}

func (f *filter) ReadTag() (*TagLine, error) {
	tl, err := f.r.ReadTag()
	for {
		if f.p(tl) {
			return tl, err
		} else if err != nil {
			return &TagLine{}, err
		}
	}
}

func NewFilenameFilter(r io.Reader, filenames ...string) Reader {
	return &filter{r: NewReader(r), p: filenamePred(filenames)}
}
