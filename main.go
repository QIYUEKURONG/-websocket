package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/QIYUEKURONG/websocket/readwrite"
)

var keyGUID = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"

func tokenListContainsValue(h http.Header, name string, value string) bool {
	for _, v := range h[name] {
		for _, s := range strings.Split(v, ",") {
			if strings.EqualFold(value, strings.TrimSpace(s)) {
				return true
			}
		}
	}
	return false
}

// Upgrade function can to upgrade protocol from http to websocket
func Upgrade(w http.ResponseWriter, r *http.Request) (*readwrite.Conn, error) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return nil, fmt.Errorf("the method not GET")
	}

	if values := r.Header["Sec-Websocket-Version"]; len(values) == 0 || values[0] != "13" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return nil, fmt.Errorf("the Sec-Websocket-Version != 13")
	}
	// 检查Connection和Upgrade
	if !(readwrite.TokenListContainsValue(r.Header, "Connection", "upgrade")) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return nil, fmt.Errorf("the Connection != upgrade")
	}
	if !readwrite.TokenListContainsValue(r.Header, "Upgrade", "websocket") {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return nil, fmt.Errorf("the Upgrade != websocket")
	}
	//计算Sec-websocket-Accept的值
	challengkey := r.Header.Get("Sec-Websocket-Key")
	if challengkey == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return nil, fmt.Errorf("the Sec-websocket-Accept if NULL")
	}
	var (
		conn net.Conn
		br   *bufio.Reader
	)

	h, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return nil, fmt.Errorf("websocket: response dose not implement http.Hijacker")
	}
	var rw *bufio.ReadWriter
	conn, rw, err := h.Hijack()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return nil, fmt.Errorf("websocket: get http.Hijacker error : %v", err)
	}
	br = rw.Reader
	if br.Buffered() > 0 {
		conn.Close()
		return nil, fmt.Errorf("websocket: client sent data before handshake is complete")
	}
	//完成后就返回response
	var p = []byte{}
	p = ([]byte)("HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Accept:")

	p = append(p, readwrite.KeyAndSecToSha1(challengkey)...)
	p = append(p, "\r\n\r\n"...)
	if _, err = conn.Write(p); err != nil {
		conn.Close()
		return nil, fmt.Errorf("net.conn to send message find error")
	}
	log.Println("Upgrade http to webcosket success")
	netconn := readwrite.Newconn(conn)
	return netconn, nil
}

// http的处理器
func echo(w http.ResponseWriter, r *http.Request) {

	conn, err := Upgrade(w, r)
	if err != nil {
		log.Printf("call Upgrade error :%v\n", err)
		return
	}

	for {
		message, err := conn.ReadData()
		if err != nil {
			log.Printf("read error: %v\n", err)
			break
		}
		log.Printf("recv: %s\n", message)
		conn.SendData(message)
	}

}
func home(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hello world")
	w.Write([]byte("hello world"))
}

func main() {
	log.SetFlags(0)
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/", home)
	http.ListenAndServe("127.0.0.1:8088", nil)

}
