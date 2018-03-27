
package net

import "errors"

//sky net errors
var (
    ReadError    	= errors.New("connection read error")
    WriteError      = errors.New("connection write error")
    BlockError      = errors.New("connection block")
    BreakError    	= errors.New("connection break")
    ClosedError   	= errors.New("connection closed")
    ShakeError      = errors.New("hand shake error")
    ShakeCodeError  = errors.New("write channel closed")
    LossedError   	= errors.New("connection loss this packet")
    IDRangeError 	= errors.New("packet id range error")
    HeaderError     = errors.New("packet header error")
    EncodeError 	= errors.New("packet encode error")
    DecodeError 	= errors.New("packet decode error")
    UnregistError   = errors.New("packet not register")
    UnlawfulError   = errors.New("unlawful incursion")
    TrafficError    = errors.New("traffic overload")
    )