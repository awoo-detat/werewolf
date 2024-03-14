package game

import (
	"fmt"

	"github.com/awoo-detat/werewolf/player"
	"github.com/awoo-detat/werewolf/role/roleset"

	"github.com/google/uuid"
)

type Game struct {
	ID      uuid.UUID
	Leader  *player.Player
	Players map[uuid.UUID]*player.Player
	Roleset *roleset.Roleset
	state   GameState
	Phase   int
}

type GameState int

const (
	Setup GameState = iota
	Running
	Finished
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

func (g *Game) State() GameState {
	return g.state
}

func (g *Game) AddPlayer(p *player.Player) {
	g.Players[p.ID] = p
}

func (g *Game) ChooseRoleset(slug string) error {
	if g.state > Setup {
		return &StateError{NeedState: Setup, InState: g.state}
	}
	rs, ok := roleset.List()[slug]
	if !ok {
		return fmt.Errorf("roleset %s not found", slug)
	}

	g.Roleset = rs
	return nil
}

func (g *Game) StartGame() error {
	if g.state > Setup {
		return &StateError{NeedState: Setup, InState: g.state}
	}
	if g.Roleset == nil {
		return fmt.Errorf("game: no roleset defined") // TODO
	}
	if len(g.Players) != len(g.Roleset.Roles) {
		return &PlayerCountError{Roleset: g.Roleset, PlayerCount: len(g.Players)}
	}

	g.state = Running
	return nil
}
