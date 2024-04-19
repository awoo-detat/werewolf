package player

import (
	"github.com/awoo-detat/werewolf/role"
)

// A View represents a player's knowledge of another's role.
// It can be both a nighttime viewing (such as a seer viewing
// for max evil) or a generic "this player now knows the roles
// of one or more other players" (such as the wolves at the start
// of the game). It may not be successful, as in the case of
// a seer viewing a non-max, and it may even be incorrect,
// if there are tinkers!
type View struct {
	Player    *Player
	Attribute role.Attribute
	Hit       bool
	GamePhase int
}
