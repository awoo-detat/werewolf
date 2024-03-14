package game

import (
	"fmt"

	"github.com/awoo-detat/werewolf/role/roleset"
)

type StateError struct {
	NeedState GameState
	InState   GameState
}

func (e *StateError) Error() string {
	return fmt.Sprintf("game: in state %v instead of required %v", e.InState, e.NeedState)
}

type PlayerCountError struct {
	Roleset     *roleset.Roleset
	PlayerCount int
}

func (e *PlayerCountError) Error() string {
	return fmt.Sprintf("game: roleset %s requires %d participants, have %d", e.Roleset.Name, len(e.Roleset.Roles), e.PlayerCount)
}
