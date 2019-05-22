package main

import (
	"bufio"
	"fmt"
	"html/template"
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
	fmt.Println("start to Upgrade")
	if r.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return nil, fmt.Errorf("the method not GET")
	}
	fmt.Println(r.Header)

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
	homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}

func main() {
	log.SetFlags(0)
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}

var homeTemplate = template.Must(template.New("").Parse(`
 <!DOCTYPE html>
 <head>
  <meta charset="utf8">
  <script>
  window.addEventListener("load", function(evt) {
 
    var output = document.getElementById("output");
     var input = document.getElementById("input");
     var ws;
 
    var print = function(message) {
        var d = document.createElement("div");
         d.innerHTML = message;
         output.appendChild(d);
     };
 
     document.getElementById("open").onclick = function(evt) {
         if (ws) {
             return false;
         }
		 ws = new WebSocket("{{.}}");
		 ws.onopen = function(evt) {
			             print("OPEN");
			       }
			        ws.onclose = function(evt) {
		           print("CLOSE");
		            ws = null;
			        }
			        ws.onmessage = function(evt) {
			            print("RESPONSE: " + evt.data);
		         }
			        ws.onerror = function(evt) {
			             print("ERROR: " + evt.data);
			        }
			        return false;
			    };
			
			   document.getElementById("send").onclick = function(evt) {
			         if (!ws) {
						return false;
						       }
						       print("SEND: " + input.value);
						        ws.send(input.value);
						        return false;
						     };
						
					     document.getElementById("close").onclick = function(evt) {
						         if (!ws) {
						             return false;
					         ws.close();
						         return false;
						     };
					   
 });
 </script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>
点击 "Open" 开始一个新的WebSocket链接,
“Send" 将内容发送到服务器，
 "Close" 将关闭链接。
 <p>
 <form>
 <button id="open">Open</button>
 <button id="close">Close</button>
<p><input id="input" type="text" value="hello world!">
 <button id="send">Send</button>
 </form>
 </td><td valign="top" width="50%">
 <div id="output"></div>
 </td></tr></table>
 </body>
 </html>
 `))
