package counter

import (
	"fmt"
	"testing"
)

func Test_all(t *testing.T) {

	c := NewCounter()
	c.Add(33)
	fmt.Println(c.GetCount())
	if c.GetCount() != 33 {
		t.Error("err")
	}

	c.Add(-44)
	fmt.Println(c.GetCount())
	if c.GetCount() != -11 {
		t.Error("err")
	}

	c.Add(44)
	fmt.Println(c.GetCount())
	if c.GetCount() != 33 {
		t.Error("err")
	}

	c.Add(44)
	fmt.Println(c.GetCount())
	if c.GetCount() != 77 {
		t.Error("err")
	}

	c.Clear()
	if c.GetCount() != 0 {
		t.Error("err")
	}

}
