package main

import (
	"fmt"
	client2 "roc/client"
)

func main() {
	//for {
	//	client:=client2.ClientRun(":8080","Helloa")
	//	// 多参数，不同类型的参数
	//	out1 := client.NewCall("DouArgs", "lijia", 11)
	//	for i, _ := range out1 {
	//		fmt.Println(out1[i])
	//	}
	//}
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
