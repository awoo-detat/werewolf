package server

import (
	"github.com/awoo-detat/werewolf/role"

	"github.com/google/uuid"
)

type RevealedPlayer struct {
	ID   uuid.UUID  `json:"id"`
	Name string     `json:"name"`
	Role *role.Role `json:"role"`
}
