package coms

import (
	"encoding/json"
	"github.com/seanlee0923/coms/logger"
	"github.com/seanlee0923/coms/protocol"
	"testing"
)

type Mock1 struct {
	Name string `json:"name"`
	Hi   string `json:"hi"`
}

type Mock2 struct {
	Age int    `json:"age"`
	Bye string `json:"bye"`
}

func TestServer(t *testing.T) {
	tServer := NewServer("0.0.0.0:9999", "/my-test/", nil)
	t.Log("creat success")
	tServer.RegisterHandler("Act1", func(client *Client, message *protocol.Message) *json.RawMessage {
		//var data Mock1
		var rm json.RawMessage

		resp := Mock1{
			Name: "test1-ack",
			Hi:   "ack",
		}
		msg, err := json.Marshal(&resp)
		if err != nil {
			logger.Error(err)
			return nil
		}

		err = json.Unmarshal(msg, &rm)
		if err != nil {
			logger.Error(err)
		}
		logger.Info("message send successfully")
		return &rm
	})

	tServer.RegisterHandler("Login", func(client *Client, message *protocol.Message) *json.RawMessage {
		var rm json.RawMessage

		return &rm
	})

	tServer.RegisterHandler("HeartBeat", func(client *Client, message *protocol.Message) *json.RawMessage {
		var rm json.RawMessage

		resp := protocol.HeartBeatResp{
			Ok: true,
		}

		msg, err := json.Marshal(&resp)
		if err != nil {
			logger.Error(err)
			return nil
		}

		err = json.Unmarshal(msg, &rm)
		if err != nil {
			logger.Error(err)
		}

		return &rm
	})

	t.Log("register success")
	if err := tServer.Start(nil); err != nil {
		t.Error(err)
	}

}
