package coms

import (
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	c := NewClient("test1", 32)
	c.SetPeriod(time.Minute, time.Minute)
	err := c.Start("ws://localhost:9999/my-test/")
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second)
	t.Log("start success")
	req := Mock1{
		Name: "test1-ack",
		Hi:   "ack",
	}

	t.Log(c.Call("Act1", req))
	time.Sleep(time.Second)

	t.Log(c.Call("Act1", req))
	time.Sleep(time.Second)
	t.Log(c.Call("Act1", req))
	time.Sleep(time.Second)
	t.Log(c.Call("Act1", req))
	time.Sleep(time.Second)
	t.Log(c.Call("Act1", req))
	time.Sleep(time.Second)
	t.Log(c.Call("Act1", req))
	time.Sleep(time.Second)
	select {}
}
