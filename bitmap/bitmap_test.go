package bitmap

import (
	"fmt"
	"testing"
)

func Test_bitmap(t *testing.T) {
	data := make([]byte, 12)
	b := NewBitMap(data)

	b.SetPositionTrue(120)
	fmt.Println(b.GetPositionValue(120))
	fmt.Println(b.Len())
}

func Benchmark_bitmap(b *testing.B) {
	data := make([]byte, 12)
	bitmap := NewBitMap(data)
	for i := 0; i < b.N; i++ {
		bitmap.SetPositionTrue(i)
		if bitmap.GetPositionValue(i) == false {
			b.Log("err ")
		}
	}
	fmt.Println("len", bitmap.Len())
	fmt.Println("counter", bitmap.DirtyCount())

}
