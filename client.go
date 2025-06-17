package coms

import (
	"github.com/gorilla/websocket"
	"github.com/seanlee0923/coms/protocol"
	"runtime"
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

		pingCh:     make(chan []byte),
		messageIn:  make(chan []byte),
		messageOut: make(chan []byte),
		closeCh:    make(chan bool),
	}
}

func (c *Client) GetId() string {
	return c.id
}

func RunClient(c *Client) {
	go func() {
		for {
			if !c.readLoop() {
				break
			}
		}
	}()

	go func() {
		if !c.writeLoop() {
			runtime.Breakpoint()
		}
	}()
}

func (c *Client) readLoop() bool {
	_, msg, err := c.conn.ReadMessage()
	if err != nil {
		c.closeCh <- true
		return false
	}

	message, err := protocol.ToMessage(msg)
	if err != nil {
		c.closeCh <- true
		return false
	}

	if message == nil {
		c.closeCh <- true
		return false
	}

	h := s.GetHandler(message.Action)
	if h == nil {
		c.closeCh <- true
		return false
	}

	resp := h(c, message)
	if resp == nil {
		c.closeCh <- true
		return false
	}

	msgOut, err := resp.ToBytes()
	if err != nil {
		c.closeCh <- true
		return false
	}

	c.messageOut <- msgOut

	return true
}

func (c *Client) writeLoop() bool {
	return true
}
