package pool

import (
	"container/heap"
	"errors"
	"fmt"
	"sync"
	"time"
)

//普通的pool 用法类似sync.Pool，区别是非临时对象池,以及有容量的概念,以及用容量限制了新创建对象的数量

type Pool struct {
	New   func() interface{}
	Heap  *PoolHeap
	Cap   int
	Count int
	mu    *sync.Mutex
}

func NewPool(cap int, f func() interface{}) *Pool {
	p := new(Pool)
	h := &PoolHeap{}
	heap.Init(h)
	p.Heap = h
	p.mu = new(sync.Mutex)
	p.New = f
	p.Cap = cap
	p.Count = 0
	return p
}

type Element struct {
	Item    interface{}
	PutTime int64
}

func (p *Pool) Len() int {
	return p.Heap.Len()
}

//如果pool满了
func (p *Pool) Get() (interface{}, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	len := p.Heap.Len()
	if len > 0 {
		ele := p.Heap.Pop().(*Element)
		return ele.Item, nil
	}
	// pool is full
	if p.Cap <= p.Count {
		return nil, errors.New("pool :can not make more items")
	}
	p.Count++
	return p.New(), nil
}

func (p *Pool) Put(x interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()
	len := p.Heap.Len()
	if len >= p.Cap {
		fmt.Println("this should not happen")
		return
	}
	ele := new(Element)
	ele.Item = x
	ele.PutTime = time.Now().Unix()
	heap.Push(p.Heap, ele)
}
