package ctags

import "bytes"

type TagLine struct {
	l  []byte
	te int
	fe int
}

func ParseTagLine(line []byte) *TagLine {
	te := bytes.Index(line, []byte("\t"))
	fe := bytes.Index(line[te+1:], []byte("\t")) + te + 1
	return &TagLine{l: line, te: te, fe: fe}
}

func (l *TagLine) IsEmpty() bool    { return len(l.l) == 0 }
func (l *TagLine) IsComment() bool  { return l.l[0] == '!' }
func (l *TagLine) Tag() []byte      { return l.l[0:l.te] }
func (l *TagLine) Filename() []byte { return l.l[l.te+1 : l.fe] }
func (l *TagLine) Bytes() []byte    { return l.l }
