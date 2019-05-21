package main

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"net"
)

var keyGUID = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"

const (
	// Fincode can record if the message is end or not end
	Fincode = 1 << 7
	// maskBit can record if message is mask or not mask
	maskBit = 1 << 7
	// TestMessage record if the file if test or not test
	TestMessage = 1
	// CloseMessage record if the message if close or not close
	CloseMessage = 8
)

// Conn can recode the  message of head
type Conn struct {
	writeBuf []byte
	maskKey  [4]byte

	conn net.Conn
}

// SendDate can send message from server
func (c *Conn) SendDate(data []byte) {

	reallsize := len(data)
	c.writeBuf = make([]byte, 10+reallsize)
	c.writeBuf[0] = Fincode & (byte)(TestMessage)
	dataLength := 2

	switch {
	case reallsize >= 127:
		c.writeBuf[1] = (byte)(127)
		binary.BigEndian.PutUint64(c.writeBuf[dataLength:], uint64(reallsize))
		dataLength += 8
	case reallsize >= 126:
		c.writeBuf[1] = (byte)(126)
		binary.BigEndian.PutUint16(c.writeBuf[dataLength:], uint16(reallsize))
		dataLength += 2
	default:
		c.writeBuf[1] = (byte)(reallsize)
	}
	copy(c.writeBuf[reallsize:], data)
	c.conn.Write(c.writeBuf)
}

// RecvData can recv message from client
func RecvData() {

}

// KeyAndSecToSha1 function  can get a hash of key and  sec-websocket-accept
func KeyAndSecToSha1(key string) string {
	h := sha1.New()
	key += keyGUID
	h.Write([]byte(key))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// maskBytes function can deal with byte to mask
func maskBytes(key [4]byte, pos int, b []byte) int {
	for i := range b {
		b[i] ^= key[pos&3] //a%(2^n) == a &(2^n-1)
		pos++
	}
	return pos & 3
}

func main() {

}
