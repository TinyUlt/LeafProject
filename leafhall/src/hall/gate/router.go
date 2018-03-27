package gate
import (
	"hall/game"
	"hall/msg"
	"hall/login"
	//"reflect"
	"reflect"
)
func init() {
	// 这里指定消息 Hello 路由到 game 模块
	// 模块间使用 ChanRPC 通讯，消息路由也不例外
	msg.Processor.SetRouter(&msg.UserLoginRequest{}, login.ChanRPC)
	msg.Processor.SetRouter(&msg.Hello{}, game.ChanRPC)

	game.ChanRPC.Go(reflect.TypeOf(&msg.UserEnterRoomRequest{}),nil)
}
