package trans

import (
	"encoding/binary"
	"io"
	"log"
	"net"
	"roc/codec"
)

// Transport 结构
type Transport struct {
	conn net.Conn
}

// NewTransport 定义一个传输通道
func NewTransport(conn net.Conn) *Transport {
	return &Transport{conn: conn}
}

func (t *Transport) Close() error {
	return t.conn.Close()
}

// Send data
// 1.对info进行编码
// 2.设置请求头(放入body的长度)
// 3.设置body
// 4.发送信息给服务端
func (t *Transport) Send(info codec.Info) error {
	b, err := codec.Encode(info)
	if err != nil {
		return err
	}
	buf := make([]byte, 4+len(b))
	binary.BigEndian.PutUint32(buf[:4], uint32(len(b)))
	copy(buf[4:], b)
	_, err = t.conn.Write(buf)
	if err != nil {
		return err
	}
	return nil
}
// Receive
// 1.从字节序中先拿到前四个，得到传输数据的长度
// 2.将后续长度的字节流中得到数据u
// 3.解码得到info数据
func (t *Transport) Receive() (codec.Info, error) {
	b := make([]byte, 4)
	_, err := io.ReadFull(t.conn, b)
	if err != nil {
		_ = t.Close()
		return codec.Info{}, err
	}
	l:=binary.BigEndian.Uint32(b)
	infoBytes :=make([]byte,l)
	_,err =io.ReadFull(t.conn,infoBytes)
	if err!=nil{
		return codec.Info{}, err
	}
	info,err:=codec.Decode(infoBytes)
	if err!=nil{
		log.Println("decode err",err)
		return codec.Info{},err
	}
	return info,nil
}
