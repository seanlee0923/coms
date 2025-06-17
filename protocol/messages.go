package protocol

type Message struct {
	Action string
	Data   Data
}

type Data any

func ToMessage(message []byte) (*Message, error) {
	return nil, nil
}
