package coms

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/seanlee0923/coms/logger"
	"github.com/seanlee0923/coms/protocol"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"time"
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

func (c *Client) heartBeat() {
	ticker := time.NewTicker(c.heartBeatPeriod)
	defer ticker.Stop()

	for range ticker.C {
		resp, err := c.Call("HeartBeat", protocol.HeartBeatReq{})
		if err != nil {
			continue
		}
		var hb protocol.HeartBeatResp
		err = json.Unmarshal(resp.Data, &hb)
		if err != nil {
			continue
		}

		logger.InfoF("Heart Beat Status:%v", hb.Ok)
	}
}

func (c *Client) start() {
	go c.readLoop(c)
	go c.writeLoop()
	go c.collectStatus()
}

func (c *Client) collectStatus() {
	ticker := time.NewTicker(c.collectPeriod)

	defer ticker.Stop()

	for range ticker.C {
		stats, err := collect()
		if err != nil {
			continue
		}

		resp, err := c.Call("Status", stats)
		if err != nil {
			continue
		}

		logger.InfoF("status resp: %v", resp)
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
