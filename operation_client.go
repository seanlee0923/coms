package coms

import (
	"github.com/gorilla/websocket"
	"github.com/seanlee0923/coms/protocol"
)

func NewClient(id string, maxPendingCalls int) *Client {
	cli := &Client{
		id: id,

		pingCh:          make(chan []byte),
		messageOut:      make(chan []byte),
		closeCh:         make(chan bool),
		maxPendingCalls: maxPendingCalls,
	}

	return cli
}

func (c *Client) SetTimeout(tc protocol.TimeOutConfig) {
	c.timout = tc
}

func (c *Client) Start(addr string) error {
	conn, _, err := websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		return err
	}
	c.conn = conn

	go c.start()
	return nil
}

func (c *Client) start() {
	go c.readLoop(c)
	go c.writeLoop()
}
