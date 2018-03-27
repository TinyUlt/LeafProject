package net

import (
	"reflect"
	"sky/jsonhelper"
	"sky/log"
	"strings"
	"sync"
)

type Packet interface {
	PacketID() uint32
}

type Call struct {
	conn   reflect.Value
	packet reflect.Value
	tls    [2]uint32
	next   *Call
}

type Proc struct {
	sync.Mutex
	recver   reflect.Value
	function reflect.Value
	packetRT reflect.Type //packet type
	free     *Call
	nfree    uint32
}

type Server struct {
	module []*Module
	Input  chan *Call
	conf   jsonhelper.JSONObject
}

func (self *Server) Start() {
	for _, v := range self.module {
		if v != nil {
			v.delegate.Start()
		}
	}
}

func (self *Server) Stop() {
	for _, v := range self.module {
		if v != nil {
			v.delegate.Stop()
		}
	}
}

func (self *Server) Register(id, ident uint32, rcv, fun reflect.Value, rt reflect.Type) {
	self.module[ident].register(id, rcv, fun, rt)
}

func (self *Server) Process(call *Call) {
	module := call.conn.Interface().(*Conn).module
	if module.ep {
		call.conn.Interface().(*Conn).Tls = call.tls
	}
	id := call.packet.Interface().(Packet).PacketID()
	if id > PacketIDReserved {
		id -= module.from
	}
	module.proc[id].function.Call([]reflect.Value{module.proc[id].recver, call.conn, call.packet})
	module.freeCall(id, call)
}

func (self *Server) GetModule(key string) *Module {
	for _, v := range self.module {
		if v != nil && v.name == key {
			return v
		}
	}
	return nil
}

func (self *Server) enqueue(id uint32, conn *Conn, call *Call) {
	if conn.module.sync {
		if id > PacketIDReserved {
			id -= conn.module.from
		}
		if conn.module.ep {
			conn.Tls = call.tls
		}
		conn.module.proc[id].function.Call([]reflect.Value{conn.module.proc[id].recver, call.conn, call.packet})
		conn.module.freeCall(id, call)
	} else {
		self.Input <- call
	}
}

func (self *Server) inflow(idx uint32, b []byte) error {
	//println("route to game", idx)
	//log.Debug(">>>>>>>>> inflow start %v %v", GameIdent, MaxModule)
	if GameIdent < MaxModule && self.module[GameIdent] != nil {
		//log.Debug(">>>>>>>> inflow 111111 %v %v", idx, uint32(len(self.module[GameIdent].conns)))
		if idx < uint32(len(self.module[GameIdent].conns)) {
			//log.Debug("-------------------- inflow %v", idx)
			self.module[GameIdent].conns[idx].forward(b)
			return nil
		}
	}
	//panic(IDRangeError)
	return nil //IDRangeError
}

func (self *Server) hashflow(ident, id uint32, b []byte) error {
	//println("hash route", ident, id)
	if ident < MaxModule && self.module[ident] != nil {
		if conn := self.module[ident].HashConn(id); conn != nil {
			//log.Debug("-------------------- hashflow %v", id)
			conn.forward(b)
			return nil
		}
	}
	//panic(IDRangeError)
	return nil
}

func (self *Server) outflow(user, idx uint32, b []byte) error {
	//println("route to user", user, idx)
	if UserIdent < MaxModule && self.module[UserIdent] != nil {
		if idx < uint32(len(self.module[UserIdent].conns)) {
			if user == self.module[UserIdent].conns[idx].Tls[0] {
				//log.Debug("-------------------- outflow %v", idx)
				err := self.module[UserIdent].conns[idx].forward(b)
				if err != nil {
					self.module[UserIdent].conns[idx].Close()
					//log.Debug("--------- forward write buf err %v", err)
				}
			}
		}
	}
	return nil
}

// NewServer return a new server object.
func NewServer(conf jsonhelper.JSONObject) *Server {
	if n := conf.GetAsObject("u").GetAsInt("max_module"); n != 0 {
		MaxModule = uint32(n)
	}
	if n := conf.GetAsObject("u").GetAsInt("packet_id_section"); n != 0 {
		PacketIDSection = uint32(n)
	}
	if n := conf.GetAsObject("u").GetAsInt("packet_reserved"); n != 0 {
		PacketIDReserved = uint32(n)
	}
	if n := conf.GetAsObject("u").GetAsInt("packet_id_count"); n != 0 {
		PacketIDCount = uint32(n)
	}
	cache := conf.GetAsObject("u").GetAsInt("packet_cache") //service
	log.Debug(">>>>>>>>> NewServer %v", cache)
	s := &Server{
		conf:   conf,
		module: make([]*Module, MaxModule, MaxModule),
		Input:  make(chan *Call, cache),
	}

	for k, v := range conf {
		if strings.HasPrefix(k, "listen_") {
			it := newModule(s, k, jsonhelper.JSONValueToObject(v))
			it.delegate = NewListener(it)
			s.module[it.ident] = it
		} else if strings.HasPrefix(k, "connect_") {
			it := newModule(s, k, jsonhelper.JSONValueToObject(v))
			it.delegate = NewConnector(it)
			s.module[it.ident] = it
		}
	}
	return s
}
