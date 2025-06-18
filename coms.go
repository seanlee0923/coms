package coms

import (
	"github.com/gorilla/websocket"
	"net/http"
)

type WebSocketInstance interface {
	getHandler(string) Handler
}

func DefaultUpgrade() (u *websocket.Upgrader) {
	return &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
}
