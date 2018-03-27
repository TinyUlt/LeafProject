package net

import (
	"encoding/binary"
	met "net"
	"sky/log"
	"sync"
	"time"
)

//Error Code
const (
	RUNING  = 0x00
	RCLOSED = 0x01
	WCLOSED = 0x02
	CLOSED  = RCLOSED | WCLOSED
)

type Conn struct {
	met.Conn
	sync.Mutex
	module *Module
	rbuff  *Buffer
	wbuff  *Buffer
	wcbuff *Buffer
	codec  Codec
	ch     chan uint32
	state  uint32
	Idx    uint32
	Tls    [2]uint32 //[0]=>id, [1]=>index
	Addr   string
}

func (self *Conn) Send(packet Packet) error {
	self.Lock()
	defer self.Unlock()
	if self.state != RUNING {
		//self.Unlock()
		return ClosedError
	}

	if self.wbuff.Len() > self.module.sendAvailableLen {
		log.Error(BlockError)
		//self.Unlock()
		return BlockError
	}

	if err := self.codec.Encode(packet); err != nil {
		log.Error("encode error:", err)
		//self.Unlock()
		return err
	}
	//self.Unlock()

	select {
	case self.ch <- 1:
	default:
	}
	return nil
}

func (self *Conn) SendTo(id, idx uint32, packet Packet) error {
	self.Lock()
	defer self.Unlock()

	self.Tls[0] = id
	self.Tls[1] = idx

	if self.state != RUNING {
		//self.Unlock()
		return ClosedError
	}

	if self.wbuff.Len() > self.module.sendAvailableLen {
		log.Error(BlockError)
		//self.Unlock()
		return BlockError
	}

	if err := self.codec.Encode(packet); err != nil {
		log.Error(err)
		//self.Unlock()
		return err
	}
	//self.Unlock()

	select {
	case self.ch <- 1:
	default:
	}
	return nil
}

func (self *Conn) Close() error {
	//log.Debug(">>>>>>>>>> Close")
	push(self.ch, 0)
	return nil
}
func (self *Conn) CloseEx() error {
	//log.Debug(">>>>>>>>>> CloseEx")
	pushex(self.ch, 0)
	return nil
}
func (self *Conn) SafeClose() {
	//log.Debug(">>>>>>>>>> SafeClose")
	self.Lock()
	if self.state == RUNING {
		self.ch <- 0
	}
	self.Unlock()
}

func (self *Conn) IsValid() bool {
	return self.state == RUNING
}

func (self *Conn) RemoteIp() (ip string) {
	if self.IsValid() {
		ip, _, _ = met.SplitHostPort(self.RemoteAddr().String())
	}
	return
}

func (self *Conn) closed(op uint32) {
	//panic(nil)
	//log.Trace("closed by op:%d", op)
	self.Lock()
	self.state |= op
	if self.state == CLOSED {
		self.module.onDisconnected(self)
		//log.Debug(">>>>>>>>>> conn closed %v", self.Idx)
		self.module.free <- self.Idx
	} else {
		self.Conn.Close()
		push(self.ch, 0)
	}
	self.Unlock()
}

func (self *Conn) work(c met.Conn) {
	self.Conn = c
	//log.Debug(">>>>>>>> conn work start %v", self.Idx)
	if tc, ok := c.(*met.TCPConn); ok {
		tc.SetLinger(self.module.lingerTime)
		tc.SetReadBuffer(self.module.readBuffLen)
		tc.SetWriteBuffer(self.module.writeBuffLen)
	}
	//log.Debug(">>>>>>>> conn work 1111")
	self.rbuff.Seek(0, 0)
	self.wbuff.Seek(0, 0)
	pop(self.ch)
	//log.Debug(">>>>>>>> conn work 2222")
	self.state = RUNING
	self.module.onConnected(self)

	//log.Debug(">>>>>>>> conn work 3333")
	go self.send()
	self.recv()
	//log.Debug(">>>>>>>> conn work end")
}

func (self *Conn) send() {
	var (
		off, end, n int
		buf         []byte
		ok          uint32
		err         error
	)
	for {
		if ok = <-self.ch; ok < 1 {
			err = ClosedError
			break
		}

		end = self.wbuff.Len()
		buf = self.wbuff.Bytes()[off:end]
		if n = end - off; n == 0 {
			//self.Unlock()
			continue
		}

		if n, err = self.Write(buf); err != nil {
			//self.Unlock()
			break
		}
		//log.Debug("Conn %d Write %d", self.module.ident, n)
		self.Lock()
		end = self.wbuff.Len()
		if off += n; off == end {
			off = 0
			self.wbuff.Seek(0, 0)
		} else if off > self.module.sendFixedLen {
			copy(self.wbuff.Bytes(), self.wbuff.Bytes()[off:end])
			self.wbuff.Seek(int64(end-off), 0)
			//log.Debug("Moved write buffer [%d-%d]", off, end)
			off = 0
		}
		self.Unlock()
	}
	self.closed(WCLOSED)
}

func (self *Conn) recv() {
	var (
		err              error
		off, dif, acc, n int
	)

	buf := self.rbuff.Bytes()
	end := self.module.recvAvailableLen
	unit := self.module.minPacketLen
Loop:
	for {
		if self.module.timeout != 0 {
			self.SetReadDeadline(time.Now().Add(self.module.timeout))
		}
		if n, err = self.Read(buf[off:end]); err != nil {
			//log.Debug(">>>>> Read err %v %v %v", err, off, end)
			break
		}
		//log.Debug("Conn %d Reads %d", self.module.ident, n)
		//log.Debug("data: %v ", buf[off:off+n])
		if off += n; off < (acc + unit) {
			//log.Debug("pending: %d-%d-%d", off, acc, unit)
			continue
		}

		for {
			unit = int(binary.BigEndian.Uint32(buf[acc:]))
			//id := int(binary.BigEndian.Uint32(buf[acc+4:]))
			//if id == 317 {
			//log.Debug(">>>>>> %v", id)
			//}
			if unit < self.module.minPacketLen || unit > self.module.maxPacketLen {
				err = HeaderError
				log.Debug("%s %d %d %d", self.module.name, len(buf), off, end)
				log.Debug("%v, %d, %d", HeaderError, unit, self.module.maxPacketLen)
				//log.Debug("%d, %d, %v", acc, self.module.minPacketLen, self)
				break Loop
			}

			if dif = unit - (off - acc); dif > 0 {
				end = off + dif
				//log.Debug(">>>>> dif > 0   %v %v  %v", unit, off-acc, dif)
				break
			}

			if err = self.codec.Decode(uint32(unit)); err != nil {
				//log.Debug("%s decode:%v", self.module.name, err)
				break Loop
			}

			if acc += unit; off == acc {
				off, acc = 0, 0
				end = self.module.recvAvailableLen
				unit = self.module.minPacketLen
				self.rbuff.Seek(0, 0)
				//log.Debug(">>>>>>>>> recv 55555 %v %v %v %v", unit, off, acc, self)
				break
			} else if (off - acc) < self.module.minPacketLen {
				copy(buf, buf[acc:off])
				self.rbuff.Seek(0, 0)
				off -= acc
				acc = 0
				end = self.module.recvAvailableLen
				unit = self.module.minPacketLen
				//log.Debug(">>>>>>>>> recv 66666 %v %v %v", off, acc, self.module.minPacketLen)
				break
			}
		}
	}
	self.closed(RCLOSED)
}

func (self *Conn) forward(b []byte) error {
	self.Lock()
	defer self.Unlock()
	var err error = nil
	if self.state == RUNING {
		//log.Debug(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>> forward %v", self.Idx)
		if _, err = self.wbuff.Write(b); err != nil {
			//log.Critical("buffer overflow!!!!!!!!!!!! %v", self.Idx)
		} else {
			select {
			case self.ch <- 1:
			default:
			}
		}
	}
	//	self.Unlock()
	return err
}

func NewConn(idx uint32, module *Module) *Conn {
	c := &Conn{
		Idx:    idx,
		state:  CLOSED,
		module: module,
		ch:     make(chan uint32, 1),
		rbuff:  NewBuffer(make([]byte, module.recvBuffLen, module.recvBuffLen)),
		wbuff:  NewBuffer(make([]byte, module.sendBuffLen, module.sendBuffLen)),
		wcbuff: NewBuffer(make([]byte, module.maxPacketLen, module.maxPacketLen)),
	}

	c.codec = newCodec(c)
	//log.Debug("New connection:%d. %v", idx, c.state)
	return c
}
