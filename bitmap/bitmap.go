package bitmap

import (

	"math/rand"
	"sync"
	"time"
)


type BitMap struct {
	m            sync.RWMutex
	byteLine         []byte
	lastTimeSyncToStorage int64 //最后同步时间
	dirtyCount  int   //最后一次同步之后又经过了多少次修改
}

func NewBitMap(dataInit []byte) *BitMap {
	bm := &BitMap{}
	bm.byteLine = dataInit
	bm.dirtyCount = 0
	bm.lastTimeSyncToStorage = time.Now().Unix()
	return bm
}

func (this *BitMap) GetBytes() []byte {
	return this.byteLine
}

//加长，如果已经够长了，什么都不操作,否则全部补充零
//注意该长度是指byte的数量长度,不是位的长度
func (this *BitMap) padWithZero(lenth int) {
	clen := len(this.byteLine)
	if clen < lenth {
		NewBiggerByteLine := make([]byte, lenth+lenth/5)
		copy(NewBiggerByteLine, this.byteLine)
		this.byteLine = NewBiggerByteLine
	}
}

//对外接口
func (this *BitMap) GetPositionValue(position int) bool {
	this.m.RLock()
	defer this.m.RUnlock()
	pos := position / 8
	if pos > len(this.byteLine)-1 {
		return false
	}
	value := this.byteLine[pos]
	mod := position % 8
	newValue := GetSingleBytePositionValue(value, mod)
	return newValue

}

//对外接口
func (this *BitMap) SetPositionTrue(position int) {
	this.SetPostion(position, true)
}
func (this *BitMap) SetPositionFalse(position int) {
	this.SetPostion(position, false)
}

func (this *BitMap) SetPostion(position int, val bool) {
	this.m.Lock()
	defer this.m.Unlock()
	whichByte := position / 8
	if whichByte > len(this.byteLine)-1 {
		this.padWithZero(whichByte + 10)
	}
	value := this.byteLine[whichByte]
	mod := position % 8
	newValue := SetSingleBytePositionValue(value, mod, val)
	this.update(whichByte, newValue)
}
func (this *BitMap) Len() int {
	return len(this.byteLine) * 8
}

//update
func (this *BitMap) update(pos int, val byte) bool {
	if len(this.byteLine)-1 < pos {
		this.padWithZero(pos + 10)
	}
	this.byteLine[pos] = val
	this.dirtyCount++
	return true

}
func (this *BitMap) DirtyCount()int{
	return this.dirtyCount
}

func minLen(t, t2 *BitMap) int {
	if len(t.byteLine) <= len(t2.byteLine) {
		return len(t.byteLine)
	}
	return len(t2.byteLine)
}

func maxLen(t, t2 *BitMap) int {
	if len(t.byteLine) >= len(t2.byteLine) {
		return len(t.byteLine)
	}
	return len(t2.byteLine)
}

func Or(b, b2 *BitMap) *BitMap {
	len := maxLen(b, b2)
	b.padWithZero(len)
	b2.padWithZero(len)
	bm := &BitMap{}
	bm.byteLine = make([]byte, len)
	for i := 0; i < len; i++ {
		bm.byteLine[i] = b.byteLine[i] | b2.byteLine[i]
	}
	return bm
}

func And(b, b2 *BitMap) *BitMap {
	len := minLen(b, b2)
	bm := &BitMap{}
	bm.byteLine = make([]byte, len)
	for i := 0; i < len; i++ {
		bm.byteLine[i] = b.byteLine[i] & b2.byteLine[i]
	}
	return bm
}

func GetSingleBytePositionValue(value byte, bitpos int) bool {

	reversedPosition := byte(7 - bitpos)
	factor := byte(1) << reversedPosition
	if factor == factor&value {
		return true
	}

	return false

}

func SetSingleBytePositionValue(haystack byte, position int, val bool) byte {

	reversedPosition := byte(7 - position)
	if val == true {
		factor := byte(1) << reversedPosition
		if factor == factor&haystack {
			return haystack
		} else {
			return haystack | factor
		}
	} else {
		factor := byte(1) << reversedPosition

		if byte(0) == factor&haystack {
			return haystack
		} else {
			return haystack & ^factor
		}
	}
}

//用于测试
func ByteToBinaryString(data byte) (str string) {
	var a byte
	for i := 0; i < 8; i++ {
		a = data
		data <<= 1
		data >>= 1

		switch a {
		case data:
			str += "0"
		default:
			str += "1"
		}
		data <<= 1
	}
	return str
}

func getRand() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(255)
}
