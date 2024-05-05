package player

import (
	"log/slog"

	"github.com/awoo-detat/werewolf/role"

	"github.com/google/uuid"
)

type Player struct {
	ID    uuid.UUID
	Name  string
	Role  *role.Role
	Views []*View
}

func NewPlayer() *Player {
	p := &Player{
		ID:    uuid.New(),
		Views: []*View{},
	}
	return p
}

func (p *Player) String() string {
	if len(p.Name) != 0 {
		return p.Name
	}
	return p.ID.String()
}

func (p *Player) SetName(name string) {
	p.Name = name
	slog.Info("setting player name", "ID", p.ID, "Name", p.Name)
}

func (p *Player) SetRole(r *role.Role) {
	p.Role = r
	slog.Info("setting player role", "player", p, "role", r)
}
