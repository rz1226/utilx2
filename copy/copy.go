package copy

import (
	"fmt"
	"reflect"
)

//copy data from slice to array

func CopySlice2Arr(data interface{}, arr interface{}) error {
	defer func() {
		if pa := recover(); pa != nil {
			fmt.Println(pa)
			fmt.Println(" in CopySlice2Arr , error : arr not a ptr?")
		}

	}()

	reflect.Copy(reflect.ValueOf(arr).Elem(), reflect.ValueOf(data))
	return nil
}

func CopyArr2Slice(data interface{}, arr interface{}) error {
	defer func() {
		if pa := recover(); pa != nil {
			fmt.Println(pa)
			fmt.Println(" in CopyArr2Slice , error : slice not a ptr?")
		}

	}()

	reflect.Copy(reflect.ValueOf(data).Elem(), reflect.ValueOf(arr))
	return nil
}
