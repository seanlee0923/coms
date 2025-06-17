package coms

import (
	"github.com/gorilla/websocket"
	"net/http"
)

func DefaultUpgrade() (u *websocket.Upgrader) {
	return &websocket.Upgrader{
		ReadBufferSize:  8096,
		WriteBufferSize: 8096,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
}
