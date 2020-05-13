package webserver

import (
	"francoisgergaud/3dGame/server"
	websocketconnector "francoisgergaud/3dGame/server/connector/websocket"
	"log"
	"net/http"
)

//NewWebServer is a factory for web-server
func NewWebServer(server server.Server, address string, upgrader websocketconnector.WebsocketUpgrader) *WebServer {
	return &WebServer{
		server:        server,
		serverAddress: address,
		upgrader:      upgrader,
	}
}

//WebServer implement a web-server
type WebServer struct {
	server        server.Server
	serverAddress string
	upgrader      websocketconnector.WebsocketUpgrader
}

//Start starts to listen  for new websocket connections
func (webServer *WebServer) Start() {
	http.HandleFunc("/join", func(w http.ResponseWriter, r *http.Request) {
		connection, err := webServer.upgrader.Upgrade(w, r, nil)
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
