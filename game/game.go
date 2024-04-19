package game

import (
	"fmt"
	"math/rand"

	"github.com/awoo-detat/werewolf/player"
	"github.com/awoo-detat/werewolf/role"
	"github.com/awoo-detat/werewolf/role/roleset"

	"github.com/google/uuid"
)

type Game struct {
	ID           uuid.UUID
	Leader       *player.Player
	AlivePlayers map[uuid.UUID]*player.Player
	Players      map[uuid.UUID]*player.Player
	playerSlice  []*player.Player
	Roleset      *roleset.Roleset
	state        GameState
	Phase        int
}

type GameState int

const (
	Setup GameState = iota
	Running
	Finished
)

func NewGame(p *player.Player) *Game {
	g := &Game{
		ID:           uuid.New(),
		Leader:       p,
		AlivePlayers: make(map[uuid.UUID]*player.Player),
		Players:      make(map[uuid.UUID]*player.Player),
		playerSlice:  []*player.Player{},
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

func (g *Game) assignRoles() error {
	if g.Roleset == nil {
		return fmt.Errorf("game: no roleset defined") // TODO
	}
	if len(g.Players) != len(g.Roleset.Roles) {
		return &PlayerCountError{Roleset: g.Roleset, PlayerCount: len(g.Players)}
	}

	for playerKey, roleKey := range rand.Perm(len(g.Players)) {
		p := g.playerSlice[playerKey]
		r := g.Roleset.Roles[roleKey]
		p.SetRole(r)
	}

	return nil
}

func (g *Game) processN0() {
	for _, p := range g.Players {
		if p.Role.CanViewForMax() && p.Role.HasRandomN0Clear() {
			view := g.randomClear(p, func(r *role.Role) bool { return r.ViewForMaxEvil() })
			p.Views = append(p.Views, &player.View{
				Player:    view,
				Attribute: role.MaxEvilAttribute,
				Hit:       false,
				GamePhase: g.Phase,
			})
		}

		if p.Role.CanViewForSeer() && p.Role.HasRandomN0Clear() {
			view := g.randomClear(p, func(r *role.Role) bool { return r.ViewForSeer() })
			p.Views = append(p.Views, &player.View{
				Player:    view,
				Attribute: role.SeerAttribute,
				Hit:       false,
				GamePhase: g.Phase,
			})
		}

		if p.Role.CanViewForAux() && p.Role.HasRandomN0Clear() {
			view := g.randomClear(p, func(r *role.Role) bool { return r.ViewForAuxEvil() })
			p.Views = append(p.Views, &player.View{
				Player:    view,
				Attribute: role.AuxEvilAttribute,
				Hit:       false,
				GamePhase: g.Phase,
			})
		}

		if p.Role.KnowsMaxes() {
			for _, m := range g.AliveMaxEvils() {
				if m != p {
					p.Views = append(p.Views, &player.View{
						Player:    m,
						Attribute: role.MaxEvilAttribute,
						Hit:       true,
						GamePhase: g.Phase,
					})
				}
			}
		}

	}
}

func (g *Game) randomClear(p *player.Player, test func(*role.Role) bool) *player.Player {
	for _, i := range rand.Perm(len(g.Players)) {
		view := g.playerSlice[i]
		if view == p {
			continue
		}
		if !test(view.Role) {
			return view
		}
	}
	return nil // should be impossible
}

func (g *Game) StartGame() error {
	if g.state > Setup {
		return &StateError{NeedState: Setup, InState: g.state}
	}

	// fill in the slice of players now that we have them all
	for _, p := range g.Players {
		g.playerSlice = append(g.playerSlice, p)
		g.AlivePlayers[p.ID] = p
	}

	if err := g.assignRoles(); err != nil {
		return err
	}

	g.processN0()

	g.state = Running
	return nil
}

func (g *Game) AliveMaxEvils() []*player.Player {
	maxes := []*player.Player{}
	for _, p := range g.AlivePlayers {
		if p.Role.IsMaxEvil() {
			maxes = append(maxes, p)
		}
	}
	return maxes
}
