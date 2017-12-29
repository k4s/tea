## Tea 游戏服务器框架
---

交流QQ群：376389675

Tea [V1.0](https://github.com/k4s/tea/tree/v1.0)：
* 像开发Web一样简单去开发Game.
* 每个用户走在单个goroutine上，更加适合多核支持并发处理.

Tea [masterr](https://github.com/k4s/tea/tree/master)：
* 支持多网关，多游戏服 分布式处理.
* 单路复用.
* 支持v1.0，leaf 游戏逻辑对接.



#### 基于Tea的开发：

1.新建生成项目工具newTea：
```
cd github.com/k4s/tea/newTea
go install
```

2.生成网关：
```
cd $GOPATH
newTea new gate appname
cd appname
```

3.生成游戏服：
```
cd $GOPATH
newTea new game appname
cd appname
```

4.配置[config]目录，选择一种msg协议作为通讯协议，对应的[protocol/process.go]：
```
Protocol  = "json"
```

5.在[hamdle]编写对应msg的处理函数.
```go 
func InfoHandle(msg *message.Message, agent network.Agent) {
	jsonMsg, err := protocol.Processor.Unmarshal(msg.Body)
	if err != nil {
		fmt.Println(err)
	}
	m := jsonMsg.(*ms.Hello)
	fmt.Println("game:", m)
	reMsg := ms.Hello{
		Name: "kkk",
	}
	agent.EndHandle(msg, reMsg)
}
```

6.在[register]做通讯消息注册
```
protocol.Processor.Register(&msg.Hello{})
```
7.在[router]做路由映射.
```
protocol.Processor.SetHandler(&msg.Hello{}, handle.InfoHandle)
```


8.分别执行网关和游戏服：
```
cd appname
go run main.go
```


#### unity 引擎游戏客户端测试

socket and websocket with json demo:

[github.com/k4s/teaUnity](http://github.com/k4s/teaUnity)




#### go 客户端测试

[https://github.com/k4s/tea/example](https://github.com/k4s/tea/tree/master/example)