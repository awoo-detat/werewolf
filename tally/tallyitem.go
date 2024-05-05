package tally

import (
	"slices"

	"github.com/awoo-detat/werewolf/player"
	"github.com/awoo-detat/werewolf/vote"
)

// A TallyItem represents a line on a tally: a player and a list
// of votes.
type TallyItem struct {
	Player *player.Player
	Votes  []*vote.Vote
}

func NewTallyItem(p *player.Player) *TallyItem {
	return &TallyItem{
		Player: p,
		Votes:  []*vote.Vote{},
	}
}

func (i *TallyItem) RemoveVote(v *vote.Vote) {
	i.Votes = slices.DeleteFunc(i.Votes, func(tallyVote *vote.Vote) bool {
		return tallyVote == v
	})
}

func (i *TallyItem) AddVote(v *vote.Vote) {
	i.Votes = append(i.Votes, v)
}
