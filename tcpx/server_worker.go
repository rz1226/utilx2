package tcpx

import (
	"errors"
	"fmt"
	"net"
	"runtime"
	"sync"
	"time"
)

//--------------------workerpool----------------/
type WorkerPool struct {
	Workers     []*Worker
	IdleWorkers []*Worker
	WorkFunc    func(*net.TCPConn) error
	MaxCount    int
	LogFunc     func(string)
	mu          *sync.Mutex
}

func NewWorkerPool(f func(*net.TCPConn) error, max int) *WorkerPool {
	wp := &WorkerPool{}
	wp.MaxCount = max
	wp.LogFunc = func(log string) {
		fmt.Println("log:", log)
	}
	wp.WorkFunc = f
	wp.Workers = make([]*Worker, 0, wp.MaxCount)
	wp.IdleWorkers = make([]*Worker, 0, wp.MaxCount)
	wp.mu = new(sync.Mutex)
	return wp
}

func (wp *WorkerPool) GetWorker() (*Worker, error) {
	if len(wp.IdleWorkers) == 0 {
		err := wp.newWorker()
		if err != nil {
			return nil, err
		}
	}

	wp.mu.Lock()
	defer wp.mu.Unlock()
	n := len(wp.IdleWorkers) - 1
	r := wp.IdleWorkers[n]
	wp.IdleWorkers = wp.IdleWorkers[:n]
	return r, nil
}

func (wp *WorkerPool) newWorker() error {
	if wp.MaxCount <= len(wp.Workers) {
		return errors.New("workerpool full ")
	}
	w := &Worker{}
	w.IsStop = true
	w.IsIdle = true
	w.IdleTime = int(time.Now().Unix())
	w.WorkChan = make(chan *net.TCPConn, 1)
	w.Pool = wp
	wp.mu.Lock()
	defer wp.mu.Unlock()
	wp.Workers = append(wp.Workers, w)
	wp.IdleWorkers = append(wp.IdleWorkers, w)
	return nil
}

//--------------------worker----------------/

type Worker struct {
	IsStop   bool
	IsIdle   bool
	WorkChan chan *net.TCPConn
	IdleTime int
	Pool     *WorkerPool
}

func (w *Worker) putBack() {
	w.Pool.mu.Lock()
	defer w.Pool.mu.Unlock()
	w.Pool.IdleWorkers = append(w.Pool.IdleWorkers, w)
}

func (w *Worker) Prepare(con *net.TCPConn) {
	w.WorkChan <- con
}

func (w *Worker) Start() {
	fmt.Println("works count:", len(w.Pool.Workers))
	fmt.Println("idle works count:", len(w.Pool.IdleWorkers))
	fmt.Println(runtime.NumGoroutine())
	//不能重复启动
	if w.IsStop == false {
		return
	}

	go func() {
		w.IsStop = false
		w.IdleTime = 0
		var err error
		defer func() {
			w.IsStop = true
			w.IdleTime = int(time.Now().Unix())
			if len(w.WorkChan) > 0 {
				fmt.Println("fatal : this should not happen")
				w.Pool.LogFunc("fatal: this should not happen")
				w.WorkChan = make(chan *net.TCPConn, 1)
			}
		}()
		for {
			err = w.Pool.WorkFunc(<-w.WorkChan)
			if err != nil {
				w.Pool.LogFunc(err.Error())
			}
			w.putBack()
		}
	}()

}
