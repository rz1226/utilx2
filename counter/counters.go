package counter

type Counters struct {
	data map[string]*Counter
}

func NewCounters(strs ...string) *Counters {
	c := &Counters{}
	c.data = make(map[string]*Counter, 0)
	for _, v := range strs {
		c.data[v] = NewCounter()
	}
	return c
}
func (c *Counters) IncBy(name string, num int64) {
	counter, ok := c.data[name]
	if !ok {
		return
	}
	counter.Add(num)
}

func (c *Counters) Inc(name string) {
	counter, ok := c.data[name]
	if !ok {
		return
	}
	counter.Add(1)
}
func (c *Counters) Get(name string) int64 {
	obj, ok := c.data[name]
	if !ok {
		return 0
	}
	return obj.Get()
}

func (c *Counters) Str(name string) string {
	obj, ok := c.data[name]
	if !ok {
		return "没有计数器\n"
	}
	return obj.Str()
}

