package gamechannel

import (
	"github.com/google/uuid"
)

type ActivityType int

const (
	Join ActivityType = iota
	SetName
	SetRoleset
	Vote
	Quit
	ResetGame
)

type Activity struct {
	Type   ActivityType
	Player uuid.UUID
	Value  interface{}
}
