package server

import (
	"encoding/json"
)

type MessageType string

const (
	Awoo MessageType = "awoo"
)

type Message struct {
	Type    MessageType `json:"messageType"`
	Payload interface{} `json:"payload"`
}

func NewMessage(t MessageType, p interface{}) ([]byte, error) {
	m := Message{
		Type:    t,
		Payload: p,
	}

	b, err := json.Marshal(m)
	if err != nil {
		return []byte{}, err
	}
	return b, nil
}
