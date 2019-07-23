package numlimit

import (
	"fmt"
	"log"
	"sync/atomic"
)

/*
var x := NewNumLimiter( 5 )
x.Add()  //超过限制返回false ,否则返回true
x.Done()

*/

type NumLimiter struct {
	maxValue   uint32
	currentValue uint32
}

func NewNumLimiter(max uint32) *NumLimiter {
	n := &NumLimiter{}
	n.maxValue = max
	n.currentValue = 0
	return n
}

func (n *NumLimiter) Done() {
	if n.currentValue == 0 {
		log.Println("numlimit ,this should not happen")
		return
	}
	atomic.AddUint32(&n.currentValue, ^uint32(0))
}

//并发的问题可以容忍
func (n *NumLimiter) Add() bool {

	if n.currentValue >= n.maxValue {
		return false
	} else {
		atomic.AddUint32(&n.currentValue, 1)
		return true
	}

}
func (n *NumLimiter) Show() {
	fmt.Println(n.currentValue)
}
