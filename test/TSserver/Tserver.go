package main

import (
	"roc/server"
	"strconv"
)

func main() {
	se:=server.NewServer(":8080")
	se.Register("Helloa",&EchoService{})
	se.Start()
}

type EchoService struct{}

func (*EchoService)DouArgs(name string,num int)string{
	return "hello "+name+"  "+strconv.Itoa(num)
}

func (*EchoService)Hello(name string)string{
	return "hello"+name
}

func (*EchoService)Do()string{
	return "猜猜没有参数能调用吗？芜湖"
}