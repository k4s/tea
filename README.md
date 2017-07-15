## Tea 游戏服务器框架
---
Tea 的基础模块来自于 [Leaf](https://github.com/name5566/leaf)，由于想做成更加简单易用的GameFrame，所以就有了[Tea](https://github.com/k4s/tea).

交流QQ群：376389675

Tea ：
* 像开发Web一样简单去开发Game.
* 每个用户走在单个goroutine上，更加适合多核支持并发处理.

#### Tea 的框架：

```
┌-----| |---| |---| |---------------| |---| |---| |------┐
|     | |   | |   | |     Gate      | |   | |   | |      |
|  ┌--| |---| |---| |---------------| |---| |---| |--┐   |
|  |  |a|   |a|   |a|               |a|   |a|   |a|  |   |
|  |  |g|   |g|   |g|   TCP Server  |g|   |g|   |g|  |   |
|  |  |e|   |e|   |e|      or       |e|   |e|   |e|  |   |
|  |  |n|   |n|   |n|   WebSocket   |n|   |n|   |n|  |   |
|  |  |t|   |t|   |t|               |t|   |t|   |t|  |   |
|  └--| |---| |---| |---------------| |---| |---| |--┘   |
|  ┌--| |---| |---| |---------------| |---| |---| |--┐   |
|  |  | |   | |   | |  Msg Process  | |   | |   | |  |   |
|  └--| |---| |---| |---------------| |---| |---| |--┘   |
└-----| |---| |---| |---------------| |---| |---| |------┘
┌-----| |---| |---| |---------------| |---| |---| |------┐
|                 Goroutine Deal with Msg                |
└--------------------------------------------------------┘
```

#### 基于Tea的开发：

1.生成项目：
```
cd github.com/k4s/tea/newTea
go install
cd $GOPATH
newTea new appname
cd appname
```
2.配置[conf]目录，选择一种msg协议作为通讯协议，编写对应的[msg/process/process.go]：
```
Protocol  = "json"
```

3.在[game]编写对应msg的处理函数.
```go 
func MsgHello(msg interface{}, agent network.InAgent) {
	m := msg.(*ms.Hello)
	fmt.Println("Hello,", m.Name)
	hi := &ms.Hello{Name: "kas"}
	agent.WriteMsg(hi)
}
```

4.在[msg/register]做通讯消息注册
```
process.Processor.Register(&Hello{})
```

5.在[router]做路由映射.
```
process.Processor.SetHandler(&msg.Hello{}, game.MsgHello)
```
6.如果需要在agent新建或者关闭执行函数，在[game/agent]里面编写对应的函数.

1) 新建agent：
```go
func initAgent(agent network.ExAgent) {
	fmt.Println("InitFunc")
}
```
2) 关闭agent：

```go
func closeAgent(agent network.ExAgent) {
	fmt.Println("CloseFunc")
}

```

7.执行你的程序：
```
cd teaserver
go run main.go
```
#### unity 引擎游戏客户端测试
websocket and json

[github.com/k4s/teaUnityWebsocket](http://github.com/k4s/teaUnityWebsocket)

#### go 客户端测试
```go
package main

import (
    "encoding/binary"
    "net"
)

func main() {
    conn, err := net.Dial("tcp", "127.0.0.1:8080")
    if err != nil {
        panic(err)
    }

    // Hello 消息（JSON 格式）
    data := []byte(`{
        "Hello": {
            "Name": "tea"
        }
    }`)

    // len + data
    m := make([]byte, 2+len(data))

    // 默认使用大端序
    binary.BigEndian.PutUint16(m, uint16(len(data)))

    copy(m[2:], data)

    // 发送消息
    conn.Write(m)
}
```
protobuf:
```go
package main

import (
    "encoding/binary"
    "net"
    "fmt"

    "appname/msg"
    "github.com/golang/protobuf/proto"
)

func main() {
    conn, err := net.Dial("tcp", "127.0.0.1:8080")
    if err != nil {
        panic(err)
    }

    login := &msg.Login{
       Name:"kas",
       Password:"123456",
   }
   // 进行编码
    data, err := proto.Marshal(login)
    if err != nil {
        panic( err)
    }
    m := make([]byte, 4+len(data))
    // // 默认使用大端序
    binary.BigEndian.PutUint16(m, uint16(2+len(data)))
    binary.BigEndian.PutUint16(m[2:], uint16(0))
    copy(m[4:], data)

    // 发送消息
    conn.Write(m)
    // 接受消息
    var buf = make([]byte, 32)
    n, err := conn.Read(buf)
    if err != nil {
        fmt.Println("read error:", err)
    } else {
        fmt.Printf("read % bytes, content is %s\n", n, string(buf[:n]))
    }
    readlogin := &msg.Login{}
    err = proto.Unmarshal(buf[4:n],readlogin)
    if err != nil{
        fmt.Println(err)
    }
    fmt.Println("readlogin:",readlogin)

}
```