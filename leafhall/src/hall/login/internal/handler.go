package internal

import (
	"LeafProject/leaf/log"
	"reflect"
	"hall/msg"
	"LeafProject/leaf/gate"
)

func init() {
	handler(&msg.UserLoginRequest{}, handleLoginRequest)
}
func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)

}


func handleLoginRequest(args []interface{}) {
	// 收到的 Hello 消息
	m := args[0].(*msg.UserLoginRequest)
	// 消息的发送者
	a := args[1].(gate.Agent)

	// 输出收到的消息的内容
	log.Debug("username: %v", m.UserName)

	// 给发送者回应一个 Hello 消息
	a.WriteMsg(&msg.Hello{
		Name: "client",
	})
}