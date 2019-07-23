package onoff

import (
	"sync/atomic"
)

type Gate struct {
	Value uint32
}

func NewGate() *Gate {
	g := &Gate{}
	atomic.StoreUint32(&g.Value, 1)
	return g
}

func (g *Gate) Open() {
	atomic.StoreUint32(&g.Value, 1)
}

func (g *Gate) Close() {
	atomic.StoreUint32(&g.Value, 0)
}

func (g *Gate) IsOpen() bool {
	return atomic.LoadUint32(&g.Value) == 1
}
