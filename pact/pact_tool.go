package pact

//生成包内容数据
func Make(no byte, data []byte) []byte {
	return append([]byte{no}, data...)
}

//解出
func Get(data []byte) (byte, []byte) {
	return data[0], data[1:]
}
