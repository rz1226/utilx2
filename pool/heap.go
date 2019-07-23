package pool

import (

//"fmt"
)

type PoolHeap []*Element

func (h PoolHeap) Len() int           { return len(h) }
func (h PoolHeap) Less(i, j int) bool { return h[i].PutTime > h[j].PutTime } //反过来排序
func (h PoolHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *PoolHeap) Push(x interface{}) {
	*h = append(*h, x.(*Element))
}

func (h *PoolHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
