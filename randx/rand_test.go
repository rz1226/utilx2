package randx

import (
	"fmt"
	"testing"
)

func Test_all(t *testing.T) {
	data := GetRandomString(33)
	fmt.Println(data)
	if len(data) != 33 {

		t.Fatal("err")
	}

	n := GetRandomInt(34)
	fmt.Println(n)

	data2 := GetRandomString(GetRandomInt(45))
	fmt.Println(data2)
}
