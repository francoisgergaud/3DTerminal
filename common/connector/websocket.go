package websocket

import (
	"github.com/gorilla/websocket"
)

//WebsocketConnection represent a websocket-connection
type WebsocketConnection interface {
	ReadJSON(v interface{}) error
	WriteJSON(v interface{}) error
	Close() error
}

type WebsocketConnectionWrapper struct {
	internalConnection *websocket.Conn
}

func (conn *WebsocketConnectionWrapper) ReadJSON(v interface{}) error {
	return conn.internalConnection.ReadJSON(v)
}

func (conn *WebsocketConnectionWrapper) WriteJSON(v interface{}) error {
	return conn.internalConnection.WriteJSON(v)
}

func (conn *WebsocketConnectionWrapper) Close() error {
	return conn.internalConnection.Close()
}
