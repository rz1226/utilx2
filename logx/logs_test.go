package logx

import (
	"fmt"
	"testing"
)

func Test_logs(t *testing.T) {
	logs := NewLogs("a", "c")
	logs.Put("a", "log message a 1")
	logs.Put("a", "log message a 2")
	logs.Put("b", "log message b 1")
	logs.Put("b", "log message b 2")

	fmt.Println(logs.Show("a", 40))
	fmt.Println(logs.Show("b", 40))
}
func BenchmarkLogs(b *testing.B) {
	logs := NewLogs("a")
	for i := 0; i < b.N; i++ {
		logs.Put("a", "log message a 1", "bb", "cc", "fuck")
	}
	fmt.Println(logs.Show("a", 30))
}
