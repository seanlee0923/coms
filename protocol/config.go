package protocol

import "time"

const (
	MinCollectDuration   = 30 * time.Second
	MinHeartBeatDuration = 30 * time.Second
)

type TimeOutConfig struct {
	PingWait   time.Duration
	WriteWait  time.Duration
	ReadWait   time.Duration
	PingPeriod time.Duration
}
