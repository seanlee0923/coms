package coms

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/seanlee0923/coms/protocol"
	"net/http"
	"strings"
	"sync"
)

var s *OperationServer

type OperationServer struct {
	addr    address
	timeOut protocol.TimeOutConfig
	clients map[string]*Client
	mu      sync.Mutex
	u       websocket.Upgrader

	handler        map[string]Handler
	maxPendingCall int
}

type address struct {
	addr string
	path string
}

type Handler func(*Client, *protocol.Message) *json.RawMessage

func NewServer(addr, path string, u *websocket.Upgrader) *OperationServer {
	if u == nil {
		u = DefaultUpgrade()
	}

	server := &OperationServer{
		addr:           address{addr: addr, path: path},
		clients:        make(map[string]*Client),
		u:              *u,
		maxPendingCall: 32,
		handler:        make(map[string]Handler),
	}

	s = server

	return server
}

func (s *OperationServer) RegisterHandler(action string, handler Handler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handler[action] = handler
}

func (s *OperationServer) getHandler(action string) Handler {
	s.mu.Lock()
	defer s.mu.Unlock()

	h, ok := s.handler[action]
	if ok {
		return h
	}
	return nil
}

func (s *OperationServer) Start(h func(http.ResponseWriter, *http.Request)) error {
	if h == nil {
		h = func(w http.ResponseWriter, r *http.Request) {
			conn, err := s.u.Upgrade(w, r, nil)
			if err != nil {
				return
			}
			path := strings.Split(r.URL.Path, "/")
			clientId := path[len(path)-1]
			client := s.makeClient(clientId, conn)
			go client.run()
			s.Add(client)
		}

	}
	http.HandleFunc(s.addr.path, h)

	return http.ListenAndServe(s.addr.addr, nil)

}

func (s *OperationServer) Remove(c *Client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.clients, c.id)
}

func (s *OperationServer) Add(c *Client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.clients[c.id] = c
}

func (s *OperationServer) GetClient(id string) (*Client, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	c, ok := s.clients[id]
	return c, ok
}
