package ratelimit

import (
	"fmt"
	"sync/atomic"
	"time"
)

//每秒钟限制次数

//  s := NewRateLimiter(12) ok := s.Get()  不阻塞
/*


func main(){
	r := tool.NewRateLimiter(100000)
	for i := 0; i < 100000; i++ {
		//time.Sleep(time.Millisecond * 80 )
		fmt.Println( r.Get() )
		fmt.Println(i)
	}

	time.Sleep( time.Second * 1000 )
}


*/

type RateLimiter struct {
	currentCount uint32
	maxCount     uint32
	ticker       *time.Ticker
}

func NewRateLimiter(max uint32) *RateLimiter {
	r := &RateLimiter{}
	r.maxCount = max
	r.currentCount = 1
	r.ticker = time.NewTicker(time.Second * 1)
	r.refresh()
	return r
}

func (r *RateLimiter) Get() bool {
	//func AddUint32(addr *uint32, delta uint32) (new uint32)
	//func CompareAndSwapInt32(addr *int32, old, new int32) (swapped bool)
	//func SwapUint32(addr *uint32, new uint32) (old uint32)
	currentCount := atomic.SwapUint32(&r.currentCount, r.currentCount+1) //无法保证r.count+1这里的取值是及时的,但可以容忍
	fmt.Println("get now value count:", currentCount)
	if currentCount > r.maxCount {
		return false
	} else {
		return true
	}
}

func (r *RateLimiter) refresh() {
	go func() {
		for {
			<-r.ticker.C
			//fmt.Println("Ticking", r.count)
			atomic.StoreUint32(&r.currentCount, 0)
		}
	}()

}
