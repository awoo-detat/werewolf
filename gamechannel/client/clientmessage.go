package client

import (
	"encoding/json"
	"log/slog"

	"github.com/google/uuid"
)

type MessageType string

const (
	Awoo        MessageType = "awoo"
	SetName                 = "setName"
	SetRoleset              = "setRoleset"
	Vote                    = "vote"
	NightAction             = "nightAction"
	Quit                    = "quit"
)

type Message struct {
	Type       MessageType `json:"messageType"`
	PlayerName string      `json:"playerName"`
	Roleset    string      `json:"roleset"`
	Target     uuid.UUID   `json:"target"`
}

func Decode(raw []byte) (*Message, error) {
	var message Message
	err := json.Unmarshal(raw, &message)
	if err != nil {
		slog.Error("error decoding client message", "error", err, "raw", string(raw[:]))
		return nil, err
	}
	if message.Type == "" {
		slog.Error("message did not decode properly", "message", message, "raw", string(raw[:]))
	}
	return &message, nil
}
