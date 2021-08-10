package codec

import (
	"bytes"
	"encoding/gob"
)

type Info struct {
	Name   string
	Method string
	Args   []interface{}
}

// Encode 对Date进行编码
// 1.先声明一个缓冲器，并且装入god的编码器
// 2.进行编码，返回错误
// 3.最后返回编码后的byte类型
func Encode(info Info) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(info); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Decode 对传来的编码进行解码
// 1.先将b装入缓冲器，并且装入god的解码器
// 2.进行解码，返回错误
// 3.最后返回解码数据
func Decode(b []byte) (Info, error) {
	buf := bytes.NewBuffer(b)
	decoder := gob.NewDecoder(buf)
	var info Info
	if err := decoder.Decode(&info); err != nil {
		return Info{}, err
	}
	return info, nil
}
