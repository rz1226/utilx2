package counter

import (
	"fmt"
	"testing"
)

func Test_counters(t *testing.T) {

	for i := 0; i < 10000; i++ {
		Inc("a")
	}

	for i := 0; i < 10022; i++ {
		Inc("b")
	}

	fmt.Println(Get("a"))
	fmt.Println(Get("b"))

}

func BenchmarkCounters(b *testing.B) {

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Inc("ax")
		}
	})

	b.Log(Get("ax"))
}
