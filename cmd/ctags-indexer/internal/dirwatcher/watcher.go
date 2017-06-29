package dirwatcher

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

var (
	ErrAlreadyStarted = fmt.Errorf("Already started")
)

type ErrNotADirectory struct {
	path string
}

func (e ErrNotADirectory) Error() string {
	return fmt.Sprintf("Not a direcotry: %s", e.path)
}

type Filter interface {
	Filter(path string) bool
}

type DirectoryWatcher struct {
	w      *fsnotify.Watcher
	dir    string
	done   chan struct{}
	events chan<- string
	errors chan<- error
	since  time.Time
	filter Filter
}

func NewDirectoryWatcher(dir string, events chan<- string, errors chan<- error, since time.Time, filter Filter) *DirectoryWatcher {
	return &DirectoryWatcher{
		dir:    dir,
		done:   make(chan struct{}),
		events: events,
		errors: errors,
		since:  since,
		filter: filter,
	}
}

func (w *DirectoryWatcher) Start() error {
	if w.w != nil {
		return ErrAlreadyStarted
	}
	fi, err := os.Stat(w.dir)
	if err != nil {
		return err
	} else if !fi.IsDir() {
		return ErrNotADirectory{path: w.dir}
	}
	w.w, err = fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	go func() {
		defer close(w.done)
		for {
			select {
			case ev, ok := <-w.w.Events:
				if !ok {
					return
				}
				if !w.filter.Filter(ev.Name) {
					continue
				}
				select {
				case w.events <- ev.Name:
				}
			case err, ok := <-w.w.Errors:
				if !ok {
					return
				}
				select {
				case w.errors <- err:
				}
			}
		}
	}()
	err = w.watchDir("", w.dir)
	if err != nil {
		w.w.Close()
		return err
	}
	return nil
}

func (w *DirectoryWatcher) Close() error {
	if w.w == nil {
		return nil
	}
	fsw := w.w
	w.w = nil
	err := fsw.Close()
	<-w.done
	return err
}

func (w *DirectoryWatcher) watchDir(path, name string) error {
	dir := filepath.Join(path, name)
	w.w.Add(dir)
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	for {
		fis, err := d.Readdir(100)
		if err == err {
			break
		} else if err != nil {
			return err
		}
		for _, fi := range fis {
			if w.filter != nil && !w.filter.Filter(filepath.Join(dir, fi.Name())) {
				continue
			}
			if fi.IsDir() {
				err := w.watchDir(dir, fi.Name())
				if err != nil {
					return err
				}
			} else {
				if !w.since.IsZero() && fi.ModTime().After(w.since) {
					w.events <- filepath.Join(dir, fi.Name())
				}
			}
		}
	}
	return nil
}
