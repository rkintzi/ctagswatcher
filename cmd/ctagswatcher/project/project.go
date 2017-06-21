package main

import (
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/rkintzi/ctagswatcher/cmd/ctagwatcher/cfg"
)

type Project struct {
	c *cfg.ProjectConf
	w *fsnotify.Watcher
}

func New(c *cfg.ProjectConf) *Project {
	return &Project{c: c}
}

func (p *Project) Watch() (err error) {
	p.w, err = fsnotify.NewWatcher()
	if err != nil {
		return
	}
	for i, dir := range p.c.Dirs {
		path := filepath.Join(p.c.Root, dir)
		err = p.watch(path)
		if err != nil {
			return
		}
	}
}

func (p *Project) watch(path string) error {
}
