package counter

import (
	"fmt"
	"testing"
)

func Test_all(t *testing.T) {

	c := NewCounter()
	c.Add(33)
	fmt.Println(c.Count())
	if c.Count() != 33 {
		t.Error("err")
	}

	c.Add(-44)
	fmt.Println(c.Count())
	if c.Count() != -11 {
		t.Error("err")
	}

	c.Add(44)
	fmt.Println(c.Count())
	if c.Count() != 33 {
		t.Error("err")
	}

	c.Add(44)
	fmt.Println(c.Count())
	if c.Count() != 77 {
		t.Error("err")
	}

}
