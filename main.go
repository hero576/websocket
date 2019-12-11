package main

import "websocket/wsconn"

func main() {
	ws := wsconn.NewWebsocketServer("127.0.0.1:8989")
	ws.SeverListener()
}
