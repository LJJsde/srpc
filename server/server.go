package server

import (
	"fmt"
	"log"
	"net"
	"reflect"
	"roc/codec"
	"roc/trans"
	"sync"
)

//声明服务器
type Serve struct {
	// 地址
	Addr string
	// 通过反射绑定函数
	Funcs map[string]reflect.Value
}

// NewServer 初始化服务对象
func NewServer(addr string) *Serve {
	return &Serve{
		Addr:  addr,
		Funcs: make(map[string]reflect.Value),
	}
}

var wg sync.WaitGroup

// Register 绑定注册方法
// 服务名，传入的函数
func (s *Serve) Register(name string, f interface{}) {
	if _, ok := s.Funcs[name]; ok {
		return
	}
	Tstruct := reflect.ValueOf(f)
	s.Funcs[name] = Tstruct
}

// Start 服务开始
func (s *Serve) Start() {
	wg.Add(1)
	lis, err := net.Listen("tcp", s.Addr)
	if err != nil {
		fmt.Println("listen err:", err)
		return
	}
	for{
		//开启监听
		conn, err := lis.Accept()
		if err != nil {
			fmt.Println("accept err", err)
			return
		}
		trs := trans.NewTransport(conn)
		go func() {
			for {
				info, err := trs.Receive()
				if err != nil {
					fmt.Println("recieve info err!:", err)
				}
				f, ok := s.Funcs[info.Name]
				if !ok {
					fmt.Println("func no life err")
				}
				// 将客户端返回的参数，转入一个数组中
				inArgs := make([]reflect.Value, 0, len(info.Args))
				for i := 0; i < len(info.Args); i++ {
					inArgs = append(inArgs, reflect.ValueOf(info.Args[i]))
				}
				if info.Method != "" {
					method := f.MethodByName(info.Method)
					log.Println("user the method ", info.Method, "from ", info.Name)
					// co in info
					if info.Args[0] == "" && len(info.Args) == 1 {
						res := method.Call(make([]reflect.Value, 0))
						//将得到的结果提取放入数组
						resArgs := make([]interface{}, 0, len(res))
						for i := 0; i < len(res); i++ {
							resArgs = append(resArgs, res[i].Interface())
						}
						err = trs.Send(codec.Info{
							Name:   info.Name,
							Method: info.Method,
							Args:   resArgs,
						})
						if err != nil {
							fmt.Println("server send err", err)
							return
						}
					} else {
						res := method.Call(inArgs)
						//将得到的结果提取放入数组
						resArgs := make([]interface{}, 0, len(res))
						for i := 0; i < len(res); i++ {
							resArgs = append(resArgs, res[i].Interface())
						}
						err = trs.Send(codec.Info{
							Name:   info.Name,
							Method: info.Method,
							Args:   resArgs,
						})
						if err != nil {
							fmt.Println("server send err", err)
							return
						}
					}
				} else {
					conn.Close()
					break
				}
				//method := f.MethodByName(info.Method)
			}
		}()
	}
	wg.Wait()
}
