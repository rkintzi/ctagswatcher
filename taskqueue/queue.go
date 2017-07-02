package taskqueue

import "sync"

type Task interface {
	Run(errs chan<- error)
}

type TaskFunc func(errs chan<- error)

func (f TaskFunc) Run(errs chan<- error) {
	f(errs)
}

type Q struct {
	c   *sync.Cond
	ts  []Task
	wg  sync.WaitGroup
	fin bool
}

func New() *Q {
	return &Q{c: sync.NewCond(&sync.Mutex{})}
}

func (q *Q) Enqueue(t Task) {
	q.c.L.Lock()
	q.ts = append(q.ts, t)
	q.c.Signal()
	q.c.L.Unlock()
}

func (q *Q) Run(errs chan<- error) {
	var t Task
	q.wg.Add(1)
	for {
		q.c.L.Lock()
		for len(q.ts) == 0 && !q.fin {
			q.c.Wait()
		}
		if q.fin {
			q.wg.Done()
			return
		}
		t, q.ts = q.ts[0], q.ts[1:len(q.ts)]
		q.c.L.Unlock()
		t.Run(errs)
	}
}

func (q *Q) Stop() {
	q.fin = true
	q.c.Broadcast()
	q.wg.Wait()
}
