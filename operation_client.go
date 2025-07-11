package coms

import "C"
import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/seanlee0923/coms/logger"
	"github.com/seanlee0923/coms/protocol"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"sync"
	"sync/atomic"
	"time"
)

type OperationClient struct {
	id      string
	conn    *websocket.Conn
	timeout protocol.TimeOutConfig

	handler map[string]ClientHandler

	pingCh     chan []byte
	messageIn  chan []byte
	messageOut chan []byte
	closeCh    chan bool

	connected bool

	pendingCalls    sync.Map
	pendingCnt      atomic.Int32
	maxPendingCalls int

	finCh           chan bool
	heartBeatPeriod time.Duration
	collectPeriod   time.Duration
}

type ClientHandler func(*OperationClient, *protocol.Message) *json.RawMessage

func NewClient(id string, maxPendingCalls int) *OperationClient {
	cli := &OperationClient{
		id: id,

		pingCh:          make(chan []byte),
		messageOut:      make(chan []byte),
		closeCh:         make(chan bool),
		finCh:           make(chan bool),
		maxPendingCalls: maxPendingCalls,
	}

	return cli
}

func (c *OperationClient) SetTimeout(tc protocol.TimeOutConfig) {
	c.timeout = tc
}

func (c *OperationClient) SetPeriod(collect, heartBeat time.Duration) {
	c.collectPeriod = collect
	c.heartBeatPeriod = heartBeat
}

func (c *OperationClient) Start(addr string, finCh chan bool) error {
	if c.collectPeriod <= protocol.MinCollectDuration {
		return errors.New("collect period too small")
	}
	if c.heartBeatPeriod <= protocol.MinHeartBeatDuration {
		return errors.New("heart beat period too small")
	}
	if c.timeout == (protocol.TimeOutConfig{}) {
		return errors.New("need to set timeout config")
	}
	//if len(finCh) != 2 {
	//	return errors.New("invalid finCh length")
	//}

	conn, _, err := websocket.DefaultDialer.Dial(addr+c.id, nil)
	if err != nil {
		return err
	}
	c.conn = conn

	go c.start(finCh)

	return nil
}

func (c *OperationClient) RegisterHandler(action string, handler ClientHandler) {
	c.handler[action] = handler
}

func (c *OperationClient) Send(action string, data any) (*protocol.Message, error) {

	raw, err := json.Marshal(data)
	if err != nil {
		logger.Error(err)
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
		logger.Error(err)
		return nil, err
	}

	c.messageOut <- msgBytes

	select {

	case resp := <-respCh:
		return resp, nil

	case <-time.After(c.timeout.ReadWait):
		return nil, errors.New("timeout")

	}

}

func (c *OperationClient) heartBeat() {
	ticker := time.NewTicker(c.heartBeatPeriod)
	defer ticker.Stop()

	for range ticker.C {

		select {

		case <-ticker.C:
			resp, err := c.Send("HeartBeat", protocol.HeartBeatReq{})
			if err != nil {
				continue
			}

			var hb protocol.HeartBeatResp
			err = json.Unmarshal(resp.Data, &hb)
			if err != nil {
				continue
			}

		case <-c.finCh:

			return

		}

	}
}

func (c *OperationClient) start(finCh chan bool) {
	go c.readLoopClient()
	go c.writeLoopClient()
	go c.collectStatus()
	go c.heartBeat()
}

func (c *OperationClient) readLoopClient() {

	defer logger.Info("read break")

	c.conn.SetPongHandler(func(appData string) error {
		return c.conn.SetReadDeadline(time.Now().Add(c.timeout.PongWait))
	})

	for {
		_, msg, err := c.conn.ReadMessage()
		logger.InfoF("%s | %s", c.id, string(msg))

		if err != nil {
			logger.Errorf("send true to close ch, %v", err)
			c.closeCh <- true
			return
		}

		message, err := protocol.ToMessage(msg)
		if err != nil {
			logger.Error(err)
			c.closeCh <- true
			break
		}

		if message == nil {
			c.closeCh <- true
			break
		}

		if message.Type == protocol.Resp {

			if call, ok := c.pendingCalls.Load(message.Id); ok {

				if callCh, ok := call.(chan *protocol.Message); ok {

					callCh <- message
				}
			}
			continue
		}

		h := c.getHandler(message.Action)
		if h == nil {
			logger.Info("not found handler")
			continue
		}

		respData := h(c, message)
		if respData == nil {
			logger.Info("nil resp data")
			c.closeCh <- true
			break
		}

		resp := protocol.Message{
			Id:     message.Id,
			Type:   protocol.Resp,
			Action: message.Action,
			Data:   *respData,
		}

		msgOut, err := resp.ToBytes()
		if err != nil {
			logger.Error(err)
			c.closeCh <- true
			return
		}

		c.messageOut <- msgOut

	}

}

func (c *OperationClient) writeLoopClient() {

	defer logger.Info("write break")
	for {

		select {

		case msg, ok := <-c.messageOut:
			if !ok {
				c.closeCh <- true
				return
			}

			writer, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				logger.Error(err)
				c.closeCh <- true
				return
			}

			mLen, err := writer.Write(msg)
			if err != nil {
				logger.Error(err)
				c.closeCh <- true
				return
			}
			err = writer.Close()
			if err != nil {
				logger.Error(err)
				c.closeCh <- true
				return
			}
			logger.InfoF("msg send to server, %d", mLen)

		case <-c.pingCh:

			err := c.conn.WriteMessage(websocket.PongMessage, []byte{})
			if err != nil {
				logger.Error(err)
				c.closeCh <- true
				return
			}

		case <-c.closeCh:

			cm := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")
			err := c.conn.WriteMessage(websocket.CloseMessage, cm)
			if err != nil {
				logger.Error(err)
				_ = c.conn.Close()
				return

			}

			close(c.finCh)

		}

	}
}

func (c *OperationClient) collectStatus() {
	ticker := time.NewTicker(c.collectPeriod)

	defer ticker.Stop()
	for {

		select {
		case <-ticker.C:
			stats, err := collect()
			if err != nil {
				continue
			}

			_, err = c.Send("Status", stats)
			if err != nil {
				continue
			}

		case <-c.finCh:
			return
		}
	}

}

func collect() (*protocol.StatusReq, error) {
	// CPU 사용량
	cpuPercents, err := cpu.Percent(0, false)
	if err != nil {
		return nil, err
	}

	cpuCores, err := cpu.Counts(true)
	if err != nil {
		return nil, err
	}

	// 메모리
	vm, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	// 디스크
	du, err := disk.Usage("/")
	if err != nil {
		return nil, err
	}

	status := &protocol.StatusReq{
		CpuUsage:    cpuPercents[0],
		CpuCores:    cpuCores,
		MemTotal:    float64(vm.Total) / (1024 * 1024), // MB
		MemUsage:    float64(vm.Used) / (1024 * 1024),  // MB
		MemPercent:  vm.UsedPercent,
		DiskTotal:   float64(du.Total) / (1024 * 1024 * 1024), // GB
		DiskUsage:   float64(du.Used) / (1024 * 1024 * 1024),  // GB
		DiskPercent: du.UsedPercent,
		ServerTime:  time.Now().Format(time.RFC3339),
	}

	return status, nil
}

func (c *OperationClient) getHandler(action string) ClientHandler {

	h, ok := c.handler[action]
	if ok {
		return h
	}

	return nil
}
