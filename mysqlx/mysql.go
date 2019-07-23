package mysqlx

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"
)

/*

var conf = "gechengzhen:123456@tcp(172.16.1.61:3306)/userdata?charset=utf8"
func main(){

	p := mysqlx.NewDbPool( conf,3 )

	for i := 0; i < 1000000 ; i++ {
		go test( p )
		time.Sleep( time.Millisecond * 1 )
	}


	time.Sleep( time.Second * 10000 )
}

func test( p *dbtool.DbPool ){

	 fmt.Println( p.Query("select * from ytk_car_test limit 1 ", nil ))

	//fmt.Println( p.Exec("insert into ytk_car_test set license = '赣B'", nil))
}
*/

type DbPool struct {
	realPool     *sql.DB
	conStr       string
	maxOpenConns int
	maxIdleConns int
	lastErr      error
}

func NewDbPool(conStr string, maxOpenConns int) *DbPool {
	p := &DbPool{}
	p.conStr = conStr
	p.realPool, p.lastErr = sql.Open("mysql", conStr)
	if p.lastErr == nil {
		p.realPool.SetMaxOpenConns(maxOpenConns)
		p.realPool.SetMaxIdleConns(maxOpenConns)
		p.realPool.SetConnMaxLifetime(time.Second * 10000)
	} else {
		p.realPool = nil
	}
	return p
}

//获取*sql.DB
func (p *DbPool) DB() *sql.DB {
	return p.realPool
}

func (p *DbPool) Query(sqlStr string, data []interface{}) (SelectRes, error) {
	rows, err := p.query(sqlStr, data)
	if err == nil {
		return rowsToMap(rows)
	} else {
		return SelectRes{nil}, err
	}

}

func (p *DbPool) query(sqlStr string, data []interface{}) (*sql.Rows, error) {
	if p.realPool == nil {
		return nil, p.lastErr
	}
	db := p.realPool
	length := len(data)
	fn := reflect.ValueOf(db.Query)
	fnParams := make([]reflect.Value, length+1)
	fnParams[0] = reflect.ValueOf(sqlStr)
	for i := 1; i <= length; i++ {
		fnParams[i] = reflect.ValueOf(data[i-1])
	}
	//fmt.Println( params )
	callResult := fn.Call(fnParams)
	if callResult[1].Interface() != nil {
		return nil, callResult[1].Interface().(error)
	}
	return callResult[0].Interface().(*sql.Rows), nil
}

// select 操作返回的结果
type SelectRes struct {
	data []map[string]interface{}
}

func (q SelectRes) Len() int {
	return len(q.data)
}

func (q SelectRes) Data() []map[string]interface{} {
	return q.data
}

//把数据库取出来d数据的rows的数据放在一个QueryResult上
// null 对应nil  数字对数字  其他对字符串

func rowsToMap(rows *sql.Rows) (SelectRes, error) {
	defer rows.Close()
	res := make([]map[string]interface{}, 0, 100)
	fields, err := rows.Columns()
	lengthRow := len(fields)
	if err != nil {
		return SelectRes{nil}, err
	}

	for {
		if result := rows.Next(); result {
			scanRes := make([]interface{}, lengthRow)
			resultData := make(map[string]interface{}, lengthRow)
			vScanRes := reflect.ValueOf(&scanRes)
			fn := reflect.ValueOf(rows.Scan)
			fnParams := make([]reflect.Value, lengthRow)
			for i := 0; i < lengthRow; i++ {
				fnParams[i] = vScanRes.Elem().Index(i).Addr()
			}
			callResult := fn.Call(fnParams)
			if callResult[0].Interface() != nil {
				return SelectRes{nil}, callResult[0].Interface().(error)
			}
			for i := 0; i < lengthRow; i++ {
				resultData[fields[i]] = scanRes[i]
			}
			res = append(res, resultData)

		} else {
			break
		}
	}
	sr := SelectRes{}
	sr.data = res
	return sr, nil
}

//datas是个二维slice,
//注意，这仍然是每个条数据遍历执行，插入的时候并不是自动生成批量sql插入的。
func (p *DbPool) ExecMany(sqlStr string, datas [][]interface{}) (int64, error) {
	if p.realPool == nil {
		return 0, p.lastErr
	}
	db := p.realPool
	var rowsAffected int64 = 0
	var lastInsertId int64 = 0

	//插入数据
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	if err != nil {
		return 0, err
	}
	fn := reflect.ValueOf(stmt.Exec)
	for _, data := range datas {
		length := len(data)
		fnParams := make([]reflect.Value, length)
		for i := 0; i < length; i++ {
			fnParams[i] = reflect.ValueOf(data[i])
		}
		callResult := fn.Call(fnParams)
		if callResult[1].Interface() != nil {
			fmt.Println("ExecMany error : ", callResult[1].Interface())
			fmt.Println("continue")
			continue
			//return 0, fv[1].Interface().(error)
		}
		result := callResult[0].Interface().(sql.Result)
		if isUpdate(sqlStr) || isDelete(sqlStr) {
			singleAffected, err := result.RowsAffected()
			if err != nil {
				fmt.Println("ExecMany error : ", err)
				fmt.Println("continue")
				continue
			}
			rowsAffected += singleAffected
		}
		if isInsert(sqlStr) {
			thisTimeLastId, err := result.LastInsertId()
			if err != nil {
				fmt.Println("ExecMany error : ", err)
				fmt.Println("continue")
				continue
			}
			lastInsertId = thisTimeLastId
		}
	}
	if isUpdate(sqlStr) || isDelete(sqlStr) {
		return rowsAffected, nil
	}
	if isInsert(sqlStr) {
		return lastInsertId, nil
	}
	return 0, errors.New("only support update insert delete 4")
}

// data 是一个slice, 里面的个数对应sqlStr里面？的数量
func (p *DbPool) Exec(sqlStr string, data []interface{}) (int64, error) {
	if p.realPool == nil {
		return 0, p.lastErr
	}
	db := p.realPool
	//插入数据
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	length := len(data)
	if err != nil {
		return 0, err
	}

	fn := reflect.ValueOf(stmt.Exec)
	fnParams := make([]reflect.Value, length)
	for i := 0; i < length; i++ {
		fnParams[i] = reflect.ValueOf(data[i])
	}
	callResult := fn.Call(fnParams)

	if callResult[1].Interface() != nil {
		return 0, callResult[1].Interface().(error)
	}
	result := callResult[0].Interface().(sql.Result)
	if isUpdate(sqlStr) || isDelete(sqlStr) {
		return result.RowsAffected() //本身就是多个返回值
	}
	if isInsert(sqlStr) {
		return result.LastInsertId() //本身就是多个返回值
	}
	return 0, errors.New("only support update insert delete 5")
}

func isInsert(sqlStr string) bool {
	str := strings.TrimSpace(strings.ToLower(sqlStr))
	if strings.HasPrefix(str, "insert") {
		return true
	}
	return false
}

func isUpdate(sqlStr string) bool {
	str := strings.TrimSpace(strings.ToLower(sqlStr))
	if strings.HasPrefix(str, "update") {
		return true
	}
	return false
}

func isDelete(sqlStr string) bool {
	str := strings.TrimSpace(strings.ToLower(sqlStr))
	if strings.HasPrefix(str, "delete") {
		return true
	}
	return false
}

//
func Int64(data interface{}) int64 {
	if data == nil {
		return 0
	}
	str := string(data.([]uint8))
	num, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		fmt.Println("int64 conv error")
		log.Fatal(err)
	}
	return num
}

func Int(data interface{}) int64 {
	if data == nil {
		return 0
	}
	str := string(data.([]uint8))
	num, err := strconv.ParseInt(str, 10, 0)
	if err != nil {
		fmt.Println("int conv error")
		log.Fatal(err)
	}
	return num
}

func String(data interface{}) string {
	if data == nil {
		return ""
	}
	str := string(data.([]uint8))
	return str
}

func Float64(data interface{}) float64 {
	if data == nil {
		return 0
	}
	str := string(data.([]uint8))
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		fmt.Println(err)
	}
	return f

}

/*************************组成批量sql插入语句*********************************/

//生成用来批量插入的参数
func BatchInsertParams(datas [][]interface{}) (string, []interface{}) {
	sqlParams := make([]interface{}, 0, 200)
	sqlStringBuffer := bytes.Buffer{}
	for _, oneData := range datas {
		length := len(oneData)
		if length == 0 {
			continue
		}
		sqlStringBuffer.WriteString("(")
		for idx, val := range oneData {
			sqlParams = append(sqlParams, val)
			if length == idx+1 {
				sqlStringBuffer.WriteString("?")
			} else {
				sqlStringBuffer.WriteString("?,")
			}
		}
		sqlStringBuffer.WriteString("),")
	}
	sql := strings.TrimRight(sqlStringBuffer.String(), ",")
	return sql, sqlParams
}
