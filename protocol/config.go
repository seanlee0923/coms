package protocol

import "time"

const (
	MinCollectDuration   = 30 * time.Second
	MinHeartBeatDuration = 30 * time.Second
)

type TimeOutConfig struct {
	PingPeriod time.Duration
	PongWait   time.Duration
	WriteWait  time.Duration
	ReadWait   time.Duration
}
