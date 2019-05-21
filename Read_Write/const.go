package Read_Write

import "net"

var keyGUID = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"

// Conn can recode the  message of head
type Conn struct {
	writeBuf []byte
	maskKey  [4]byte

	conn net.Conn
}

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
