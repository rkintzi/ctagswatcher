package project

import (
	"os"
	"path/filepath"
	"time"

	"github.com/rkintzi/ctagswatcher/cmd/ctags-indexer/internal/dirwatcher"
	"github.com/rkintzi/ctagswatcher/cmd/ctags-indexer/internal/project/cfg"
)

type Project struct {
	c      *cfg.ProjectConf
	dws    []*dirwatcher.DirectoryWatcher
	events chan string
	errors chan error
}

func New(c *cfg.ProjectConf, exts []string) *Project {
	p := &Project{
		c:      c,
		dws:    make([]*dirwatcher.DirectoryWatcher, len(c.Dirs)),
		events: make(chan string),
		errors: make(chan error),
	}
	return p
}

func (p *Project) Start() (err error) {
	var filter dirwatcher.Filter
	var since time.Time
	tags := filepath.Join(p.c.Root, p.c.Tags)
	fi, err := os.Stat(tags)
	if err != nil {
		return
	}
	since = fi.ModTime()
	for i, d := range p.c.Dirs {
		d = filepath.Join(p.c.Root, d)
		p.dws[i] = dirwatcher.NewDirectoryWatcher(d, p.events, p.errors, since, filter)
		err = p.dws[i].Start()
		if err != nil {
			return
		}
	}
	go p.debunce()
	return
}

func (p *Project) debounce() {
	for {
		fs := make(map[string]bool)
		f <- p.events
		fs[f] = true
		for {
			t := time.NewTimer(100 * time.Millisecond)
			select {
			case f <- p.events:
				fs[f] = true
				if !t.Stop() {
					<-t.C
				}
				t.Reset(100 * time.Microsecond)
				if len(fs) == 10 {
					break
				}
			case <-t.C:
				break
			}
		}

	}
}

func (p *Project) Close() (err error) {
	for _, w := range p.dws {
		if w == nil {
			break
		}
		e := w.Close()
		if err == nil && e != nil {
			err = e
		}
	}
	return
}
