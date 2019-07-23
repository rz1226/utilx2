package ejson

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

type Ejson struct {
	j          interface{} //json对象，一般是[]interface{}或者map[string]interface{}
	sReplicate string      //json字符串， j的数据的字符串形式，内容是相同的
}

func NewEjson(jsonData interface{}) (*Ejson, error) {
	e := Ejson{}
	//如果传过来的参数是个json字符串
	if jsonString, ok := jsonData.(string); ok {
		err := json.Unmarshal([]byte(jsonString), &e.j)
		if err != nil {
			fmt.Println("newejson:无法解析成json", err, jsonString)
			return nil, err
		}
		e.sReplicate = jsonString
		return &e, nil
	}
	//如果传过来的参数是byte数组，其实还是字符串
	if jsonBytes, ok := jsonData.([]byte); ok {
		err := json.Unmarshal(jsonBytes, &e.j)
		if err != nil {
			fmt.Println("newejson:无法解析成json", err, jsonBytes)
			return nil, err
		}
		e.sReplicate = string(jsonBytes)
		return &e, nil
	}

	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		return nil, errors.New("newjson:传来的json无法解析成字符串")
	}
	e.j = jsonData
	e.sReplicate = string(jsonBytes)
	return &e, nil
}

//重新计算sreplicate
func (e *Ejson) RenewSReplicate() {
	jsonBytes, _ := json.Marshal(e.j)
	e.sReplicate = string(jsonBytes)
}

func (e *Ejson) Json() interface{} {
	return e.j
}

func (e *Ejson) Bytes() []byte {
	return []byte(e.sReplicate)
}
func (e *Ejson) String() string {
	return e.sReplicate
}

//设置内部的值
//where =( "map", "taginfo", "map", "expire", "float64") 第一个参数是自己本身的类型
//支持的类型 array map string float64 bool
//如果最后两个参数指明的键和类型已经存在，那么修改，否则添加
//如果需要修改或者添加的值所在的数据不存在，不存在则报错
func (e *Ejson) Update(value interface{}, where ...string) error {

	v := reflect.ValueOf(e.j)
	err := updateJson(&v, value, where...)
	e.RenewSReplicate()
	if err != nil {
		return err
	}
	return nil

}

//第一个参数是要修改的数据的reflect.value类型
func updateJson(haystack *reflect.Value, value interface{}, where ...string) error {
	lengthWhere := len(where)
	if lengthWhere < 2 {
		//这个限制表明，无法整体替换数据
		return errors.New("SetValue:参数长度至少2")
	}

	jsonTypeWhere := where[0]
	keyOrIndex := where[1]

	jtype := (*haystack).Type().Kind().String()
	if jtype == "map" && jsonTypeWhere == "map" {

	} else if jtype == "slice" && jsonTypeWhere == "array" {

	} else {
		return errors.New("setvalue:数据本身不是指定的map或者array类型")
	}

	if lengthWhere == 2 {
		setvalue := value
		if jtype == "map" {
			(*haystack).SetMapIndex(reflect.ValueOf(keyOrIndex), reflect.ValueOf(setvalue))
			return nil
		} else if jtype == "slice" {
			index, err := strconv.Atoi(keyOrIndex)
			if err != nil {
				return errors.New("SetValue:数组的键必须是数字")
			}
			if index > ((*haystack).Len() - 1) {
				return errors.New("setvalue:数组越界:" + keyOrIndex)
			}
			(*haystack).Index(index).Set(reflect.ValueOf(setvalue))
			return nil
		}

	}
	unProcessedwhere := where[2:]
	if jtype == "map" {
		nextLevelHayStack := (*haystack).MapIndex(reflect.ValueOf(keyOrIndex)).Elem()
		return updateJson(&nextLevelHayStack, value, unProcessedwhere...)
	} else if jtype == "slice" {
		index, err := strconv.Atoi(keyOrIndex)
		if err != nil {
			return errors.New("SetValue:数组的键必须是数字")
		}
		if index > ((*haystack).Len() - 1) {
			return errors.New("setvalue:数组越界:" + keyOrIndex)
		}
		nextLevelHayStack := (*haystack).Index(index).Elem()
		return updateJson(&nextLevelHayStack, value, unProcessedwhere...)
	}

	return errors.New("SetValue:仅支持对map,array类型的修改")

}

//获取内部的值
//where =( "map", taginfo", "map", "expire", "float64") 第一个参数是自己本身的类型
//支持的类型 array map string float64 bool
func (e *Ejson) GetInnerValue(where ...string) (interface{}, error) {
	lengthWhere := len(where)
	if lengthWhere == 1 {
		if where[0] == "map" {
			mapData, ok := e.j.(map[string]interface{})
			if ok {
				return mapData, nil
			}
			return nil, errors.New("GetValue:不是map")
		}
		if where[0] == "array" {
			arrData, ok := e.j.([]interface{})
			if ok {
				return arrData, nil
			}
			return nil, errors.New("GetValue:不是array")
		}
		return nil, errors.New("GetValue:复合类型只有map,array")
	}

	if lengthWhere%2 != 1 {
		return nil, errors.New("参数的个数必须是奇数")
	}
	jsonTypeWhere := where[0]
	if jsonTypeWhere == "map" {
		if !e.IsMap() {
			return nil, errors.New("GetValue:数据自身不是map")
		}
	} else if jsonTypeWhere == "array" {
		if !e.IsArray() {
			return nil, errors.New("GetValue:数据自身不是array")
		}
	} else {
		return nil, errors.New("GetValue:此操作仅支持类型是map和array")
	}
	headerWhere := where[:2]
	footerWhere := where[2:]
	k := headerWhere[1]
	t := footerWhere[0]
	if jsonTypeWhere == "map" {
		if t == "map" {
			ej, err := e.MapGetMap(k)
			if err != nil {
				return nil, err
			}
			return ej.GetInnerValue(footerWhere...)

		} else if t == "array" {
			ej, err := e.MapGetArray(k)
			if err != nil {
				return nil, err
			}
			return ej.GetInnerValue(footerWhere...)
		} else if t == "string" {
			ej, err := e.MapGetString(k)
			if err != nil {
				return nil, err
			}
			return ej, nil
		} else if t == "float64" {
			ej, err := e.MapGetFloat64(k)
			if err != nil {
				return nil, err
			}
			return ej, nil
		} else if t == "bool" {
			ej, err := e.MapGetBool(k)
			if err != nil {
				return nil, err
			}
			return ej, nil
		}
	} else if jsonTypeWhere == "array" {
		index, err := strconv.Atoi(k)
		if err != nil {
			return nil, errors.New("GetValue:从数组取数据的key必须可以转换为int")
		}
		if t == "map" {
			ej, err := e.ArrayGetMap(index)
			if err != nil {
				return nil, err
			}
			return ej.GetInnerValue(footerWhere...)
		} else if t == "array" {
			ej, err := e.ArrayGetArray(index)
			if err != nil {
				return nil, err
			}
			return ej.GetInnerValue(footerWhere...)
		} else if t == "string" {
			ej, err := e.ArrayGetString(index)
			if err != nil {
				return nil, err
			}
			return ej, nil
		} else if t == "float64" {
			ej, err := e.ArrayGetFloat64(index)
			if err != nil {
				return nil, err
			}
			return ej, nil
		} else if t == "bool" {
			ej, err := e.ArrayGetBool(index)
			if err != nil {
				return nil, err
			}
			return ej, nil
		}
	} else {
		return nil, errors.New("GetValue: 类型错误，复合类型只有array和map")
	}
	return nil, errors.New("GetValue:未知错误，或类型不支持,支持map,array,float64,string,bool")

}

//拷贝map类型的json , 就是新创建，用s解析填充
func (e *Ejson) GetDataInNewMap() (map[string]interface{}, error) {
	if e.IsMap() {
		var result map[string]interface{}
		err := json.Unmarshal([]byte(e.sReplicate), &result)
		if err != nil {
			return nil, errors.New("ejson:copymap: there is error :" + err.Error())
		}
		return result, nil
	}
	return nil, errors.New("ejson:copymap: this is not map")
}

//如果不是map类型返回空值
func (e *Ejson) GetMapData() map[string]interface{} {
	data, ok := e.j.(map[string]interface{})
	if ok {
		return data
	}
	return make(map[string]interface{})
}

//如果不是数组类型返回空值
func (e *Ejson) GetArrayData() []interface{} {
	data, ok := e.j.([]interface{})
	if ok {
		return data
	}
	return make([]interface{}, 1)
}

func (e *Ejson) MergeMap(a *Ejson) (*Ejson, error) {
	if !e.IsMap() || !a.IsMap() {
		return nil, errors.New("ejson:merge 两个值必须都是map")
	}
	res := &Ejson{j: make(map[string]interface{})}

	for k, v := range e.j.(map[string]interface{}) {
		res.j.(map[string]interface{})[k] = v
	}
	for k2, v2 := range a.j.(map[string]interface{}) {
		res.j.(map[string]interface{})[k2] = v2
	}
	return res, nil
}

func (e *Ejson) Len() int {
	if e.IsMap() {
		return len(e.j.(map[string]interface{}))
	}
	if e.IsArray() {
		return len(e.j.([]interface{}))
	}
	return 0
}

func (e *Ejson) ArrayGetData(key int) (*Ejson, error) {
	json, ok := e.j.([]interface{})
	if !ok {
		return nil, errors.New("ejson:ArrayGetData主体对象必须是array")
	}
	if key > (len(json) - 1) {
		return nil, errors.New("ejson:ArrayGetData找不到index")
	}
	m := json[key]
	return NewEjson(m)

	//return &Ejson{j: m}, nil
}

func (e *Ejson) ArrayGetMap(key int) (*Ejson, error) {
	json, ok := e.j.([]interface{})
	if !ok {
		return nil, errors.New("ejson:arraygetmap主体对象必须是array")
	}
	if key > (len(json) - 1) {
		return nil, errors.New("ejson:arraygetmap找不到index")
	}
	m, ok := (json[key]).(map[string]interface{})
	if !ok {
		return nil, errors.New("ejson:arraygetmap找到的index不是map类型")
	}
	return NewEjson(m)
}

func (e *Ejson) ArrayGetArray(key int) (*Ejson, error) {
	json, ok := e.j.([]interface{})
	if !ok {
		return nil, errors.New("ejson:ArrayGetArray主体对象必须是array")
	}
	if key > (len(json) - 1) {
		return nil, errors.New("ejson:ArrayGetArray找不到index")
	}
	m, ok := (json[key]).([]interface{})
	if !ok {
		return nil, errors.New("ejson:ArrayGetArray找到的index不是array类型")
	}
	return NewEjson(m)
}

func (e *Ejson) MapGetData(key string) (interface{}, error) {
	json, ok := e.j.(map[string]interface{})
	if !ok {
		return nil, errors.New("ejson:MapGetData主体对象必须是map")
	}
	find, ok := json[key]
	if !ok {
		return nil, errors.New("ejson:MapGetData找不到index")
	}
	return find, nil
}

func (e *Ejson) MapGetMap(key string) (*Ejson, error) {
	json, ok := e.j.(map[string]interface{})
	if !ok {
		return nil, errors.New("ejson:MapGetMap主体对象必须是map")
	}
	find, ok := json[key]
	if !ok {
		return nil, errors.New("ejson:MapGetMap找不到index")
	}
	arr, ok := find.(map[string]interface{})
	if !ok {
		return nil, errors.New("ejson:MapGetMap找到的index不是map类型")
	}
	return NewEjson(arr)
}

func (e *Ejson) ArrayGetFloat64(key int) (float64, error) {
	json, ok := e.j.([]interface{})
	if !ok {
		return 0, errors.New("ejson:ArrayGetFloat64主体对象必须是数组")
	}
	if key > (len(json) - 1) {
		return 0, errors.New("ejson:ArrayGetFloat64找不到index")
	}
	f, ok := json[key].(float64)
	if !ok {
		return 0, errors.New("ejson:ArrayGetFloat64找到的值不是float64类型")
	}
	return f, nil
}

func (e *Ejson) ArrayGetBool(key int) (bool, error) {
	json, ok := e.j.([]interface{})
	if !ok {
		return false, errors.New("ejson:ArrayGetBool 主体对象必须是map")
	}
	if key > (len(json) - 1) {
		return false, errors.New("ejson:ArrayGetBool :key not exist")
	}
	value, ok := json[key].(bool)
	if !ok {
		return false, errors.New("ejson:ArrayGetBool :value not bool")
	}
	return value, nil

}

func (e *Ejson) ArrayGetString(key int) (string, error) {
	json, ok := e.j.([]interface{})
	if !ok {
		return "", errors.New("ejson:ArrayGetString主体对象必须是数组")
	}
	if key > (len(json) - 1) {
		return "", errors.New("ejson:ArrayGetString找不到index")
	}
	str, ok := json[key].(string)
	if !ok {
		return "", errors.New("ejson:ArrayGetString找到的值不是String类型")
	}
	return str, nil
}

func (e *Ejson) MapGetArray(key string) (*Ejson, error) {
	json, ok := e.j.(map[string]interface{})
	if !ok {
		return nil, errors.New("ejson:MapGetArray主体对象必须是map")
	}
	find, ok := json[key]

	if !ok {

		return nil, errors.New("ejson:MapGetArray找不到index")
	}
	arr, ok := find.([]interface{})
	if !ok {
		return nil, errors.New("ejson:MapGetArray找到的index不是array类型")
	}
	return NewEjson(arr)

}

func (e *Ejson) MapGetFloat64(key string) (float64, error) {
	json, ok := e.j.(map[string]interface{})
	if !ok {
		return 0, errors.New("ejson:MapGetFloat64主体对象必须是map")
	}
	find, ok := json[key]
	if !ok {
		return 0, errors.New("ejson:MapGetFloat64找不到index")
	}
	num, ok := find.(float64)
	if !ok {
		return 0, errors.New("ejson:MapGetFloat64找到的index不是float64")
	}
	return num, nil

}

func (e *Ejson) MapGetBool(key string) (bool, error) {
	json, ok := e.j.(map[string]interface{})
	if !ok {
		return false, errors.New("ejson:MapGetBool主体对象必须是map")
	}
	find, ok := json[key]
	if !ok {
		return false, errors.New("ejson:MapGetBool:key not exist")
	}
	value, ok := find.(bool)
	if !ok {
		return false, errors.New("ejson:MapGetBool:value not bool")
	}
	return value, nil

}

func (e *Ejson) MapGetString(key string) (string, error) {
	json, ok := e.j.(map[string]interface{})
	if !ok {
		return "", errors.New("ejson:MapGetString对象本身必须是map")
	}
	find, ok := json[key]
	if !ok {
		return "", errors.New("ejson:MapGetString找不到index")
	}
	str, ok := find.(string)
	if !ok {
		return "", errors.New("ejson:MapGetString找到的数据不是string")
	}
	return str, nil

}

func (e *Ejson) IsMap() bool {
	_, ok := e.j.(map[string]interface{})
	if ok {
		return true
	}
	return false
}

func (e *Ejson) IsArray() bool {
	_, ok := e.j.([]interface{})
	if ok {
		return true
	}
	return false
}

//其作用是复制一个map json, 解决在收集筛选结果的时候复制了map本身，底层数据仍然公用
//当然，数据的复制是有性能代价的
func CopyMap(m map[string]interface{}) map[string]interface{} {
	tmp, err := NewEjson(m)

	if err == nil {
		result, err := tmp.GetDataInNewMap()
		if err == nil {
			return result
		}
	}

	return make(map[string]interface{})
}
