package game

import (
	"fmt"
	"log/slog"
	"math"
	"math/rand"

	"github.com/awoo-detat/werewolf/player"
	"github.com/awoo-detat/werewolf/role"
	"github.com/awoo-detat/werewolf/role/roleset"
	"github.com/awoo-detat/werewolf/tally"

	"github.com/google/uuid"
)

type Game struct {
	ID           uuid.UUID
	Leader       *player.Player
	VotingMethod VotingMethod
	AlivePlayers map[uuid.UUID]*player.Player
	Players      map[uuid.UUID]*player.Player
	playerSlice  []*player.Player
	Roleset      *roleset.Roleset
	state        GameState
	Phase        int
	Tally        *tally.Tally
	Winner       role.PlayerType
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
		AlivePlayers: make(map[uuid.UUID]*player.Player),
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
	slog.Info("player added", "player", p)
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
	slog.Info("roleset chosen", "roleset", rs)
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

func (g *Game) Start() error {
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

	g.Tally = tally.New(g.playerSlice)
	g.state = Running
	return nil
}

func (g *Game) Vote(from *player.Player, to *player.Player) error {
	if g.state != Running {
		return &StateError{NeedState: Running, InState: g.state}
	}
	if g.IsNight() {
		return &PhaseError{GamePhase: g.Phase}
	}

	g.Tally.Vote(from, to)

	switch g.VotingMethod {
	case InstaKill:
		g.checkForInstaKillDayEnd()
	case InstaKillWithDelay:
		return fmt.Errorf("InstaKillWithDelay is not supported")
	case Timed:
		return fmt.Errorf("Timed is not supported")
	}

	return nil
}

func (g *Game) checkForInstaKillDayEnd() {
	leader := g.Tally.List[0]
	// TODO? if there's an even number this is first-to-half...
	need := math.Ceil(float64(len(g.AlivePlayers)) / 2)
	have := float64(len(leader.Votes))
	if have < need {
		slog.Info("day not over", "have", have, "need", need)
		return
	}

	g.KillPlayer(leader.Player)
	if g.state == Running {
		g.Phase++
	}
}

func (g *Game) KillPlayer(p *player.Player) {
	killed := p.Role.Kill()
	if !killed {
		return
	}
	delete(g.AlivePlayers, p.ID)

	maxes, nonmaxes := g.AlivePlayersByType()
	parity := g.Parity()
	equality := len(maxes) == len(nonmaxes)
	switch {
	case len(maxes) == 0:
		g.EndGame(role.Good)
	case parity < 0:
		g.EndGame(role.Evil)
	case parity == 0 && equality:
		g.EndGame(role.Evil)
	case equality:
		// due to a hunter, evil loses
		// TODO: what will/should happen with ancient WW vs hunter?
		g.EndGame(role.Good)
	}
}

func (g *Game) Parity() int {
	parity := 0
	for _, p := range g.AlivePlayers {
		parity += p.Role.Parity
	}
	slog.Info("parity calculated", "parity", parity)
	return parity
}

func (g *Game) EndGame(winner role.PlayerType) {
	slog.Info("game over", "winner", winner)
	g.state = Finished
	g.Winner = winner
}

func (g *Game) AliveMaxEvils() []*player.Player {
	maxes, _ := g.AlivePlayersByType()
	return maxes
}

func (g *Game) AlivePlayersByType() (maxes []*player.Player, nonmaxes []*player.Player) {
	for _, p := range g.AlivePlayers {
		if p.Role.IsMaxEvil() {
			maxes = append(maxes, p)
		} else {
			nonmaxes = append(nonmaxes, p)
		}
	}
	return
}

func (g *Game) IsDay() bool {
	return g.Phase%2 == 0
}

func (g *Game) IsNight() bool {
	return !g.IsDay()
}
