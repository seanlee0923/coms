package protocol

import "encoding/json"

type Message struct {
	Action string `json:"action"`
	Data   Data   `json:"data"`
}

type Data any

func ToMessage(message []byte) (*Message, error) {
	var rawMsg []json.RawMessage
	var msg Message
	err := json.Unmarshal(message, &rawMsg)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(rawMsg[0], &msg)

	return nil, nil
}

func (m *Message) ToBytes() ([]byte, error) {
	return json.Marshal(m)
}
