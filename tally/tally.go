package tally

import (
	"sort"

	"github.com/awoo-detat/werewolf/player"
	"github.com/awoo-detat/werewolf/vote"
)

// A Tally is a list of players and the votes they have received.
// It is in descending order by number of votes and by longest
// held last vote (LHLV).
type Tally struct {
	// Item is a map of votes ordered by the person being voted for.
	List []*TallyItem
	// Inverted is a map ordered by the player doing the voting.
	voteMap     map[*player.Player]*TallyItem
	Inverted    map[*player.Player]*vote.Vote
	playerCount int
}

func New(players []*player.Player) *Tally {
	t := &Tally{
		List:        []*TallyItem{},
		voteMap:     make(map[*player.Player]*TallyItem),
		Inverted:    make(map[*player.Player]*vote.Vote),
		playerCount: len(players),
	}
	sort.Slice(players, func(i, j int) bool {
		return players[i].Name < players[j].Name
	})
	for _, p := range players {
		ti := NewTallyItem(p)
		t.List = append(t.List, ti)
		t.voteMap[p] = ti
		t.Inverted[p] = nil
	}
	return t
}

func (t *Tally) Vote(fp *player.FingerPoint) {
	// if they've voted for anyone before, remove it from the tally
	if current := t.Inverted[fp.From]; current != nil {
		t.voteMap[current.Candidate].RemoveVote(current)
	}
	v := vote.New(fp)
	// add to the tally
	t.voteMap[fp.To].AddVote(v)
	// update the inverted tally
	t.Inverted[fp.From] = v

	sort.Slice(t.List, func(i, j int) bool {
		return len(t.List[i].Votes) > len(t.List[j].Votes)
	})
}

func (t *Tally) Unvote(from *player.Player) {
	v := t.Inverted[from]
	if v == nil {
		return
	}
	t.voteMap[v.Candidate].RemoveVote(v)
	t.Inverted[from] = nil
}
