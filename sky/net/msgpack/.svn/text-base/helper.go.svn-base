package msgpack

import (
	"reflect"
	"fmt"
	"time"
	"errors"
)

type ContainerType byte

const (
	ContainerRawBytes = ContainerType('b')
	ContainerList = ContainerType('a')
	ContainerMap = ContainerType('m')
)

var (
	nilIntfSlice = []interface{}(nil)
	intfSliceTyp = reflect.TypeOf(nilIntfSlice)
	intfTyp = intfSliceTyp.Elem()
	byteSliceTyp = reflect.TypeOf([]byte(nil))
	timeTyp = reflect.TypeOf(time.Time{})
	mapStringIntfTyp = reflect.TypeOf(map[string]interface{}(nil))
	mapIntfIntfTyp = reflect.TypeOf(map[interface{}]interface{}(nil))
)

func getContainerByteDesc(ct ContainerType) (cutoff int, b0, b1, b2 byte) {
	switch ct {
	case ContainerRawBytes:
		cutoff = 32
		b0, b1, b2 = 0xa0, 0xda, 0xdb
	case ContainerList:
		cutoff = 16
		b0, b1, b2 = 0x90, 0xdc, 0xdd
	case ContainerMap:
		cutoff = 16
		b0, b1, b2 = 0x80, 0xde, 0xdf
	default:
		panic(fmt.Errorf("getContainerByteDesc: Unknown container type: %v", ct))
	}
	return
}

func reflectValue(v interface{}) (rv reflect.Value) {
	rv, ok := v.(reflect.Value)
	if !ok {
		rv = reflect.ValueOf(v)
	}
	return 
}

func panicToErr(err *error) {
	if x := recover(); x != nil { 
		switch xerr := x.(type) {
		case error:
			*err = xerr
		case string:
			*err = errors.New(xerr)
		default:
			*err = fmt.Errorf("%v", x)
		}
	}
}

