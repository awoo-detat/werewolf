package game

import (
	"fmt"
	"log/slog"
	"math"
	"math/rand"
	"slices"

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
	nightActions map[*player.Player]*player.FingerPoint
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
		nightActions: make(map[*player.Player]*player.FingerPoint),
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

// this selects random clears for those that get them, and informs
// the roles that know maxes of those players
func (g *Game) processN0() {
	for _, p := range g.Players {
		if p.Role.CanViewForMax() && p.Role.HasRandomN0Clear() {
			view := g.randomClear(p, func(r *role.Role) bool { return r.ViewForMaxEvil() })
			g.nightActions[p] = &player.FingerPoint{From: p, To: view}
		}

		if p.Role.CanViewForSeer() && p.Role.HasRandomN0Clear() {
			view := g.randomClear(p, func(r *role.Role) bool { return r.ViewForSeer() })
			g.nightActions[p] = &player.FingerPoint{From: p, To: view}
		}

		if p.Role.CanViewForAux() && p.Role.HasRandomN0Clear() {
			view := g.randomClear(p, func(r *role.Role) bool { return r.ViewForAuxEvil() })
			g.nightActions[p] = &player.FingerPoint{From: p, To: view}
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
	g.processNightActions()
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

	g.state = Running
	g.processN0()

	return nil
}

func (g *Game) nextPhase() {
	g.Phase++
	// new day new me, reset everything
	if g.IsDay() {
		g.Tally = tally.New(g.playerSlice)
		g.nightActions = make(map[*player.Player]*player.FingerPoint)
	}
}

func (g *Game) Vote(fp *player.FingerPoint) error {
	if g.state != Running {
		return &StateError{NeedState: Running, InState: g.state}
	}
	if g.IsNight() {
		return &PhaseError{GamePhase: g.Phase}
	}

	if !fp.From.Role.Alive {
		return fmt.Errorf("%s is dead and cannot vote", fp.From)
	}
	if !fp.To.Role.Alive {
		return fmt.Errorf("%s is dead and cannot be voted for", fp.To)
	}

	g.Tally.Vote(fp)

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
		g.nextPhase()
	}
}

func (g *Game) KillPlayer(p *player.Player) {
	killed := p.Role.Kill()
	if !killed {
		return
	}
	delete(g.AlivePlayers, p.ID)
	g.RevealPlayer(p)

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
	return g.Phase%2 == 1
}

func (g *Game) IsNight() bool {
	return !g.IsDay()
}

func (g *Game) SetNightAction(fp *player.FingerPoint) error {
	if !fp.From.Role.Alive {
		return fmt.Errorf("%s is dead and cannot have a night action", fp.From)
	}
	if !fp.To.Role.Alive {
		return fmt.Errorf("%s is dead and cannot be targeted by a night action", fp.To)
	}

	// this allows you to change your mind and choose someone else
	g.nightActions[fp.From] = fp
	neededPlayers := g.alivePlayersWithNightActions()
	neededPlayers = slices.DeleteFunc(neededPlayers, func(p *player.Player) bool {
		_, ok := g.nightActions[p]
		return ok
	})
	if len(neededPlayers) == 0 {
		g.processNightActions()
	} else {
		slog.Info("still need night actions", "needed", neededPlayers)
	}
	return nil
}

func (g *Game) alivePlayersWithNightActions() []*player.Player {
	players := []*player.Player{}
	for _, p := range g.AlivePlayers {
		switch {
		case p.Role.CanViewForMax():
			players = append(players, p)
		case p.Role.CanNightKill():
			players = append(players, p)
		case p.Role.CanViewForSeer():
			players = append(players, p)
		case p.Role.CanViewForAux():
			players = append(players, p)
		}
	}
	return players
}

// this assumes you will only ever have one night action...
func (g *Game) processNightActions() {
	var nk *player.Player
	for _, fp := range g.nightActions {
		switch {
		case fp.From.Role.CanViewForMax():
			fp.From.Views = append(fp.From.Views, &player.View{
				Player:    fp.To,
				Attribute: role.MaxEvilAttribute,
				Hit:       fp.To.Role.ViewForMaxEvil(),
				GamePhase: g.Phase,
			})
		case fp.From.Role.CanNightKill():
			nk = fp.To
		case fp.From.Role.CanViewForSeer():
			fp.From.Views = append(fp.From.Views, &player.View{
				Player:    fp.To,
				Attribute: role.SeerAttribute,
				Hit:       fp.To.Role.ViewForSeer(),
				GamePhase: g.Phase,
			})
		case fp.From.Role.CanViewForAux():
			fp.From.Views = append(fp.From.Views, &player.View{
				Player:    fp.To,
				Attribute: role.AuxEvilAttribute,
				Hit:       fp.To.Role.ViewForAuxEvil(),
				GamePhase: g.Phase,
			})
		default:
			// TODO keep track of "most suspicious" (#4)
		}
	}

	// TODO this prioritizes the "last" wolf if there are multiple
	if nk != nil {
		g.KillPlayer(nk)
	}
	if g.state == Running {
		g.nextPhase()
	}
}

func (g *Game) RevealPlayer(p *player.Player) {
	// TODO #5
}
