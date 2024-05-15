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
	NightAction
	Quit
	ResetGame
)

type Activity struct {
	Type  ActivityType
	From  uuid.UUID
	Value interface{}
}
