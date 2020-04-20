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

//NewWebServer is a factory for web-server
func NewWebServer(server server.Server, address string) *WebServer {
	return &WebServer{
		server:        server,
		serverAddress: address,
	}
}

//WebServer implement a web-server
type WebServer struct {
	server        server.Server
	serverAddress string
}

//Start starts to listen  for new websocket connections
func (webServer *WebServer) Start() {
	http.HandleFunc("/join", func(w http.ResponseWriter, r *http.Request) {
		connection, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		websocketconnector.RegisterWebSocketClientConnectionToServer(connection, webServer.server)
	})
	err := http.ListenAndServe(webServer.serverAddress, nil)
	if err != nil {
		log.Fatal("Error from server: ", err)
	}
}
