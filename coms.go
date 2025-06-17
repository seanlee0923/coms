package coms

import (
	"github.com/gorilla/websocket"
	"net/http"
)

func DefaultUpgrade() (u *websocket.Upgrader) {
	return &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
}
