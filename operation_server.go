package coms

import (
	"github.com/gorilla/websocket"
	"github.com/seanlee0923/coms/protocol"
	"net/http"
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
	path string
}

func NewServer(addr, path string, u *websocket.Upgrader) *OperationServer {
	if u == nil {
		u = DefaultUpgrade()
	}
	return &OperationServer{
		addr:    address{addr: addr, path: path},
		clients: make(map[string]*Client),
		u:       *u,
	}
}

func (s *OperationServer) Start(h func(http.ResponseWriter, *http.Request)) error {

	http.HandleFunc(s.addr.path, h)

	return http.ListenAndServe(s.addr.addr, nil)

}
