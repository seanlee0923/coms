package coms

import (
	"github.com/gorilla/websocket"
	"github.com/seanlee0923/coms/protocol"
	"sync"
)

type Client struct {
	id     string
	conn   *websocket.Conn
	mu     sync.Mutex
	timout protocol.TimeOutConfig
}

func NewClient(id string, conn *websocket.Conn) *Client {
	return &Client{
		id:   id,
		conn: conn,
	}
}
