package coms

import (
	"github.com/gorilla/websocket"
	"github.com/seanlee0923/coms/protocol"
	"sync"
)

type OperationServer struct {
	addr    address
	timout  protocol.TimeOutConfig
	clients map[string]*Client
	mu      sync.Mutex
	u       websocket.Upgrader
}

type address struct {
	addr string
	port string
	path string
}
