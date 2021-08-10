# README

一款基本的rpc框架，利用反射，编码等完成而来远程的函数调用。

## 服务端

```go
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
```

我们通过绑定结构体来得到他所对应的方法，注册服务名。

## 客户端

```go
func main() {
	client:=client2.ClientRun(":8080","Helloa")
	// 多参数，不同类型的参数
	out1 := client.NewCall("DouArgs", "lijia", 11)
	for i, _ := range out1 {
		fmt.Println(out1[i])
	}

	// 单参数
	out := client.NewCall("Hello", "lijia")
	for j, _ := range out {
		fmt.Println(out[j])
	}
	//// 无参数类型
	out2:=client.NewCall("Do","")
	for k,_:=range out2{
		fmt.Println(out2[k])
	}
}
```

在得到客户端的时候输入的是端口号以及服务端所注册的服务名，必须与服务端一致。

调用函数的时候使用NewCall(),第一个参数是该服务有的方法名，后续参数则是传入的参数。如果你的远程注册方法不需要传入参数则你只需要传入一个“ ”就能达到效果。

![image](https://user-images.githubusercontent.com/73375890/128823344-a4ff780c-5f42-4e33-b5a6-586c97205813.png)
