package webserver

import (
	"francoisgergaud/3dGame/server"
	websocketconnector "francoisgergaud/3dGame/server/connector/websocket"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

//WebServer implement a web-server
type WebServer struct {
	server server.Server
}

func (webServer *WebServer) launch(serverAddress *string) {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		connection, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		websocketconnector.RegisterWebSocketClientConnectionToServer(connection, webServer.server)
	})
	err := http.ListenAndServe(*serverAddress, nil)
	if err != nil {
		log.Fatal("Error from server: ", err)
	}
}
