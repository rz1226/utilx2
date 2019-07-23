package circleq

import (
	"sync/atomic"
)

//只保留最近写入的部分数据
type CQ struct {
	seqNum   uint64
	size     uint32
	dataList []ele
}

type ele struct {
	seqNum       uint64
	contentValue interface{}
}

func NewCQ(size uint32) *CQ {
	c := &CQ{}
	c.seqNum = 0
	c.size = minQuantity(size)
	c.dataList = make([]ele, c.size)
	ele := &(c.dataList[0])
	ele.seqNum = 0
	return c
}

func (c *CQ) Put(val interface{}) {
	//func AddUint64(addr *uint64, delta uint64) (new uint64)
	nextseqNum := atomic.AddUint64(&c.seqNum, 1)
	positionInList := nextseqNum & uint64((c.size - 1))
	ele := &(c.dataList[positionInList])
	ele.seqNum = nextseqNum
	ele.contentValue = val
}

func (c *CQ) GetSeveral(count int) ([]interface{}, uint64) {
	res := make([]interface{}, count)
	newestSeqNum := atomic.LoadUint64(&c.seqNum)
	for i := 0; i < count-1; i++ {
		pos := (newestSeqNum - uint64(i)) & uint64((c.size - 1))
		ele := &(c.dataList[pos])
		if ele.seqNum == newestSeqNum-uint64(i) {
			res[i] = ele.contentValue
		} else {
			break
		}
	}
	return res, newestSeqNum
}

// round 到最近的2的倍数
func minQuantity(v uint32) uint32 {
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v++
	return v
}
