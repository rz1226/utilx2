package copy

import (
	"fmt"
	"testing"
)

func Test_all(t *testing.T) {

	test := []byte{1, 2, 3, 4, 5, 6}
	arr := [7]byte{}
	CopySlice2Arr(test, &arr)

	fmt.Println(arr)

}
