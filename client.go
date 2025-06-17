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

	pingCh     chan []byte
	messageIn  chan []byte
	messageOut chan []byte
	closeCh    chan bool

	connected bool
}

func NewClient(id string, conn *websocket.Conn) *Client {
	return &Client{
		id:   id,
		conn: conn,
	}
}
