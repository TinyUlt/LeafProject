package net

import (
	"encoding/binary"
	"reflect"
	"sky/log"
	"sky/net/msgpack"
	"time"
)

type Codec interface {
	Encode(packet Packet) error
	Decode(n uint32) error
}

//server oriented
type codec struct {
	conn *Conn
	enc  *msgpack.Encoder
	dec  *msgpack.Decoder
}

func (self *codec) Encode(packet Packet) (err error) {
	//pos := self.conn.wbuff.Len()
	buf := self.conn.wcbuff.Bytes()[:]
	self.conn.wcbuff.Seek(16, 0)
	if err = self.enc.Encode(packet); err != nil {
		//self.conn.wbuff.Seek(int64(pos), 0)
		//log.Debug(">>>>>>>>?????????? Encode err %v", err)
		return
	}
	sz := self.conn.wcbuff.Len() //- pos
	//max_packet_len := 102400
	//if sz <= 0 || sz > max_packet_len {
	//log.Debug("==============++++++++++++++++++++??????????????????????????????? codec Encode %v", sz)
	//panic("codec Encode err")
	//}

	binary.BigEndian.PutUint32(buf, uint32(sz))
	binary.BigEndian.PutUint32(buf[4:], packet.PacketID())
	binary.BigEndian.PutUint32(buf[8:], self.conn.Tls[0])
	binary.BigEndian.PutUint32(buf[12:], self.conn.Tls[1])
	self.conn.wbuff.Write(buf[:sz])
	//if packet.PacketID() == 317 {
	//	log.Debug(">>>>>>>>>> %v", packet.PacketID())
	//	}
	//encrypt here
	return
}

func (self *codec) Decode(n uint32) error {
	//decrypt here
	pos := uint32(self.conn.rbuff.Len())
	buf := self.conn.rbuff.Bytes()[pos:]
	id := binary.BigEndian.Uint32(buf[4:])
	call := self.conn.module.allocCall(id)

	//log.Debug("===== >>>>>>>>>>>>> codec Decode %v", id)
	if call == nil {
		//recovery and skp this packet, not return an error
		self.conn.rbuff.Seek(int64(n), 1)
		log.Debug("packet not register<id:%d len:%d>", id, n)
		return nil
	}
	call.tls[0] = binary.BigEndian.Uint32(buf[8:])
	call.tls[1] = binary.BigEndian.Uint32(buf[12:])
	//log.Debug("read tls:%d %d", call.tls[0], call.tls[1])
	self.conn.rbuff.Seek(16, 1)
	if err := self.dec.DecodeValue(call.packet); err != nil {
		//recovery and skp this packet, not return an error
		log.Debug("===== Decode err %v %v", id, err)
		self.conn.rbuff.Seek(int64(pos+n), 0)
		return err
	}

	call.conn = reflect.ValueOf(self.conn)
	self.conn.module.ns.enqueue(id, self.conn, call)
	return nil
}

//client oriented
type firewall struct {
	conn  *Conn
	enc   *msgpack.Encoder
	dec   *msgpack.Decoder
	count int
	prev  time.Time
}

/*func (self *firewall) Encode(packet Packet) (err error) {
	pos := self.conn.wbuff.Len()
	buf := self.conn.wbuff.Bytes()[pos:]
	self.conn.wbuff.Seek(16, 1)
	if err = self.enc.Encode(packet); err != nil {
		self.conn.wbuff.Seek(int64(pos), 0)
		return
	}
	sz := self.conn.wbuff.Len() - pos
	binary.BigEndian.PutUint32(buf, uint32(sz))
	binary.BigEndian.PutUint32(buf[4:], packet.PacketID())
	//encrypt here
	return
}*/
func (self *firewall) Encode(packet Packet) (err error) {
	//pos := self.conn.wbuff.Len()
	buf := self.conn.wcbuff.Bytes()[:]
	self.conn.wcbuff.Seek(16, 0)
	if err = self.enc.Encode(packet); err != nil {
		//self.conn.wbuff.Seek(int64(pos), 0)
		return
	}
	sz := self.conn.wcbuff.Len() //- pos
	binary.BigEndian.PutUint32(buf, uint32(sz))
	binary.BigEndian.PutUint32(buf[4:], packet.PacketID())
	self.conn.wbuff.Write(buf[:sz])
	//encrypt here
	return
}
func (self *firewall) Decode(n uint32) error {
	self.count++
	if self.count > self.conn.module.flowLimit {
		if time.Since(self.prev) < 30*time.Second {
			log.Debug("%s [id:%d, ip:%s] %v", TrafficError,
				self.conn.Tls[0], self.conn.RemoteAddr().String(), self.count)
			return TrafficError
		}
		self.count = 0
		self.prev = time.Now()
	}
	//decrypt here
	pos := uint32(self.conn.rbuff.Len())
	buf := self.conn.rbuff.Bytes()[pos:]
	id := binary.BigEndian.Uint32(buf[4:])
	//log.Debug("recv packet<id:%d len:%d> %v %v", id, n, self.conn.Tls[0], PacketIDReserved)
	if self.conn.Tls[0] > 0 && id > PacketIDReserved {
		ident := id / PacketIDSection
		idx := binary.BigEndian.Uint32(buf[12:])
		//log.Debug("write tls:%d %d", self.conn.Tls[0], self.conn.Idx)
		binary.BigEndian.PutUint32(buf[8:], self.conn.Tls[0])
		binary.BigEndian.PutUint32(buf[12:], self.conn.Idx)
		self.conn.rbuff.Seek(int64(n), 1)
		//log.Debug(">>>>>>>>>>>>>> %v %v %v %v", id, ident, idx, GameIdent)
		//	log.Debug(">>>> firewall Decode id[%v] route[%v] %v", id, idx, ident == GameIdent)
		switch ident {
		case UserIdent:
			return UnlawfulError
		case GameIdent:
			self.conn.Tls[1] = idx
			//log.Debug(">>>>>>>>> %v %v %v", id, ident, GameIdent)
			return self.conn.module.ns.inflow(idx, buf[:n])
		}
		return self.conn.module.ns.hashflow(ident, self.conn.Tls[0], buf[:n])
	}
	call := self.conn.module.allocCall(id)
	if call == nil {
		//recovery and skp this packet, not return an error
		self.conn.rbuff.Seek(int64(n), 1)

		log.Debug("packet not register<id:%d len:%d> %v %v", id, n, self.conn.Tls[0], self.conn.state)
		return nil
	}
	self.conn.rbuff.Seek(16, 1)
	if err := self.dec.DecodeValue(call.packet); err != nil {
		//recovery and skp this packet, not return an error
		self.conn.rbuff.Seek(int64(pos+n), 0)
		return nil
	}

	call.conn = reflect.ValueOf(self.conn)
	self.conn.module.ns.enqueue(id, self.conn, call)
	return nil
}

//gateway oriented
type route struct {
	conn *Conn
	enc  *msgpack.Encoder
	dec  *msgpack.Decoder
}

/*
func (self *route) Encode(packet Packet) (err error) {
	pos := self.conn.wbuff.Len()
	buf := self.conn.wbuff.Bytes()[pos:]
	self.conn.wbuff.Seek(16, 1)
	if err = self.enc.Encode(packet); err != nil {
		self.conn.wbuff.Seek(int64(pos), 0)
		return
	}
	sz := self.conn.wbuff.Len() - pos
	binary.BigEndian.PutUint32(buf, uint32(sz))
	binary.BigEndian.PutUint32(buf[4:], packet.PacketID())
	binary.BigEndian.PutUint32(buf[8:], 0)
	binary.BigEndian.PutUint32(buf[12:], 0)
	//encrypt here
	return
}*/
func (self *route) Encode(packet Packet) (err error) {
	//pos := self.conn.wbuff.Len()
	buf := self.conn.wcbuff.Bytes()[:]
	self.conn.wcbuff.Seek(16, 0)
	if err = self.enc.Encode(packet); err != nil {
		//self.conn.wbuff.Seek(int64(pos), 0)
		return
	}
	sz := self.conn.wcbuff.Len() //- pos
	binary.BigEndian.PutUint32(buf, uint32(sz))
	binary.BigEndian.PutUint32(buf[4:], packet.PacketID())
	binary.BigEndian.PutUint32(buf[8:], 0)
	binary.BigEndian.PutUint32(buf[12:], 0)
	self.conn.wbuff.Write(buf[:sz])
	//encrypt here
	return
}
func (self *route) Decode(n uint32) error {
	//decrypt here
	pos := uint32(self.conn.rbuff.Len())
	buf := self.conn.rbuff.Bytes()[pos:]
	id := binary.BigEndian.Uint32(buf[4:])
	user := binary.BigEndian.Uint32(buf[8:])
	//log.Debug("===== route Decode 1111 %v", id)
	if id > PacketIDReserved && self.conn.Tls[0] > 0 && user > 0 {
		idx := binary.BigEndian.Uint32(buf[12:])
		//log.Debug(">>>>>>>>>>>>==== Decode id[%v] route[%v]", id, idx)
		self.conn.module.ns.outflow(user, idx, buf[:n])
		self.conn.rbuff.Seek(int64(n), 1)
		return nil
	}
	call := self.conn.module.allocCall(id)
	if call == nil {
		//recovery and skp this packet, not return an error
		self.conn.rbuff.Seek(int64(n), 1)
		log.Debug("packet not register<id:%d len:%d>", id, n)
		return nil
	}
	self.conn.rbuff.Seek(16, 1)
	if err := self.dec.DecodeValue(call.packet); err != nil {
		//recovery and skp this packet, not return an error
		self.conn.rbuff.Seek(int64(pos+n), 0)
		log.Debug("[%d] %s", id, err)
		return nil
	}

	call.conn = reflect.ValueOf(self.conn)
	self.conn.module.ns.enqueue(id, self.conn, call)
	return nil
}

func newCodec(conn *Conn) Codec {
	switch conn.module.codec {
	case "firewall":
		return &firewall{
			conn: conn,
			enc:  msgpack.NewEncoder(conn.wcbuff),
			dec:  msgpack.NewDecoder(conn.rbuff, nil),
		}
	case "route":
		return &route{
			conn: conn,
			enc:  msgpack.NewEncoder(conn.wcbuff),
			dec:  msgpack.NewDecoder(conn.rbuff, nil),
		}
	}
	return &codec{
		conn: conn,
		enc:  msgpack.NewEncoder(conn.wcbuff),
		dec:  msgpack.NewDecoder(conn.rbuff, nil),
	}
}
