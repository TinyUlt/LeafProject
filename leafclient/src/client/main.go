package main

import (
	"encoding/binary"
	"net"
	"client/msg"
	"github.com/golang/protobuf/proto"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:3563")
	if err != nil {
		panic(err)
	}

	// Hello 消息（JSON 格式）
	// 对应游戏服务器 Hello 消息结构体

	d :=&msg.Hello{}
	d.Name = "hello"
	data, err := proto.Marshal(d)

	//id + data
	i := make([]byte, 2+len(data))
	binary.BigEndian.PutUint16(i, uint16(0))
	copy(i[2:], data)

	// len + data
	m := make([]byte, 2+len(i))

	// 默认使用大端序
	binary.BigEndian.PutUint16(m, uint16(len(i)))

	copy(m[2:], i)


	// 发送消息
	conn.Write(m)
}