package msg

import (

	"LeafProject/leaf/network/protobuf"
)

var Processor = protobuf.NewProcessor()

func init() {
	Processor.Register(&Hello{},uint16(MessageId__Hello))
	Processor.Register(&UserLoginRequest{},uint16(MessageId__UserLoginRequest))
}
