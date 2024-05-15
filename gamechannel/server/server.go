package server

import (
	"encoding/json"
)

type MessageType string

const (
	Awoo            MessageType = "awoo"
	IDSet                       = "idSet"
	NameSet                     = "nameSet"
	PlayerJoin                  = "playerJoin"
	PlayerLeave                 = "playerLeave"
	AlivePlayerList             = "alivePlayerList"
	RolesetList                 = "rolesetList"
	RolesetSelected             = "rolesetSelected"
	LeaderSet                   = "leaderSet"
	Password                    = "password"
	TallyChanged                = "tallyChanged"
	RoleAssigned                = "roleAssigned"
	PhaseChanged                = "phaseChanged"
	View                        = "view"
	PlayerKilled                = "playerKilled"
	GameOver                    = "gameOver"
	Error                       = "error"
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
