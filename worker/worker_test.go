package worker

import (
	"fmt"
	"testing"
	"time"
	"utilx/console"
)

func Test_all(t *testing.T) {

	fn := func(data interface{}) {
		a := data.(int)
		//time.Sleep( time.Second*1)
		fmt.Println(a)
	}

	workers := NewWorkers(2, fn)

	for i := 0; i < 12; i++ {
		time.Sleep(time.Millisecond * 1)
		w, err := workers.Get()
		if err != nil {
			fmt.Println(err)
		} else {
			w.C <- i

		}

	}

	time.Sleep(time.Second * 3)
	console.Blue(workers.Len())

}
