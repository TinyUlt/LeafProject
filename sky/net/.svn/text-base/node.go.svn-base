package net

import (
	"sky/jsonhelper"
	//"sky/log"
	"time"
)

type Node struct {
	//conn variable
	minPacketLen       int
	maxPacketLen       int
	recvBuffLen        int
	sendBuffLen        int
	readBuffLen        int
	writeBuffLen       int
	lingerTime         int
	maxConnection      uint32
	optimalConnection  uint32
	preallocConnection uint32

	//computing variable
	recvAvailableLen int
	sendAvailableLen int
	sendFixedLen     int
	timeout          time.Duration
	flowLimit        int
	//string variable
	secure string
	codec  string
}

func newNode(conf jsonhelper.JSONObject) *Node {
	node := new(Node)
	node.minPacketLen = conf.GetAsInt("min_packet_size")
	if node.minPacketLen == 0 {
		node.minPacketLen = PacketHeadSize
	}

	node.maxPacketLen = conf.GetAsInt("max_packet_size")
	if node.maxPacketLen == 0 {
		node.maxPacketLen = PacketMaxSize
	}

	node.readBuffLen = conf.GetAsInt("read_buff_size")
	if node.readBuffLen == 0 {
		node.readBuffLen = node.maxPacketLen * 4
	}

	node.writeBuffLen = conf.GetAsInt("write_buff_size")
	if node.writeBuffLen == 0 {
		node.writeBuffLen = node.maxPacketLen * 4
	}

	node.recvBuffLen = conf.GetAsInt("recv_buff_size")
	if node.recvBuffLen == 0 {
		node.recvBuffLen = node.maxPacketLen * 4
	}

	node.sendBuffLen = conf.GetAsInt("send_buff_size")
	if node.sendBuffLen == 0 {
		node.sendBuffLen = node.maxPacketLen * 4
	}
	node.lingerTime = conf.GetAsInt("linger_time")
	if node.lingerTime == 0 {
		node.lingerTime = -1
	}

	node.maxConnection = uint32(conf.GetAsInt("max_connection"))

	node.optimalConnection = uint32(conf.GetAsInt("optimal_connection"))

	node.preallocConnection = uint32(conf.GetAsInt("prealloc_connection"))

	if t := conf.GetAsInt("timeout"); t > 0 {
		node.timeout = time.Millisecond * time.Duration(t)
	}

	node.flowLimit = conf.GetAsInt("flow_limit")
	if node.flowLimit == 0 {
		node.flowLimit = 100
	}

	node.codec = conf.GetAsString("codec")
	node.secure = conf.GetAsString("secure")

	if node.optimalConnection == 0 {
		if node.maxConnection == 0 {
			node.maxConnection = 2
			node.optimalConnection = 1
		} else {
			node.optimalConnection = node.maxConnection / 2
		}
	}
	if node.maxConnection == 0 {
		node.maxConnection = node.optimalConnection * 2
	}

	if node.preallocConnection == 0 || node.preallocConnection > node.optimalConnection {
		node.preallocConnection = node.optimalConnection
	}

	node.recvAvailableLen = node.recvBuffLen - node.maxPacketLen
	node.sendAvailableLen = node.sendBuffLen - node.maxPacketLen
	node.sendFixedLen = node.sendAvailableLen - node.maxPacketLen
	return node
}
