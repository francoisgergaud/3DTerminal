package webserver

import (
	websocket "francoisgergaud/3dGame/common/connector"
	"francoisgergaud/3dGame/common/event"
	"francoisgergaud/3dGame/common/runner"
	"francoisgergaud/3dGame/server"
	websocketconnector "francoisgergaud/3dGame/server/connector/websocket"
	"log"
	"net/http"
)

//NewWebServer is a factory for web-server
func NewWebServer(server server.Server, address string, upgrader websocketconnector.WebsocketUpgrader) *WebServer {
	return &WebServer{
		playerJoinHandler: &PlayerJoinHandler{
			runner:                           &runner.AsyncRunner{},
			server:                           server,
			upgrader:                         upgrader,
			websocketClientConnectionFactory: websocketconnector.NewWebSocketClientConnection,
			websocketClientListenerFactory:   websocketconnector.NewClientWebSocketListener,
			websocketClientSenderFactory:     websocketconnector.NewClientWebSocketSender,
		},
		httpServer:    &HttpServerWrapper{},
		serverAddress: address,
	}
}

//WebServer implements a web-server
type WebServer struct {
	serverAddress     string
	httpServer        HttpServer
	playerJoinHandler *PlayerJoinHandler
}

//Run starts a blocking loop to listen  for new websocket connections
//TODO: handle shutdown gracefully (ob both server and client side)
func (webServer *WebServer) Run() {
	webServer.httpServer.Handle("/join", webServer.playerJoinHandler)
	err := webServer.httpServer.ListenAndServe(webServer.serverAddress, nil)
	if err != nil {
		log.Print("Error from server: ", err)
	}
}

//PlayerJoinHandler defines the handler for new-player join http event
type PlayerJoinHandler struct {
	runner                           runner.Runner
	server                           server.Server
	upgrader                         websocketconnector.WebsocketUpgrader
	websocketClientConnectionFactory func(eventToSendToCLient chan event.Event) *websocketconnector.WebSocketClientConnection
	websocketClientListenerFactory   func(playerID string, wsConnection websocket.WebsocketConnection, server server.Server) *websocketconnector.ClientWebSocketListener
	websocketClientSenderFactory     func(wsConnection websocket.WebsocketConnection, eventToSendToCLient chan event.Event) *websocketconnector.ClientWebSocketSender
}

func (joinHandler *PlayerJoinHandler) ServeHTTP(writer http.ResponseWriter, reader *http.Request) {
	connection, err := joinHandler.upgrader.Upgrade(writer, reader, nil)
	if err != nil {
		log.Println(err)
		return
	}
	eventsToSendToClient := make(chan event.Event)
	webSocketClientConnection := joinHandler.websocketClientConnectionFactory(eventsToSendToClient)
	//TODO: beware of the order: ClientSender must be ready to unqueue events from server for the player before register-player,
	//as register-player would block when sending the initialization-event otherwise: make it non-blocking
	joinHandler.runner.Start(joinHandler.websocketClientSenderFactory(connection, eventsToSendToClient))
	playerID := joinHandler.server.RegisterPlayer(webSocketClientConnection)
	joinHandler.runner.Start(joinHandler.websocketClientListenerFactory(playerID, connection, joinHandler.server))

}
