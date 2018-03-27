package net

import (
	"errors"
	"io"
)
 
// SocketBuffer, A Buffer object can only using read or write interface, 
// Using both the read and write interface will produce unpredictable errors
type Buffer struct {
	s 		[]byte
	i       int // current r/w index
}

func (self *Buffer) Bytes() []byte { return self.s }

// Len returns the number of bytes of the unread portion of the
// slice.
func (self *Buffer) Len() int {
	return self.i
}

func (self *Buffer) Read(b []byte) (n int, err error) {
	if len(b) == 0 {
		return 0, nil
	}
	if self.i >= len(self.s) {
		return 0, io.EOF
	}
	n = copy(b, self.s[self.i:])
	self.i += n
	return
}

func (self *Buffer) ReadByte() (b byte, err error) {
	if self.i >= len(self.s) {
		return 0, io.EOF
	}
	b = self.s[self.i]
	self.i++
	return
}

func (self *Buffer) Write(b []byte) (n int, err error) {
	if (self.i + len(b)) >= len(self.s) {
		return 0, io.ErrShortWrite
	}

	n = copy(self.s[self.i:], b)
	self.i += n	
	return
}


// Seek implements the io.Seeker interface.
func (self *Buffer) Seek(offset int64, whence int) (int64, error) {
	var abs int64
	switch whence {
	case 0:
		abs = offset
	case 1:
		abs = int64(self.i) + offset
	case 2:
		abs = int64(len(self.s)) + offset
	default:
		return 0, errors.New("bytes: invalid whence")
	}
	if abs < 0 {
		return 0, errors.New("bytes: negative position")
	}
	if abs >= 1<<31 {
		return 0, errors.New("bytes: position out of range")
	}
	self.i = int(abs)
	return abs, nil
}

// NewReader returns a new Buffer reading from b.
func NewBuffer(b []byte) *Buffer { return &Buffer{b, 0} }