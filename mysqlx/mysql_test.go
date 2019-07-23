package mysqlx

import (
	"fmt"
	"testing"
)

func Test_query(t *testing.T) {

	conf := "root:123456@tcp(127.0.0.1:3306)/task?charset=utf8"
	p := NewDbPool(conf, 3)
	result, err := p.Query("select * from screen_user_new limit 100", nil)

	fmt.Println(result.Len(), err)
}

func BenchmarkSliceAppend(b *testing.B) {
	b.N = 1000
	a := make([]string, 0, b.N)
	for i := 0; i < b.N; i++ {
		a = append(a, "sdfsdfsdsgggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggg")
	}
}

func BenchmarkSliceAppend2(b *testing.B) {
	b.N = 1000
	a := make([]string, 0)
	for i := 0; i < b.N; i++ {
		a = append(a, "sdfsdfsdsgggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggg")
	}
}
func BenchmarkSliceAppend3(b *testing.B) {
	b.N = 1000
	a := make([]string, 0, 1000)
	for i := 0; i < b.N; i++ {
		a = append(a, "sdfsdfsdsgggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggg")
	}
}
