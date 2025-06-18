package protocol

import "encoding/json"

type Message struct {
	Action string          `json:"action"`
	Data   json.RawMessage `json:"data"`
}

func ToMessage(message []byte) (*Message, error) {
	var msg Message

	err := json.Unmarshal(message, &msg)
	if err != nil {
		return nil, err
	}

	return &msg, nil
}

func (m *Message) ToBytes() ([]byte, error) {
	return json.Marshal(m)
}
