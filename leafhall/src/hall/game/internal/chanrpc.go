package internal

import (
	"LeafProject/leaf/gate"
	"fmt"
	"reflect"
	"hall/msg"
)

func init() {
	skeleton.RegisterChanRPC("NewAgent", rpcNewAgent)
	skeleton.RegisterChanRPC("CloseAgent", rpcCloseAgent)
	skeleton.RegisterChanRPC(reflect.TypeOf(&msg.UserEnterRoomRequest{}), test)
}
func test(args []interface{}) {
	fmt.Print("test")
}
func rpcNewAgent(args []interface{}) {
	a := args[0].(gate.Agent)
	_ = a
}

func rpcCloseAgent(args []interface{}) {
	a := args[0].(gate.Agent)
	_ = a
}
