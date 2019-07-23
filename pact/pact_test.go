package pact

import (
	"bytes"
	"fmt"
	"testing"
)

func Test_all(t *testing.T) {
	a := new(bytes.Buffer)
	data := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}

	fmt.Println(Write(a, data))
	fmt.Println(Write(a, []byte{}))
	fmt.Println(Write(a, data))
	fmt.Println(Read(a))
	fmt.Println(Read(a))
	fmt.Println(Read(a))

	fmt.Println(Make(9, []byte{10, 11}))
	fmt.Println(Get([]byte{9, 10, 11}))

}
