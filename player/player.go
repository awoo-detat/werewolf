package player

import (
	"github.com/awoo-detat/werewolf/role"

	"github.com/google/uuid"
)

type Player struct {
	ID   uuid.UUID
	Name string
	Role *role.Role
}

func NewPlayer() *Player {
	p := &Player{
		ID: uuid.New(),
	}
	return p
}

func (p *Player) SetName(name string) {
	p.Name = name
}

func (p *Player) SetRole(r *role.Role) {
	p.Role = r
}
