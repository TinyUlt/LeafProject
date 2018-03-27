package net

import (
	"reflect"
	"sky/jsonhelper"
	"sky/log"
)

type Delegate interface {
	Start()
	Stop()
}

type Module struct {
	*Node
	delegate Delegate
	ident    uint32
	name     string
	ns       *Server
	proc     []*Proc
	conns    []*Conn
	free     chan uint32
	sync     bool
	ep       bool //endpoint
	stoped   bool
	from     uint32
	to       uint32
	prealloc uint32
	conf     jsonhelper.JSONObject
}

func (self *Module) register(id uint32, rcv, fun reflect.Value, rt reflect.Type) {
	log.Debug("register %v %v %v ", id, self.from, self.to)
	if id > PacketIDReserved {
		id -= self.from
	}
	//log.Debug(">>>>>>>>> register %v %v", id, self.from)
	if self.proc[id] != nil {
		log.Error("duplicated packet[%d][%v] [%v][%v]", id, rt, rcv, fun)
		panic("duplicated packet register")
	}
	self.proc[id] = &Proc{
		recver:   rcv,
		function: fun,
		packetRT: rt,
		nfree:    self.prealloc,
	}

	for i := uint32(0); i < self.prealloc; i++ {
		c := &Call{
			packet: reflect.New(self.proc[id].packetRT),
			next:   self.proc[id].free,
		}
		self.proc[id].free = c
	}
}

//var connType  = reflect.TypeOf(Conn{})
func (self *Module) allocCall(id uint32) (c *Call) {
	if id > PacketIDReserved {
		id -= self.from
	}
	if id < self.to && self.proc[id] != nil {
		self.proc[id].Lock()
		defer self.proc[id].Unlock()
		if c = self.proc[id].free; c != nil {
			self.proc[id].free = c.next
			self.proc[id].nfree--
		} else {
			c = &Call{packet: reflect.New(self.proc[id].packetRT)}
		}

	}
	return
}

func (self *Module) freeCall(id uint32, c *Call) {
	if self.proc[id].nfree < self.prealloc {
		self.proc[id].Lock()
		defer self.proc[id].Unlock()
		c.next = self.proc[id].free
		self.proc[id].free = c
		self.proc[id].nfree++

	}
}

func (self *Module) onConnected(conn *Conn) {
	//log.Debug(">>>>>>> module onConnected start")
	if call := self.allocCall(0); call != nil {
		call.conn = reflect.ValueOf(conn)
		//Note: this event may precede send/recv processed
		//so that in OnConnected call the send data cannot be greater than the send buffer size
		if self.sync {
			self.proc[0].function.Call([]reflect.Value{self.proc[0].recver, call.conn, call.packet})
			self.freeCall(0, call)
		} else {
			self.ns.Input <- call
		}
	}
	//log.Debug(">>>>>>> module onConnected end")
}

func (self *Module) onDisconnected(conn *Conn) {
	if call := self.allocCall(1); call != nil {
		call.conn = reflect.ValueOf(conn)
		//Note: this event may precede send/recv processed
		//so that in OnConnected call the send data cannot be greater than the send buffer size
		if self.sync {
			self.proc[1].function.Call([]reflect.Value{self.proc[1].recver, call.conn, call.packet})
			self.freeCall(1, call)
		} else {
			self.ns.Input <- call
		}
	}
}

func (self *Module) HashConn(id uint32) *Conn {
	//todo: idx may be a negative value in 32bit system
	idx := int(id) % len(self.conns)
	for i := idx; i > -1; i-- {
		if self.conns[i] != nil && self.conns[i].IsValid() {
			return self.conns[i]
		}
	}
	for i := idx + 1; i < len(self.conns); i++ {
		if self.conns[i] != nil && self.conns[i].IsValid() {
			return self.conns[i]
		}
	}
	return nil
}

func (self *Module) Retrieve(idx uint32) *Conn {
	if idx < uint32(len(self.conns)) {
		if self.conns[idx].IsValid() {
			return self.conns[idx]
		}
	}
	return nil
}
func (self *Module) RetrieveByUsrID(usrid uint32) *Conn {
	for _, _conn := range self.conns {
		if _conn.IsValid() && _conn.Tls[0] == usrid {
			return _conn
		}
	}
	return nil
}
func (self *Module) ToConnector() *Connector {
	if self.delegate != nil {
		if c, ok := self.delegate.(*Connector); ok {
			return c
		}
	}
	return nil
}

func (self *Module) Connections() []*Conn {
	return self.conns
}

func newModule(s *Server, key string, conf jsonhelper.JSONObject) *Module {
	module := &Module{
		ns:       s,
		name:     key,
		conf:     conf,
		sync:     conf.GetAsBool("sync"),
		Node:     newNode(conf),
		ep:       conf.GetAsString("codec") == "",
		ident:    uint32(conf.GetAsInt("ident")),
		prealloc: uint32(conf.GetAsInt("packet_prealloc")),
	}
	if ids := conf.GetAsArray("packet_id_range"); ids.Len() > 1 {
		module.from = uint32(ids.GetAsInt(0))
		module.to = uint32(ids.GetAsInt(1)) - module.from
	} else {
		module.from = module.ident * PacketIDSection
		module.to = PacketIDCount
	}
	module.proc = make([]*Proc, module.to, module.to)
	return module
}
