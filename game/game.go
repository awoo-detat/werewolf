package game

import (
	"github.com/awoo-detat/werewolf/player"
	"github.com/awoo-detat/werewolf/role/roleset"

	"github.com/google/uuid"
)

type Game struct {
	ID      uuid.UUID
	Leader  *player.Player
	Players map[uuid.UUID]*player.Player
	Roleset *roleset.Roleset
	State   GameState
	Phase   int
}

type GameState int

const (
	Setup    GameState = iota
	Running            = iota
	Finished           = iota
)

func NewGame(p *player.Player) *Game {
	g := &Game{
		ID:      uuid.New(),
		Leader:  p,
		Players: make(map[uuid.UUID]*player.Player),
	}
	g.AddPlayer(p)

	return g
}

func (g *Game) AddPlayer(p *player.Player) {
	g.Players[p.ID] = p
}

// start the game etc
