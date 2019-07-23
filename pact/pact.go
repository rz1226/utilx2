package pact

import (
	"bytes"
	"encoding/binary"
	"errors"

	"io"
	"log"
)

//简单协议
//读取
/*
分解
x := bytes.NewBuffer([]byte(str))
somestr, err := pact.Read(x)
somestr2, err := pact.Read(x)
somestr3, err := pact.Read(x)


组装
str := new(bytes.Buffer)
pact.Write(str, somestr1)
pact.Write(str, somestr2)
pact.Write(str, somestr3)
all := str.String()
*/

//读出一个完整的数据包
func Read(r io.Reader) ([]byte, error) {
	header := make([]byte, 4)
	lengthReadHeader, err := io.ReadFull(r, header)

	if err != nil {
		return nil, err
	}
	if lengthReadHeader != 4 {
		return nil, errors.New("can not read header ")
	}

	headerReader := bytes.NewReader(header)
	var lengthBody int32
	binary.Read(headerReader, binary.BigEndian, &lengthBody)
	body := make([]byte, int(lengthBody))

	lengthReadBody, err := io.ReadFull(r, body)
	if err != nil {
		return nil, err
	}
	if lengthReadBody == 0 {
		//长度是0的数据也是合法的
		return body, nil
	}
	return body, nil
}

func Write(w io.Writer, data []byte) error {
	lengthData := len(data)
	headerBuffer := new(bytes.Buffer)
	err := binary.Write(headerBuffer, binary.BigEndian, int32(lengthData))
	if err != nil {
		return err
	}
	lengthHeader, err := w.Write(headerBuffer.Bytes())
	if err != nil {
		return err
	}
	if lengthHeader != len(headerBuffer.Bytes()) {
		log.Println("this should not happen")
	}


	body := []byte{}
	body = append(body, data...)
	lengthBody, err := w.Write(body)
	if err != nil {
		return err
	}
	if lengthBody != len(body) {
		log.Println("this should not happening 2")
		return errors.New("write part of the body ")
	}
	return nil

}
