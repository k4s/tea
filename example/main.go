package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:4444")
	if err != nil {
		panic(err)
	}

	// Hello 消息（JSON 格式）
	data := []byte(`{"Hello": {"Name": "tea"}}`)

	// len + data
	m := make([]byte, 4+len(data))

	// 默认使用大端序
	binary.LittleEndian.PutUint32(m, uint32(len(data)))
	// binary.BigEndian.PutUint32(m, uint32(len(data)))

	copy(m[4:], data)

	// 发送消息
	conn.Write(m)

	var buf = make([]byte, 32)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("read error:", err)
	} else {
		fmt.Printf("read % bytes, content is %s\n", n, string(buf[:n]))
	}
	fmt.Println(string(buf))

	fmt.Println(binary.LittleEndian.Uint32(buf))
	fmt.Println(len(buf))
	time.Sleep(time.Second * 3)
}
