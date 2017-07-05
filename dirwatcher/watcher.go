package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
)

var ErrNotStarted = fmt.Errorf("Watcher not started")

type Watcher struct {
	events chan<- string
	errors chan<- error
	done   chan struct{}
	w      *fsnotify.Watcher
	wg     sync.WaitGroup
	dir    string
}

func New(dir string, events chan<- string, errors chan<- error) *Watcher {
	return &Watcher{dir: dir, events: events, errors: errors, done: make(chan struct{})}
}

func (w *Watcher) Start() error {
	var err error
	w.w, err = fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	w.wg.Add(1)
	go w.monitor()
	return w.walkdir(w.dir)
}

func (w *Watcher) Stop() error {
	if w.w == nil {
		return ErrNotStarted
	}
	err := w.w.Close()
	close(w.done)
	w.wg.Wait()
	return err
}

func (w *Watcher) walkdir(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	werr := w.w.Add(dir)
	if werr != nil {
		return werr
	}
	for {
		fis, err := d.Readdir(100)
		for _, fi := range fis {
			if !fi.IsDir() {
				continue
			}
			dirpath := filepath.Join(dir, fi.Name())
			err := w.walkdir(dirpath)
			if err != nil {
				return err
			}
		}
		if err != nil && err != io.EOF {
			return err
		} else if err == io.EOF {
			return nil
		}
	}
}

func (w *Watcher) monitor() {
	defer w.wg.Done()
	for {
		select {
		case ev, ok := <-w.w.Events:
			if !ok {
				return
			}
			w.enqueue(ev.Name)
		case err, ok := <-w.w.Errors:
			if !ok {
				return
			}
			select {
			case w.errors <- err:
			}
		}
	}
}

func (w *Watcher) Enqueue(filename string) {

}

func main() {
	fmt.Println("vim-go")
}
