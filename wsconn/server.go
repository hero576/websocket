package wsconn

import (
	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

type WSServer interface {
	SeverListener()
	WsHandle(http.ResponseWriter, *http.Request)
}

type WebsocketServer struct {
	addr string
}

func NewWebsocketServer(addr string) *WebsocketServer {
	ws := new(WebsocketServer)
	ws.addr = addr
	return ws
}

func (w *WebsocketServer) SeverListener() {
	http.HandleFunc("/ws", w.WsHandle)
	if err := http.ListenAndServe(w.addr, nil); err != nil {
		log.Panic(err)
	}
}

func (w *WebsocketServer) WsHandle(writer http.ResponseWriter, request *http.Request) {
	//writer.Write([]byte("hello world"))

	var (
		data     []byte
		err      error
		upgrader = websocket.Upgrader{
			//允许跨域
			CheckOrigin: func(r *http.Request) bool {
				return true
			}}
		wsConn *websocket.Conn
		conn   *Connection
	)

	//upgrade websocket
	wsConn, err = upgrader.Upgrade(writer, request, nil)
	defer func() {
		conn.Close()
	}()
	if err != nil {
		if _, err := writer.Write([]byte(err.Error())); err != nil {
			logs.Error(err)
		}
	}

	if conn, err = initConnection(wsConn); err != nil {
		logs.Error("init websocket connect err", err)
		goto ERR
	}

	go func() {
		var (
			err error
		)
		for {
			if err = conn.WriteMessage([]byte("heartbeat")); err != nil {
				logs.Error("send heartbeat err", err)
				return
			}
			time.Sleep(1 * time.Second)
		}
	}()

	for {
		if data, err = conn.ReadMessage(); err != nil {
			goto ERR
		}
		if err = conn.WriteMessage(data); err != nil {
			goto ERR
		}
	}

ERR:
	conn.Close()
}
