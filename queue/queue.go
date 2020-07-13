package queue

import (
	"errors"
	"time"
)

/****************************基于chan的队列,可批量取出数据,有超时设定**************************************************************************/
type BatchQueue struct {
	C              chan interface{}
	TimeoutSeconds int
}

func NewBatchQueue(len, seconds int) *BatchQueue {
	b := &BatchQueue{}
	b.C = make(chan interface{}, len)
	b.TimeoutSeconds = seconds
	return b
}
func (b *BatchQueue) Len() int {
	return len(b.C)
}
func (b *BatchQueue) Close()   {
	close( b.C )
}
func (b *BatchQueue) Put(data interface{}) {
	b.C <- data
}

func (b *BatchQueue) PutWithError(data interface{}) error {
	select {
	case b.C <- data:
		return nil
	default:
		return errors.New("full queue")
	}

}

//第三个参数代表chan是否可用，关闭后不可用
func (b *BatchQueue) Get(num int) (interface{}, int, bool ) {
	result := make([]interface{}, 0, num)
	timeoutSeconds := b.TimeoutSeconds
	timer := time.After(time.Second * time.Duration(timeoutSeconds))
	counter := 0
	for i := 0; i < num; i++ {
		select {
		case d , ok  := <-b.C:
			if !ok{
				if counter == 0 {
					return result ,counter, false
				}else{
					return result ,counter, true
				}

			}else{
				result = append( result , d )
				counter++
			}

		case <-timer:
			if counter > 0 {
				return result[:counter], counter, true
			}
			return result, 0, true
		}

	}
	return result, counter, true

}
