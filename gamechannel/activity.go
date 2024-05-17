package gamechannel

import (
	"github.com/google/uuid"
)

type ActivityType int

const (
	SetName ActivityType = iota
	SetRoleset
	Reconnect
	Start
	Vote
	NightAction
	Quit
	ResetGame
	Awoo
)

type Activity struct {
	Type  ActivityType
	From  uuid.UUID
	Value interface{}
}
