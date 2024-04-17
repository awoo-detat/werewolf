package game

import (
	"fmt"

	"github.com/awoo-detat/werewolf/player"
	"github.com/awoo-detat/werewolf/role/roleset"

	"github.com/google/uuid"
)

type Game struct {
	ID           uuid.UUID
	Leader       *player.Player
	Players      map[uuid.UUID]*player.Player
	Roleset      *roleset.Roleset
	state        GameState
	VotingMethod VotingMethod
	Phase        int
}

type GameState int

const (
	Setup GameState = iota
	Running
	Finished
)

type VotingMethod int

const (
	// InstaKill means that as soon as one player has 51% of the vote they are killed and the day ends.
	InstaKill VotingMethod = iota
	// InstaKillWithDelay is the same as InstaKill, but with a delay to allow for claiming and changing of votes. It is not currently supported.
	InstaKillWithDelay
	// Timed will end the day after a set number of minutes, killing whoever is in the lead (breaking ties with Longest Held Last Vote). It is not currently supported.
	Timed
)

func NewGame(p *player.Player) *Game {
	g := &Game{
		ID:           uuid.New(),
		Leader:       p,
		Players:      make(map[uuid.UUID]*player.Player),
		VotingMethod: InstaKill,
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

	// TODO N0 actions

	g.state = Running
	return nil
}

func (g *Game) Vote(from *player.Player, to *player.Player) error {
	if g.state != Running {
		return &StateError{NeedState: Running, InState: g.state}
	}
	if g.Phase%2 != 0 {
		return &PhaseError{GamePhase: g.Phase}
	}
	// TODO
	return nil
}
