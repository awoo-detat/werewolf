package client

import (
	"encoding/json"
	"log/slog"
)

type MessageType string

const (
	Awoo MessageType = "awoo"
)

type Message struct {
	Type MessageType `json:"messageType"`
}

func Decode(raw []byte) *Message {
	var message Message
	err := json.Unmarshal(raw, &message)
	if err != nil {
		slog.Error("error decoding client message", "error", err)
	}
	return &message
}
