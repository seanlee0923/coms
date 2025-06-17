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

		pingCh:     make(chan []byte),
		messageIn:  make(chan []byte),
		messageOut: make(chan []byte),
		closeCh:    make(chan bool),
	}
}

func (c *Client) GetId() string {
	return c.id
}

func (c *Client) Run() {
	go c.readLoop()
	go c.writeLoop()
}

func (c *Client) readLoop() {

	defer s.Remove(c)

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			c.closeCh <- true
			return
		}

		message, err := protocol.ToMessage(msg)
		if err != nil {
			c.closeCh <- true
			return
		}

		if message == nil {
			c.closeCh <- true
			return
		}

		h := s.GetHandler(message.Action)
		if h == nil {
			c.closeCh <- true
			return
		}

		respData := h(c, message)
		if respData == nil {
			c.closeCh <- true
			return
		}

		resp := protocol.Message{
			Action: message.Action,
			Data:   respData,
		}

		msgOut, err := resp.ToBytes()
		if err != nil {
			c.closeCh <- true
			return
		}

		c.messageOut <- msgOut

	}

}

func (c *Client) writeLoop() {

	defer s.Remove(c)

	for {
		return
	}
}
