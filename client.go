package coms

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/seanlee0923/coms/protocol"
	"sync"
	"sync/atomic"
	"time"
)

type Client struct {
	id     string
	conn   *websocket.Conn
	timout protocol.TimeOutConfig

	handler map[string]Handler

	pingCh     chan []byte
	messageIn  chan []byte
	messageOut chan []byte
	closeCh    chan bool

	connected bool

	pendingCalls    sync.Map
	pendingCnt      atomic.Int32
	maxPendingCalls int
}

func (s *OperationServer) makeClient(id string, conn *websocket.Conn) *Client {

	cli := &Client{
		id:   id,
		conn: conn,

		pingCh:          make(chan []byte),
		messageIn:       make(chan []byte),
		messageOut:      make(chan []byte),
		closeCh:         make(chan bool, 1),
		handler:         make(map[string]Handler),
		maxPendingCalls: s.maxPendingCall,
	}

	cli.conn.SetPingHandler(func(appData string) error {
		cli.pingCh <- []byte(appData)
		return cli.conn.SetWriteDeadline(time.Now().Add(cli.timout.PingWait))
	})

	return cli
}

func (c *Client) getId() string {
	return c.id
}

func (c *Client) run() {
	go c.readLoop(s)
	go c.writeLoop()
}

func (c *Client) readLoop(w WebSocketInstance) {

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
		if message.Type == protocol.Resp {
			if call, ok := c.pendingCalls.Load(message.Id); ok {
				if callCh, ok := call.(chan *protocol.Message); ok {
					callCh <- message
				}
			}
			continue
		}

		h := w.getHandler(message.Action)
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
			Id:     uuid.NewString(),
			Type:   protocol.Resp,
			Action: message.Action,
			Data:   *respData,
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

		select {

		case msg, ok := <-c.messageOut:
			if !ok {
				c.closeCh <- true
				return
			}

			writer, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				c.closeCh <- true
				return
			}

			_, err = writer.Write(msg)
			if err != nil {
				c.closeCh <- true
				return
			}

		case <-c.pingCh:

			err := c.conn.WriteMessage(websocket.PongMessage, []byte{})
			if err != nil {
				c.closeCh <- true
				return
			}

		case <-c.closeCh:

			cm := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")
			err := c.conn.WriteMessage(websocket.CloseMessage, cm)
			if err != nil {
				_ = c.conn.NetConn().Close()
				break
			}

		}

	}
}

func (c *Client) Call(data any, action string) (*protocol.Message, error) {

	if c.pendingCnt.Load() >= int32(c.maxPendingCalls) {
		return nil, errors.New("max pending calls exceeded")
	}

	raw, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req := &protocol.Message{
		Id:     uuid.NewString(),
		Type:   protocol.Req,
		Action: action,
		Data:   raw,
	}

	respCh := make(chan *protocol.Message, 1)
	c.pendingCalls.Store(req.Id, respCh)
	c.pendingCnt.Add(1)
	defer func() {
		c.pendingCalls.Delete(req.Id)
		c.pendingCnt.Add(-1)
	}()

	msgBytes, err := req.ToBytes()
	if err != nil {
		return nil, err
	}

	c.messageOut <- msgBytes
	select {

	case resp := <-respCh:
		return resp, nil

	case <-time.After(c.timout.ReadWait):
		return nil, errors.New("timeout")

	}

}

func (c *Client) getHandler(action string) Handler {
	return c.handler[action]
}
