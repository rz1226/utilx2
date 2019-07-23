package logx

import (
	"bytes"
	"fmt"
	"github.com/xiaobai22/utilx2/circleq"
	"strconv"
	"time"
)

type Logs struct {
	logs map[string]*circleq.CQ
}

func NewLogs(lognames ...string) *Logs {
	l := &Logs{}
	l.logs = make(map[string]*circleq.CQ)
	for _, v := range lognames {
		l.logs[v] = circleq.NewCQ(500)
	}
	l.logs["defaultInfo"] = circleq.NewCQ(500)
	l.logs["defaultError"] = circleq.NewCQ(500)
	return l
}
func (l *Logs) GetContents(logname string, count int) string {
	q, ok := l.logs[logname]
	if !ok {
		return logname + " 没有这个日志队列\n"
	}
	values, seq := q.GetSeveral(count)
	return formatLog(values, seq)
}

//返回完整的字符串
func (l *Logs) PutContentsAndFormat(logname string, a ...interface{}) string {
	if len(logname) == 0 {
		return ""
	}
	q, ok := l.logs[logname]
	if !ok {
		return ""
	}
	b := bytes.Buffer{}
	b.WriteString(logname)
	b.WriteString(" ")
	b.WriteString(time.Now().Format("2006-01-02 15:04:05"))
	b.WriteString(" ")
	b.WriteString(fmt.Sprintln(a...))
	logStr := b.String()
	q.Put(logStr)
	return logStr
}

func formatLog(values []interface{}, seq uint64) string {
	b := bytes.Buffer{}
	b.WriteString("日志序号: " + strconv.FormatUint(seq, 10) + "\n")
	for _, v := range values {
		str, ok := v.(string)
		if ok {
			b.WriteString(str)
		}
	}
	return b.String()
}
