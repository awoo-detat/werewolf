package server

import (
	"github.com/awoo-detat/werewolf/role"
)

type GameOverMessage struct {
	Winner role.PlayerType   `json:"winner"`
	Roles  []*RevealedPlayer `json:"roles"`
}
