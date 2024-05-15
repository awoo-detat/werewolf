package gamechannel

import (
	"github.com/google/uuid"
)

type ActivityType int

const (
	SetName ActivityType = iota
	SetRoleset
	Reconnect
	Vote
	Quit
	ResetGame
)

type Activity struct {
	Type  ActivityType
	From  uuid.UUID
	To    uuid.UUID
	Value interface{}
}
