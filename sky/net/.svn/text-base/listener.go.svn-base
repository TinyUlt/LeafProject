package net

import (
	met "net"
	"sky/log"
	"time"
)

type Listener struct {
	*Module
	listener []met.Listener
}

func (self *Listener) Start() {
	var err error
	addrs := self.conf.GetAsArray("listen")
	for i := 0; i < addrs.Len(); i++ {
		it := self.conf.GetAsArray("listen").GetAsObject(i)
		self.listener[i], err = met.Listen(it.GetAsString("net"), it.GetAsString("addr"))
		if err != nil {
			panic(err)
		}
		go self.accept(self.listener[i])
		log.Trace("%s at %s", self.name, addrs.GetAsObject(i).GetAsString("addr"))
	}
}

func (self *Listener) Stop() {
	self.stoped = true
	addrs := self.conf.GetAsArray("listen")
	for i := 0; i < len(self.listener); i++ {
		if self.listener[i] != nil {
			self.listener[i].Close()
			log.Trace("%s stoped listener %s", self.name, addrs.GetAsObject(i).GetAsString("addr"))
		}
	}
	//close(self)
	for i := 0; i < len(self.conns); i++ {
		if self.conns[i].IsValid() {
			self.conns[i].Close()
		}
	}
	log.Trace("Listener Stoped<%s>", self.name)
}

func (self *Listener) accept(l met.Listener) {
	var idx uint32
	delay := time.Millisecond * 5

	for {
		select {
		case idx = <-self.free:
		default:
			if idx = uint32(len(self.conns)); idx < self.maxConnection {
				self.conns = append(self.conns, NewConn(idx, self.Module))
			} else {
				idx = <-self.free
			}
		}

		c, err := l.Accept()
		if err != nil {
			if self.stoped {
				break
			}
			if ne, ok := err.(met.Error); ok && ne.Temporary() {
				if delay *= 2; delay > time.Second {
					delay = time.Second
				}
				log.Error("%s: %v; retrying in %v", self.name, err, delay)
				time.Sleep(delay)
				continue
			}
			log.Critical("%s: %v", self.name, err)
			break
		}
		//log.Trace("%s accept a connect<%s>.", self.name, c.RemoteAddr())

		go self.conns[idx].work(c)
	}
}

func NewListener(module *Module) *Listener {
	addrs := module.conf.GetAsArray("listen")
	l := &Listener{
		Module:   module,
		listener: make([]met.Listener, addrs.Len()),
	}

	l.free = make(chan uint32, module.optimalConnection)
	l.conns = make([]*Conn, module.optimalConnection, module.optimalConnection)
	for i := uint32(0); i < module.optimalConnection; i++ {
		l.conns[i] = NewConn(i, module)
		l.free <- i
	}
	return l
}
