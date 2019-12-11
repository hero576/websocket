package wsconn

import (
	"fmt"
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
		msgType  int
		data     []byte
		err      error
		upgrader = websocket.Upgrader{
			//允许跨域
			CheckOrigin: func(r *http.Request) bool {
				return true
			}}
		con *websocket.Conn
	)

	//upgrade websocket
	con, err = upgrader.Upgrade(writer, request, nil)
	defer func() {
		if err := con.Close(); err != nil {
			logs.Error(err)
		}
	}()
	if err != nil {
		if _, err := writer.Write([]byte(err.Error())); err != nil {
			logs.Error(err)
		}
	}

	//websocket conn
	for {
		msgType, data, err = con.ReadMessage()
		if err != nil {
			if _, err := writer.Write([]byte(err.Error())); err != nil {
				logs.Error(err)
			}
			break
		}
		fmt.Println("client msg:", msgType, string(data))
		if err := con.WriteMessage(websocket.TextMessage, []byte(time.Now().String())); err != nil {
			logs.Error(err)
		}
	}
}
