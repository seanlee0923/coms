package coms

import (
	"github.com/seanlee0923/coms/protocol"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	c := NewClient("test1", 32)
	c.SetPeriod(time.Minute, time.Minute)
	c.SetTimeout(protocol.TimeOutConfig{
		PingPeriod: 30 * time.Second,
		ReadWait:   30 * time.Second,
		WriteWait:  15 * time.Second,
	})
	err := c.Start("ws://localhost:9999/my-test/")
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(5 * time.Second)
	t.Log("start success")
	req := Mock1{
		Name: "test1-ack",
		Hi:   "ack",
	}

	_, err = c.Send("Act1", req)
	if err != nil {
		t.Log(err)
	}
	//time.Sleep(time.Second)
	//
	//t.Log(c.Call("Act1", req))
	//time.Sleep(time.Second)
	//t.Log(c.Call("Act1", req))
	//time.Sleep(time.Second)
	//t.Log(c.Call("Act1", req))
	//time.Sleep(time.Second)
	//t.Log(c.Call("Act1", req))
	//time.Sleep(time.Second)
	//t.Log(c.Call("Act1", req))
	//time.Sleep(time.Second)
	select {}
}
