package ctags

import (
	"bytes"
	"io"
)

type scanner struct {
	r   Reader
	err error
	eof bool
	tl  *TagLine
}

func (s *scanner) Scan() bool {
	if s.eof {
		s.err = io.EOF
	}
	if s.err != nil {
		return false
	}
	var err error
	s.tl, err = s.r.ReadTag()
	if err == io.EOF && !s.tl.IsEmpty() {
		s.eof = true
	} else if err != nil {
		s.err = err
		return false
	}
	return true
}
func (s *scanner) Err() error {
	return s.err
}
func (s *scanner) Tagline() *TagLine {
	return s.tl
}

func skipHeader(s *scanner) error {
	for s.Scan() {
		if !s.Tagline().IsComment() {
			break
		}
	}
	return s.Err()
}

func copyHeader(w io.Writer, s *scanner) error {
	for s.Scan() {
		if !s.Tagline().IsComment() {
			break
		}
		_, err := w.Write(s.Tagline().Bytes())
		if err != nil {
			return err
		}
	}
	return s.Err()
}

func cmpTags(l, r *TagLine) int {
	if l.IsEmpty() {
		return -1
	} else if r.IsEmpty() {
		return -1
	} else if l.IsComment() {
		return -1
	} else if r.IsComment() {
		return 1
	} else {
		return bytes.Compare(l.Bytes(), r.Bytes())
	}
}

func MergeTags(w io.Writer, rs ...Reader) error {
	scanners := make([]scanner, len(rs))
	for i, r := range rs {
		scanners[i].r = r
		if i == 0 {
			err := copyHeader(w, &scanners[i])
			if err != nil && err != io.EOF {
				return err
			}
		} else {
			err := skipHeader(&scanners[i])
			if err != nil && err != io.EOF {
				return err
			}
		}
	}
	var previous TagLine
	for {
		var min TagLine
		var selected = -1
		for i := range scanners {
			s := &scanners[i]
			if s.Err() == io.EOF {
				continue
			}
			if cmpTags(s.Tagline(), &min) == -1 {
				min = *s.Tagline()
				selected = i
			}
		}
		if selected == -1 {
			break
		}
		if res := cmpTags(&previous, &min); res == -1 && !min.IsEmpty() {
			_, err := w.Write(min.Bytes())
			if err != nil {
				return err
			}
		} else if res == 1 {
			panic("Bug")
		}
		scanners[selected].Scan()
		if err := scanners[selected].Err(); err != nil && err != io.EOF {
			return err
		}
		previous = min
	}
	return nil
}
