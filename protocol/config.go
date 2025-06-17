package protocol

import "time"

type TimeOutConfig struct {
	PingWait   time.Duration
	WriteWait  time.Duration
	ReadWait   time.Duration
	PingPeriod time.Duration
}
