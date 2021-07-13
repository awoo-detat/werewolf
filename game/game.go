package game

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/awoo-detat/werewolf/player"
	"github.com/awoo-detat/werewolf/role"
	"github.com/awoo-detat/werewolf/role/roleset"
	"github.com/awoo-detat/werewolf/tally"
)

const (
	NotRunning = iota
	Setup      = iota
	Running    = iota
	Finished   = iota
)

// A Game is a representation of a game of Werewolf.
type Game struct {
	Players          map[string]player.Player `json:"-"`
	PlayerList       []player.Player          `json:"players"`
	Roleset          *roleset.Roleset         `json:"roleset"`
	Tally            []*tally.TallyItem       `json:"tally"`
	State            int                      `json:"game_state"`
	Phase            int                      `json:"phase"`
	NightActionQueue []*FingerPoint           `json:"-"`
	votes            map[player.Player]string
}

// New creates a new Game.
func New() *Game {
	game := &Game{
		Players: make(map[string]player.Player),
		State:   NotRunning,
		Phase:   0,
	}
	return game
}

// UpdatePlayerList populates the list of players with those who are alive.
func (g *Game) UpdatePlayerList() {
	var l []player.Player
	for _, p := range g.Players {
		if p.Role().Name == "" || p.Role().Alive {
			l = append(l, p)
		}
	}
	g.PlayerList = l
}

// AddPlayer adds a player to the game. If there are currently no players, it will make
// the player the game's Leader.
func (g *Game) AddPlayer(p player.Player) error {
	if g.State == NotRunning {
		g.State = Setup
	}
	if g.State != Setup {
		return fmt.Errorf("cannot add player: game is not in setup phase")
	}

	log.Printf("new player: %s", p.GetIdentifier())
	if len(g.Players) == 0 {
		log.Printf("setting leader: %s", p.GetIdentifier())
		p.SetLeader()
	}

	g.Players[p.GetID()] = p
	g.UpdatePlayerList()

	if g.ShouldStart() {
		g.Start()
	}
	return nil
}

// ShouldStart determines if the game is ready to begin.
func (g *Game) ShouldStart() bool {
	if g.Roleset != nil {
		log.Printf("%v/%v players", len(g.Players), len(g.Roleset.Roles))
	}
	return g.Roleset != nil && len(g.Players) == len(g.Roleset.Roles)
}

// Start begins the game.
func (g *Game) Start() error {
	if g.State != Setup {
		return fmt.Errorf("cannot start game: game is not in setup phase")
	}

	log.Println("== role selection ==")
	rand.Seed(time.Now().UnixNano())
	roleOrder := rand.Perm(len(g.Roleset.Roles))
	for k, v := range g.PlayerList {
		r := g.Roleset.Roles[roleOrder[k]]
		log.Printf("%s: %s", v.GetIdentifier(), r.Name)
		g.Players[v.GetID()].SetRole(r)
	}

	log.Printf("== starting game ==")

	g.State = Running
	g.NextPhase()
	g.ProcessStartActionQueue()

	return nil
}

// NextPhase moves the game to its next phase: from day to night or night to day. If a
// victory condition has been satisfied, it will end the game.
func (g *Game) NextPhase() {
	g.UpdatePlayerList()
	maxes := g.AliveMaxEvils()
	g.NightActionQueue = []*FingerPoint{}
	g.votes = make(map[player.Player]string)
	g.RebuildTally()

	log.Printf("alive maxes: %v", maxes)
	log.Printf("parity: %v", g.Parity())

	if len(maxes) == 0 {
		g.End(role.Good)
		return
	} else if g.Parity() <= 0 {
		g.End(role.Evil)
		return
	} else if len(g.PlayerList) == 2 {
		g.End(role.Good)
		return
	}

	g.Phase++

	log.Printf("== game is now on phase %v ==", g.Phase)
}

// Parity returns the parity count of the game: typically this will be max evils counting
// for -1 and all other roles counting for 1, though that can change depending on the role.
func (g *Game) Parity() int {
	parity := 0
	for _, p := range g.Players {
		if p.Role().Alive {
			parity += p.Role().Parity
		}
	}
	return parity
}

// RebuildTally calculates the current tally based on individual votes.
// If there are no votes, it essentially clears it.
func (g *Game) RebuildTally() {
	t := make(map[player.Player][]player.Player)

	for _, p := range g.PlayerList {
		t[p] = []player.Player{}
	}

	for from, to := range g.votes {
		t[g.Players[to]] = append(t[g.Players[to]], from)
	}

	list := []*tally.TallyItem{}
	for c, v := range t {
		item := tally.Item(c, v)
		list = append(list, item)
	}
	g.Tally = list
}

// SetRoleset chooses the roleset for the game. It can only be done before the game starts.
func (g *Game) SetRoleset(r *roleset.Roleset) error {
	if g.State != Setup {
		return fmt.Errorf("cannot set roleset: game is not in 'not running' phase")
	} else if len(r.Roles) < len(g.Players) {
		return fmt.Errorf("roleset %s only has %v roles; %v players in lobby", r.Name, len(r.Roles), len(g.Players))
	}

	g.Roleset = r

	if g.ShouldStart() {
		g.Start()
	}
	return nil
}

// Vote handles a vote from a player to another. The voted player's ID is given. Currently
// a game only supports instakill, so if a majority is voting for one player then it will
// end the day.
func (g *Game) Vote(from player.Player, to string) error {
	if !g.Day() {
		err := fmt.Errorf("vote failed; not day")
		log.Println(err)
		return err
	}

	g.votes[from] = to
	g.RebuildTally()

	ousted := g.VotedOut()
	if ousted != nil {
		g.EndDay(ousted)
	}
	return nil
}

// EndDay finalizes a vote for a player and ends the day.
//
// TODO this should have a better name
func (g *Game) EndDay(ousted player.Player) {
	// if the player died, regenerate our list
	if !ousted.Role().Kill() {
		g.UpdatePlayerList()
	}
	revealed := ousted.Reveal()

	log.Printf("killed: %s (%s)", revealed.Name, revealed.RoleName)
	g.NextPhase()
}

// AliveMaxEvils returns a list of the identifiers of the alive max evils.
func (g *Game) AliveMaxEvils() []string {
	var maxes []string
	for _, p := range g.PlayerList {
		if p.Role().IsMaxEvil() {
			maxes = append(maxes, p.GetIdentifier())
		}
	}
	return maxes
}

// End concludes the game, declaring a victor.
func (g *Game) End(victor int) {
	log.Printf("== game over. victor: %v ==", victor)
	g.State = Finished
}

// VotedOut determines the player who has been voted out, if there is one.
func (g *Game) VotedOut() player.Player {
	// if an even number, first to half. otherwise, 50%+1
	threshold := len(g.PlayerList) / 2
	if len(g.PlayerList)%2 == 1 {
		threshold++
	}

	votes := 0
	var ousted player.Player
	for _, item := range g.Tally {
		votes += len(item.Votes)
		if len(item.Votes) >= threshold {
			ousted = item.Candidate
		}
	}
	if votes == len(g.PlayerList) && ousted != nil {
		return ousted
	}
	return nil
}

// QueueNightAction adds a FingerPoint, potentially replacing one if it already exists.
// If we have a number of night actions equal to the number of alive players, process the
// queue.
func (g *Game) QueueNightAction(fp *FingerPoint) {
	for i, a := range g.NightActionQueue {
		if a.From.GetID() == fp.From.GetID() {
			log.Printf("%s: replacing night action", fp.From.GetIdentifier())
			g.NightActionQueue = append(g.NightActionQueue[:i], g.NightActionQueue[i+1:]...)
			break
		}
	}
	g.NightActionQueue = append(g.NightActionQueue, fp)
	log.Printf("%v / %v actions", len(g.NightActionQueue), len(g.PlayerList))
	if len(g.NightActionQueue) >= len(g.PlayerList) {
		g.ProcessNightActionQueue()
	}
}

// ProcessStartActionQueue handles all of the random N0s, etc, to begin the game.
func (g *Game) ProcessStartActionQueue() {
	for _, p := range g.Players {
		_ = g.StartAction(p)
		// TODO use the result of that action (PM the user, etc)
	}
}

// ProcessNightActionQueue iterates over the night actions and handles all of them.
func (g *Game) ProcessNightActionQueue() {
	//var death *player.Revealed
	for _, action := range g.NightActionQueue {
		// TODO message
		_ = g.NightAction(action)
		/*
			if result.PlayerMessage != "" {
				action.From.Message(message.Private, result.PlayerMessage)
			}

			// TODO maybe more than one death?
			if result.Killed != nil && death == nil && !result.Killed.Role().Kill() {
				death = result.Killed.Reveal()
				result.Killed.Message(message.Dead, "")
			}
		*/
	}

	g.NextPhase()

	/*
		if death != nil {
			g.Broadcast(message.Targeted, death)
		}
	*/
}

// Day returns whether or not it is daytime.
func (g *Game) Day() bool {
	return g.Phase%2 == 1
}

// RemovePlayer removes them from the game. This is different from them being no longer
// alive.
func (g *Game) RemovePlayer(id string) {
	delete(g.Players, id)
	g.UpdatePlayerList()

	if len(g.PlayerList) == 0 {
		g.Reset()
	}
}

// Reset brings us back to square one.
func (g *Game) Reset() {
	for id := range g.Players {
		delete(g.Players, id)
	}
	g.PlayerList = []player.Player{}
	g.Roleset = nil
	g.votes = make(map[player.Player]string)
	g.Tally = []*tally.TallyItem{}
	g.State = NotRunning
	g.Phase = 0
	g.NightActionQueue = []*FingerPoint{}
}
