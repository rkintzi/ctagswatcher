package project

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/rkintzi/ctagswatcher/cmd/ctags-indexer/internal/dirwatcher"
	"github.com/rkintzi/ctagswatcher/cmd/ctags-indexer/internal/globfilter"
	"github.com/rkintzi/ctagswatcher/cmd/ctags-indexer/internal/project/cfg"
	"github.com/rkintzi/ctagswatcher/cmdrunner"
	"github.com/rkintzi/ctagswatcher/ctags"
	"github.com/rkintzi/ctagswatcher/taskqueue"
)

type Project struct {
	c      *cfg.ProjectConf
	q      *taskqueue.Q
	dws    []*dirwatcher.DirectoryWatcher
	filter dirwatcher.Filter
	events chan string
	errors chan error
}

func New(c *cfg.ProjectConf, q *taskqueue.Q, exts []string) *Project {
	f := globfilter.New(globfilter.Deny)
	f.Append(globfilter.Allow, exts...)
	p := &Project{
		c:      c,
		dws:    make([]*dirwatcher.DirectoryWatcher, len(c.Dirs)),
		events: make(chan string),
		errors: make(chan error),
		filter: f,
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
	go p.debounce()
	return
}

func (p *Project) debounce() {
	for {
		fs := make(map[string]bool)
		f := <-p.events
		fs[f] = true
		for {
			t := time.NewTimer(100 * time.Millisecond)
			select {
			case f = <-p.events:
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
		go p.reindex(fs)
	}
}

func (p *Project) reindex(fs map[string]bool) {
	files := make([]string, 0, len(fs))
	for f := range fs {
		files = append(files, f)
	}
	args := make([]string, 0, 3+len(fs))
	args = append(args, "ctags", "-f", "-")
	args = append(args, files...)
	cmd := exec.Command(args[0], args[1:len(args)]...)
	f := func(errs chan<- error) {
		var tmp, tags *os.File
		var err error
		defer func() {
			if err != nil {
				errs <- err
			}
			if tmp != nil {
				err = tmp.Close()
				if err != nil {
					errs <- err
				}
				err = os.Remove(tmp.Name())
				if err != nil {
					errs <- err
				}
			}
			if tags != nil {
				err = tags.Close()
				if err != nil {
					errs <- err
				}
			}
		}()
		tmp, err = ioutil.TempFile("", fmt.Sprintf("%s-tags", p.c.Name))
		if err != nil {
			errs <- err
		}
		defer tmp.Close()
		tags, err = os.OpenFile(filepath.Join(p.c.Root, p.c.Tags), os.O_RDWR, 0666)
		if err != nil {
			return
		}
		defer tags.Close()
		err = cmdrunner.Run(cmd, func(r io.Reader) error {
			return ctags.MergeTags(tmp, ctags.NewFilenameFilter(ctags.NewReader(tags), files...), ctags.NewReader(r))
		})
		if err != nil {
			return
		}
		_, err = tmp.Seek(0, os.SEEK_SET)
		if err != nil {
			return
		}
		_, err = tags.Seek(0, os.SEEK_SET)
		if err != nil {
			return
		}
		err = tags.Truncate(0)
		if err != nil {
			return
		}
		_, err = io.Copy(tags, tmp)
		if err != nil {
			return
		}
	}
	p.q.Enqueue(taskqueue.TaskFunc(f))
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
