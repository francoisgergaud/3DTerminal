package websocketconnector

import (
	"net/http"

	gorillaWebsocket "github.com/gorilla/websocket"

	websocket "francoisgergaud/3dGame/common/connector"
)

type WebsocketUpgrader interface {
	Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (websocket.WebsocketConnection, error)
}

func NewWebsocketUpgraderWwrapper(upgrader *gorillaWebsocket.Upgrader) WebsocketUpgrader {
	return &WebsocketUpgraderWwrapper{
		internalUpgrader: upgrader,
	}
}

type WebsocketUpgraderWwrapper struct {
	internalUpgrader *gorillaWebsocket.Upgrader
}

func (upgraderWwrapper WebsocketUpgraderWwrapper) Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (websocket.WebsocketConnection, error) {
	return upgraderWwrapper.internalUpgrader.Upgrade(w, r, responseHeader)
}
