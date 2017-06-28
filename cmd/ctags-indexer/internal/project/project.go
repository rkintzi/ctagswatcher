package project

import (
	"os"
	"path/filepath"
	"time"

	"github.com/rkintzi/ctagswatcher/cmd/ctags-indexer/internal/project/cfg"
	"github.com/rkintzi/ctagswatcher/cmd/ctags-indexer/internal/watcher"
)

type Project struct {
	c      *cfg.ProjectConf
	dws    []*watcher.DirectoryWatcher
	events chan string
	errors chan error
}

func New(c *cfg.ProjectConf, exts []string) *Project {
	p := &Project{
		c:      c,
		dws:    make([]*watcher.DirectoryWatcher, len(c.Dirs)),
		events: make(chan string),
		errors: make(chan error),
	}
	return p
}

func (p *Project) Start() (err error) {
	var filter watcher.Filter
	var since time.Time
	tags := filepath.Join(p.c.Root, p.c.Tags)
	fi, err := os.Stat(tags)
	if err != nil {
		return
	}
	since = fi.ModTime()
	for i, d := range p.c.Dirs {
		d = filepath.Join(p.c.Root, d)
		p.dws[i] = watcher.NewDirectoryWatcher(d, p.events, p.errors, since, filter)
		err = p.dws[i].Start()
		if err != nil {
			return
		}
	}
	return
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
