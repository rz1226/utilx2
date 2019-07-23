package show

import (
	"fmt"
)

//用十六进制打印[]byte
func ShowBytes(data []byte) {
	fmt.Print("[")
	for _, v := range data {
		fmt.Printf(" %2x ", v)
	}
	fmt.Print("]")
	fmt.Println("")

}
