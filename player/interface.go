package player

import (
	"github.com/awoo-detat/werewolf/role"
)

// A Player is an abstract representation of a person playing in a Game.
type Player interface {
	GetID() string
	GetIdentifier() string
	Reveal() *Revealed
	//Vote(to string)
	//NightAction(to string)
	//Play()
	//InGame() bool
	//LeaveGame()
	SetLeader()
	SetName(name string)
	SetRole(r *role.Role)
	Role() *role.Role
	//Quit()
}
