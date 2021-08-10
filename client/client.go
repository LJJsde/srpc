package client

import (
	"fmt"
	"log"
	"net"
	"reflect"
	"roc/codec"
	"roc/trans"
)



type Client struct {
	ServerName string
	Conn net.Conn
	trs  *trans.Transport

}
//type ReChan struct {
//	sync.Mutex
//	reGo chan bool
//}

func ClientRun(router string,serverName string)*Client{
	conn, err := net.Dial("tcp", router)
	if err != nil {
		fmt.Println("dial err:", err)
		return nil
	}
	return NewClient(conn,serverName)
}

func NewClient(conn net.Conn,serverName string) *Client {
	trs:=trans.NewTransport(conn)
	return &Client{serverName,conn,trs}
}
// 在客户端重构函数式的调用
func (c *Client) Call(rpcN string, in interface{}) {
	fn := reflect.ValueOf(in).Elem()
	f := func(args []reflect.Value) []reflect.Value {
		inArgs := make([]interface{}, 0, len(args))
		for _, arg := range args {
			inArgs = append(inArgs, arg.Interface())
		}
		inInfo := codec.Info{
			Name: rpcN,
			Args: inArgs,
		}
		err := c.trs.Send(inInfo)
		if err != nil {
			log.Println("client send err:", err)
		}
		getInfo, err := c.trs.Receive()
		if err != nil {
			log.Println("client get err:", err)
		}
		getArgs := make([]reflect.Value, len(getInfo.Args))
		for i := 0; i < len(getInfo.Args); i++ {
			getArgs[i] = reflect.ValueOf(getInfo.Args[i])
		}
		return getArgs

	}
	funcs := reflect.MakeFunc(fn.Type(), f)
	fn.Set(funcs)
}

// 远程调用函数得到他的返回值
func (c *Client) NewCall(methodName string,in ...interface{})[]reflect.Value{
	//inArg:=make([]reflect.Value,len(in))
	//for i:=0;i<len(in);i++{
	//	inArg[i]=reflect.ValueOf(in[i])
	//}
	err:=c.trs.Send(codec.Info{
		Name:   c.ServerName,
		Method: methodName,
		Args:   in,
	})
	if err!=nil{
		fmt.Println("client send err",err)
		return nil
	}
	var outArgs codec.Info
	outArgs, err = c.trs.Receive()
	if err!=nil{
		fmt.Println("client receive err:",err)
		return nil
	}
	getInfo:=make([]reflect.Value,len(outArgs.Args))
	for i := 0; i < len(outArgs.Args); i++ {
		getInfo[i] = reflect.ValueOf(outArgs.Args[i])
	}

	return getInfo
}