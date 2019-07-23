package pool

import (
	"fmt"
	"testing"
	"time"
)

func Test_a(t *testing.T) {
	f := func() interface{} {
		fmt.Println("f running")
		return time.Now().Unix()
	}
	p := NewPool(2, f)

	for i := 0; i < 12; i++ {
		item, err := p.Get()
		fmt.Println(item, err)
		//p.Put( item )
	}

}
