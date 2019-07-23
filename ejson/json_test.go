package ejson

import (
	"fmt"
	"testing"
)

func xTest_a(t *testing.T) {

	str := `[
    "and",
    {
        "name": "品牌",
        "op": "eq",
        "value": 12,
        "logs":[120,130,40],
        "haha":false
    }
]`

	e, _ := NewEjson([]byte(str))

	fmt.Println(str)
	fmt.Println(e)
	fmt.Println(e.IsMap())
	fmt.Println(e.IsArray())
	fmt.Println(e.ArrayGetString(0))
	fmt.Println(e.ArrayGetMap(1))
	m, _ := e.ArrayGetMap(1)
	fmt.Println(m.MapGetString("name"))
	fmt.Println(m.MapGetBool("haha"))

	m2, _ := e.ArrayGetMap(1)
	fmt.Println(m2.MapGetFloat64("value"))

	a1, _ := m.MapGetArray("logs")
	fmt.Println(a1.ArrayGetFloat64(1))

	fmt.Println("len:", e.Len())

	str2 := `{
        "name": "品牌",
        "op": "eq",
        "value": 12,
        "logs":[120,130,40]
    }`

	str3 := `{
        "name3": "品牌",
        "op3": "eq",
        "value3": 12,
        "logs3":[120,130,40]
    }`

	e2, _ := NewEjson([]byte(str2))
	e3, _ := NewEjson([]byte(str3))
	fmt.Println(e3)
	fmt.Println(e2.Merge(e3))

	m4 := make(map[string]interface{})
	fmt.Println(NewEjson(m4))

}

func xTest_axx(t *testing.T) {

	str := `[
    "and",
    {
        "name": "品牌",
        "op": "eq",
        "value": 12,
        "logs":[120,130,40],
        "haha":false
    }
]`

	e, err := NewEjson(str)
	fmt.Println("err:", err)
	fmt.Println(e.GetValue("array", "0", "string"))
	fmt.Println(e.GetValue("array", "1", "map"))
	fmt.Println(e.GetValue("array", "1", "map", "name", "string"))
	fmt.Println(e.GetValue("array", "1", "map", "haha", "bool"))
	fmt.Println(e.GetValue("array", "1", "map", "logs", "array"))
	fmt.Println(e.GetValue("array", "1", "map", "logs", "array", "1", "float64"))
	fmt.Println(e.GetValue("array", "1", "array"))
	fmt.Println(e.GetValue("map"))
}

func xTest_axxx(t *testing.T) {

	str := `{"a":[
    "and",
    {
        "name": "品牌",
        "op": "eq",
        "value": 12,
        "logs":[120,130,40],
        "haha":false
    }
]}`

	e, err := NewEjson(str)
	fmt.Println("err:", err)
	err2 := e.SetValue("ffffffffff", "map", "a", "array", "1", "array", "op")
	fmt.Println(err2)

	fmt.Println(e, err2)

	err2 = e.SetValue("sdfsdfsdf", "map", "a", "array", "1")
	fmt.Println(e, err2)

}
func Test_axxx33(t *testing.T) {

	str := `{}`

	e, err := NewEjson(str)
	fmt.Println("err:", err)
	err2 := e.SetValue(make(map[string]interface{}), "map", "a")
	fmt.Println(err2)

	fmt.Println(e, err2)

	err2 = e.SetValue(make([]interface{}, 3), "map", "a", "map", "c")
	fmt.Println(e, err2)

	err2 = e.SetValue("sdfsdfsdf", "map", "a", "map", "c", "array", "2")
	fmt.Println(e, err2)

}

func Test_axxx333(t *testing.T) {

	str := `[null,null,null]`

	e, err := NewEjson(str)
	fmt.Println("err:", err)
	err2 := e.SetValue(make(map[string]interface{}), "array", "1")
	fmt.Println(err2)

	fmt.Println(e, err2)

	err2 = e.SetValue(make([]interface{}, 3), "array", "1", "map", "c")
	fmt.Println(e, err2)

	err2 = e.SetValue("sdfsdfsdf", "array", "1", "map", "c", "array", "2")
	fmt.Println(e, err2)

}
