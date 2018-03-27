package net

import (
	met "net"
	"sky/log"
	"time"
)

type Connector struct {
	*Module
}

func (self *Connector) Start() {
	self.stoped = false
	addrs := self.conf.GetAsArray("connect")
	for i := 0; i < addrs.Len(); i++ {
		if addrs.GetAsObject(i).GetAsString("addr") == "" {
			continue
		}
		go self.connect(i, addrs.GetAsObject(i).GetAsString("net"),
			addrs.GetAsObject(i).GetAsString("addr"))
	}
	if self.conf.GetAsBool("auto_reconnect") {
		go self.watch()
	}
}

func (self *Connector) Stop() {
	self.stoped = true
	for i := 0; i < len(self.conns); i++ {
		if self.conns[i].IsValid() {
			self.conns[i].Close()
		}
	}
	log.Trace("Connector Stoped <%s>", self.name)
}

//todo: deputy of temp
func (self *Connector) ConnectTo(tp, addr string) {
	addrs := self.conf.GetAsArray("connect")
	for i := 0; i < addrs.Len(); i++ {
		it := addrs.GetAsObject(i)
		if it.GetAsString("addr") != "" {
			continue
		}
		it.Set("net", tp)
		it.Set("addr", addr)
		self.connect(i, tp, addr)
		break
	}
}

func (self *Connector) connect(idx int, net, addr string) {
	log.Trace("%s connect to %s", self.name, addr)
	retry := self.conf.GetAsInt("retry_times")
	if retry == 0 && self.conf.GetAsBool("auto_reconnect") {
		retry = 1 << 30
	}
	for {
		c, err := met.Dial(net, addr)
		if err != nil {
			if retry--; retry < 0 {
				log.Error("%v", err)
				return
			}
			log.Trace("%v", err)
			time.Sleep(time.Second * 2)
			continue
		}
		log.Trace("%s succed connect<%s>.", self.name, addr)
		self.conns[idx].Addr = addr
		go self.conns[idx].work(c)
		break
	}
}

func (self *Connector) watch() {
	for {
		idx := <-self.free
		if self.stoped {
			break
		}
		c := self.conf.GetAsArray("connect").GetAsObject(int(idx))
		log.Debug(">>>>>>> watch %v", idx)
		log.Trace("%s auto reconnect<%s>.", self.name, c.GetAsString("addr"))
		go self.connect(int(idx), c.GetAsString("net"), c.GetAsString("addr"))
	}
}

func NewConnector(module *Module) *Connector {
	c := &Connector{
		Module: module,
	}

	num := c.conf.GetAsArray("connect").Len()
	c.free = make(chan uint32, num)
	c.conns = make([]*Conn, num, num)
	for i := 0; i < num; i++ {
		c.conns[i] = NewConn(uint32(i), module)
	}
	return c
}
