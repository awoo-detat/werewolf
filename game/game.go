package game

import (
	"fmt"
	"log/slog"
	"math"
	"math/rand"
	"slices"

	"github.com/awoo-detat/werewolf/gamechannel"
	"github.com/awoo-detat/werewolf/gamechannel/server"
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
	gameChannel  gamechannel.GameChannel
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
		Players:      make(map[uuid.UUID]*player.Player),
		VotingMethod: InstaKill,
		AlivePlayers: make(map[uuid.UUID]*player.Player),
		nightActions: make(map[*player.Player]*player.FingerPoint),
		playerSlice:  []*player.Player{},
		gameChannel:  make(gamechannel.GameChannel),
	}
	g.AddPlayer(p)

	go g.ListenToGameChannel()
	return g
}

func (g *Game) SetLeader(p *player.Player) {
	slog.Info("setting leader", "player", p)
	g.Leader = p
	p.Message(server.RolesetList, roleset.List())
}

func (g *Game) State() GameState {
	return g.state
}

func (g *Game) AddPlayer(p *player.Player) {
	if len(g.Players) == 0 {
		g.SetLeader(p)
	}
	p.SetGameChannel(g.gameChannel)
	g.Players[p.ID] = p
	slog.Info("player added", "player", p)
	g.Broadcast(server.PlayerJoin, p)
	g.BroadcastPlayerList()
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
	g.Broadcast(server.RolesetSelected, rs)
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
		slog.Info("assigning role", "player", p, "role", r)
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
					p.AddView(player.NewAttributeView(m, role.MaxEvilAttribute, true, g.Phase))
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
			slog.Info("random clear obtained", "for", p, "clear", view)
			return view
		}
	}
	return nil // should be impossible
}

func (g *Game) Start() error {
	if g.state > Setup {
		return &StateError{NeedState: Setup, InState: g.state}
	}
	slog.Info("starting game")

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
	slog.Info("new phase", "phase", g.Phase)

	// new day new me, reset everything
	if g.IsDay() {
		g.Broadcast(server.PhaseChanged, &server.Phase{Phase: server.Day, Count: g.Phase})
		g.Tally = tally.New(g.playerSlice)
		g.Broadcast(server.TallyChanged, g.Tally)
		g.nightActions = make(map[*player.Player]*player.FingerPoint)
	} else {
		g.Broadcast(server.PhaseChanged, &server.Phase{Phase: server.Night, Count: g.Phase})
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
	g.Broadcast(server.TallyChanged, g.Tally)

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
	slog.Info("killing player", "player", p)
	killed := p.Role.Kill()
	if !killed {
		return
	}
	delete(g.AlivePlayers, p.ID)
	g.RevealPlayer(p)

	maxes, nonmaxes := g.AlivePlayersByType()
	parity := g.Parity()
	equality := len(maxes) == len(nonmaxes)
	slog.Info("checking for game end", "parity", parity, "equality", equality)
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

	slog.Info("setting night action", "fingerpoint", fp)

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
	slog.Info("processing night actions", "phase", g.Phase)
	var nk *player.Player
	for _, fp := range g.nightActions {
		var view *player.View
		switch {
		case fp.From.Role.CanViewForMax():
			view = player.NewAttributeView(fp.To, role.MaxEvilAttribute, fp.To.Role.ViewForMaxEvil(), g.Phase)
		case fp.From.Role.CanNightKill():
			nk = fp.To
		case fp.From.Role.CanViewForSeer():
			view = player.NewAttributeView(fp.To, role.SeerAttribute, fp.To.Role.ViewForSeer(), g.Phase)
		case fp.From.Role.CanViewForAux():
			view = player.NewAttributeView(fp.To, role.AuxEvilAttribute, fp.To.Role.ViewForAuxEvil(), g.Phase)
		default:
			// TODO keep track of "most suspicious" (#4)
		}
		if view != nil {
			fp.From.AddView(view)
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
	slog.Info("revealing player", "player", p, "role", p.Role)
	g.BroadcastView(player.NewRoleView(p, p.Role, g.Phase))
}

// BroadcastView sends a view to every player, alive and dead. It
// is primarily used for revealing the roles of dead players.
func (g *Game) BroadcastView(v *player.View) {
	for _, p := range g.Players {
		p.AddView(v)
	}
}

func (g *Game) Broadcast(t server.MessageType, payload interface{}) {
	for _, p := range g.Players {
		p.Message(t, payload)
	}
}

// probably needs to be better but hackathon
func (g *Game) alivePlayerList() []*player.Player {
	var list []*player.Player
	for _, p := range g.Players {
		list = append(list, p)
	}
	return list
}

func (g *Game) BroadcastPlayerList() {
	g.Broadcast(server.AlivePlayerList, g.alivePlayerList())
}

func (g *Game) ListenToGameChannel() {
	for {
		slog.Info("waiting for message on game channel...")
		activity := <-g.gameChannel

		switch activity.Type {
		case gamechannel.SetName:
			g.BroadcastPlayerList()
		case gamechannel.SetRoleset:
			g.ChooseRoleset(activity.Value.(string))
		case gamechannel.Reconnect:
			p := g.Players[activity.From]
			p.Message(server.AlivePlayerList, g.alivePlayerList())
			if g.Leader == p && g.state == Setup {
				slog.Info("sending roleset list to leader", "player", p)
				p.Message(server.RolesetList, roleset.List())
			}
		case gamechannel.Vote:
			from := g.Players[activity.From]
			to := g.Players[activity.Value.(uuid.UUID)]
			g.Vote(&player.FingerPoint{From: from, To: to})
		case gamechannel.NightAction:
			from := g.Players[activity.From]
			to := g.Players[activity.Value.(uuid.UUID)]
			g.SetNightAction(&player.FingerPoint{From: from, To: to})
		case gamechannel.Quit:
			p := g.Players[activity.From]
			delete(g.Players, activity.From)
			g.Broadcast(server.PlayerLeave, p)
			g.BroadcastPlayerList()
		}
	}
}
