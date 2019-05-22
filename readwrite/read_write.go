package readwrite

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
)

// Conn can recode the  message of head
type Conn struct {
	writeBuf []byte
	maskKey  [4]byte
	conn     net.Conn
}

// Newconn can get a new object
func Newconn(cs net.Conn) *Conn {
	br := &Conn{}
	br.conn = cs
	return br
}

// SendData can send message from server
func (c *Conn) SendData(data []byte) {

	reallsize := len(data)
	c.writeBuf = make([]byte, 10+reallsize)
	c.writeBuf[0] = Fincode | (byte)(TestMessage)
	dataLength := 2

	switch {
	case reallsize >= 127:
		c.writeBuf[1] = byte(0x00) | 127
		binary.BigEndian.PutUint64(c.writeBuf[dataLength:], uint64(reallsize))
		dataLength += 8
	case reallsize >= 126:
		c.writeBuf[1] = byte(0x00) | 126
		binary.BigEndian.PutUint16(c.writeBuf[dataLength:], uint16(reallsize))
		dataLength += 2
	default:
		c.writeBuf[1] = byte(0x00) | byte(dataLength)
	}
	copy(c.writeBuf[reallsize:], data)
	c.conn.Write(c.writeBuf)
}

// ReadData can recv message from client
func (c *Conn) ReadData() ([]byte, error) {
	var buff [8]byte
	_, err := c.conn.Read(buff[0:2])
	if err != nil {
		return nil, fmt.Errorf("read data error:%v", err)
	}
	final := buff[0]&Fincode != 0
	if !final {
		log.Println("this is websocket don't support fragmented")
		return nil, fmt.Errorf("this is websocket don't support fragmented")
	}
	filetype := int(buff[0] & 0xf)
	if filetype == CloseMessage || filetype != TestMessage {
		log.Println("the websocket only support text file")
		return nil, fmt.Errorf("the websocket only support text file")
	}
	mask := buff[1]&(byte)(maskBit) != 0

	length := (uint64)(buff[1] & 0x7F)
	datalength := (uint64)(length)
	switch length {
	case 127:
		if _, err := c.conn.Read(buff[:2]); err != nil {
			return nil, fmt.Errorf("read data length error: %v", err)
		}
		datalength = uint64(binary.BigEndian.Uint16(buff[:2]))
	case 126:
		if _, err := c.conn.Read(buff[:8]); err != nil {
			return nil, fmt.Errorf("read data length error：%v", err)
		}
		datalength = uint64(binary.BigEndian.Uint64(buff[0:8]))
	}
	//读取真正的数据和mask效验阶段
	if mask {
		if _, err := c.conn.Read(c.maskKey[:]); err != nil {
			return nil, fmt.Errorf("read  mask error: %v", err)
		}

	}
	databuff := make([]byte, datalength)
	if _, err := c.conn.Read(databuff); err != nil {
		return nil, fmt.Errorf("read data  error: %v", err)
	}
	if mask {
		maskBytes(c.maskKey, 0, databuff)
	}
	return databuff, nil

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
func TokenListContainsValue(h http.Header, name string, value string) bool {
	for _, v := range h[name] {
		for _, s := range strings.Split(v, ",") {
			if strings.EqualFold(value, strings.TrimSpace(s)) {
				return true
			}
		}
	}
	return false
}
