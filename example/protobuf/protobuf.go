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