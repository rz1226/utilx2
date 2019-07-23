package counter

import (
	"strconv"
	"sync/atomic"
)

type Counter struct {
	value int64
}

func NewCounter() *Counter {
	return new(Counter)
}

func (c *Counter) Add(i int64) {
	atomic.AddInt64(&c.value, i)
}

func (c *Counter) GetCount() int64 {
	return c.value
}

func (c *Counter) Get() int64 {
	return c.value
}

func (c *Counter) Str() string {
	return strconv.FormatInt(c.value, 10)
}
