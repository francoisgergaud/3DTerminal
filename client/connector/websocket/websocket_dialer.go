package websocketconnector

import (
	"net/http"

	gorillaWebsocket "github.com/gorilla/websocket"

	websocket "francoisgergaud/3dGame/common/connector"
)

//WebsocketDialer is the client dialer
type WebsocketDialer interface {
	Dial(urlStr string, requestHeader http.Header) (websocket.WebsocketConnection, *http.Response, error)
}

//NewWebsocketDialerWrapper WebsocketDialeris a factory for WebsocketDialer
func NewWebsocketDialerWrapper() WebsocketDialer {
	return &WebsocketDialerWrapper{
		internalDialer: gorillaWebsocket.DefaultDialer,
	}
}

//WebsocketDialerWrapper implements WebsocketDialer by wrapping a Gorilla websocket
type WebsocketDialerWrapper struct {
	internalDialer *gorillaWebsocket.Dialer
}

//Dial connect a client websocket to a server
func (dialerWrapper WebsocketDialerWrapper) Dial(urlStr string, requestHeader http.Header) (websocket.WebsocketConnection, *http.Response, error) {
	return dialerWrapper.internalDialer.Dial(urlStr, requestHeader)
}
