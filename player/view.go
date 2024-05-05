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
	Role      *role.Role
	Hit       bool
	GamePhase int
}

func NewAttributeView(p *Player, a role.Attribute, hit bool, phase int) *View {
	return &View{
		Player:    p,
		Attribute: a,
		Hit:       hit,
		GamePhase: phase,
	}
}

func NewRoleView(p *Player, r *role.Role, phase int) *View {
	return &View{
		Player:    p,
		Role:      r,
		Hit:       true, // by definition you get the real role
		GamePhase: phase,
	}
}

func (v *View) For() string {
	if v.Role != nil {
		return v.Role.Name
	}
	return v.Attribute.String()
}
