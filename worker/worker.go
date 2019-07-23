package worker

import (
	//"fmt"
	"github.com/xiaobai22/utilx2/pool"
)

/*
	workers := NewWorkers( 12, f func( interface{} ){} )
	worker := workers.Get()
	worker.C <- "task"


*/

type Workers struct {
	Pool *pool.Pool
}

func NewWorkers(cap int, f func(interface{})) *Workers {
	workers := &Workers{}
	fn := func() interface{} {
		return newWorker(workers, f)
	}
	p := pool.NewPool(cap, fn)
	workers.Pool = p
	return workers
}

type Worker struct {
	C     chan interface{}
	Group *Workers
	F     func(interface{})
}

func (workers *Workers) Len() int {
	return workers.Pool.Len()
}

func (workers *Workers) Get() (*Worker, error) {
	item, err := workers.Pool.Get()
	if err != nil {
		return nil, err
	}
	return item.(*Worker), nil
}

func newWorker(group *Workers, f func(interface{})) *Worker {
	w := &Worker{}
	w.C = make(chan interface{}, 0)
	w.Group = group
	w.F = f
	go w.Start()
	return w
}

func (w *Worker) Start() {
	for v := range w.C {
		w.F(v)
		w.PutBack()
	}
}
func (w *Worker) PutBack() {
	w.Group.Pool.Put(w)
}
