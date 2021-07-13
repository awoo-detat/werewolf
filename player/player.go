package player

import (
	"github.com/awoo-detat/werewolf/role"
)

// A GamePlayer implements the Player interface. It's likely that consuming code will
// need to extend this class, to implement private messaging and so forth.
type GamePlayer struct {
	ID     string `json:"id"`
	Name   string `json:"name,omitempty"`
	Leader bool   `json:"leader"`
	role   *role.Role
}

// A Revealed is a player who is either alive or whose role has been revealed.
// TODO revisit this?
type Revealed struct {
	ID       string `json:"id"`
	Name     string `json:"name,omitempty"`
	RoleName string `json:"role"`
	Alive    bool   `json:"alive"`
}

// New created a new GamePlayer. The name should be set separately.
func New(id string) *GamePlayer {
	p := &GamePlayer{
		ID:   id,
		role: &role.Role{},
	}

	return p
}

// GetID returns the string representation of the player's unique identifier.
func (p *GamePlayer) GetID() string {
	return p.ID
}

// GetIdentifier gets the best human-readable way to identify a player, prioritizing
// their name if it exists and defaulting to their ID,
func (p *GamePlayer) GetIdentifier() string {
	if p.Name != "" {
		return p.Name
	}
	return p.ID
}

// Reveal creates a Revealed from a GamePlayer.
func (p *GamePlayer) Reveal() *Revealed {
	r := &Revealed{
		ID:   p.ID,
		Name: p.Name,
	}
	if p.Role().Alive {
		r.Alive = true
	} else {
		r.RoleName = p.Role().Name
	}
	return r
}

// Role returns the GamePlayer's role.
func (p *GamePlayer) Role() *role.Role {
	return p.role
}

// SetLeader marks the player as the game's leader, able to start it.
func (p *GamePlayer) SetLeader() {
	p.Leader = true
}

// SetName changes the player's name.
func (p *GamePlayer) SetName(name string) {
	p.Name = name
}

// SetRole assigns the player's role.
func (p *GamePlayer) SetRole(r *role.Role) {
	p.role = r
}
