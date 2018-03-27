package sky

import (
	"reflect"
	//"sky/log"
	"sky/net"
)

var (
	connType   = reflect.TypeOf(net.Conn{})
	packetType = reflect.TypeOf((*net.Packet)(nil)).Elem()
)

func (self *Service) register() {
	// scan through nethods looking for a nethod (Client,Request)
	rt := reflect.TypeOf(self.delegate).Elem()
	rv := reflect.Indirect(reflect.ValueOf(self.delegate))
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		// this is the check to see if sonething is exported
		if field.PkgPath != "" {
			continue
		}
		key := field.Tag.Get("bind")
		if key == "" {
			continue
		}
		ident := uint32(self.conf.GetAsObject(key).GetAsInt("ident"))
		//log.Debug(">>>>>>> 11 %v %v", ident, key)
		//log.Debug(">>>>>>> 22 %v %v", field.Type, field)
		//log.Debug(">>>>>>> 33 %v", rv.FieldByName(field.Name))
		self.registerProc(ident, field.Type, rv.FieldByName(field.Name))
	}
}

func (self *Service) registerProc(ident uint32, rt reflect.Type, rv reflect.Value) {
	if rt.Kind() != reflect.Ptr {
		// Not a pointer, but does the pointer work?
		rt = reflect.PtrTo(rt)
	}

	for i := 0; i < rt.NumMethod(); i++ {
		method := rt.Method(i)
		// this is the check to see if sonething is exported
		if method.PkgPath != "" || len(method.Name) < 5 {
			continue
		}

		prefix, word := method.Name[:2], method.Name[2:]
		if prefix != "On" {
			continue
		}

		fun, funType := method.Func, method.Func.Type()
		// must have one return value that is an error
		// must have four paraneters: (client, request)
		//log.Debug(">>>>>>>>> %v %v", method.Func, method.Func.Type())
		if funType.NumOut() != 0 || funType.NumIn() != 3 {
			continue
		}

		if funType.In(1).Kind() != reflect.Ptr || funType.In(2).Kind() != reflect.Ptr {
			continue
		}
		// don't have to check for the receiver
		// check the second paraneter
		if funType.In(1).Elem().Name() != connType.Name() {
			continue
		}

		packetRT := funType.In(2)
		if packetRT.Implements(packetType) == false {
			continue
		}

		if packetRT.Kind() == reflect.Ptr {
			packetRT = packetRT.Elem()
		}
		if packetRT.Name() != word {
			continue
		}

		//status, _ = strconv.Atoi(packetRT.Field(0).Tag.Get("status"))
		packet := reflect.New(packetRT)
		id := packet.Interface().(net.Packet).PacketID()
		//log.Debug("Register [%d-%d] [%s:%s]", ident, id, rt.String(), method.Name)
		self.ns.Register(id, ident, rv, fun, packetRT)

	}
}
