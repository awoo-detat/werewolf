package tally

import (
	"github.com/awoo-detat/werewolf/player"
)

// A ShortTally is a more human-readable representation of the game's tally.
type ShortTally struct {
	Candidate string   `json:"candidate"`
	Votes     []string `json:"votes"`
}

// Short creates a new ShortTally from a list of TallyItems.
func Short(verbose []*TallyItem) []*ShortTally {
	var short []*ShortTally
	for _, item := range verbose {
		c := item.Candidate.GetID()
		votes := []string{}
		for _, v := range item.Votes {
			votes = append(votes, v.GetID())
		}
		short = append(short, &ShortTally{
			Candidate: c,
			Votes:     votes,
		})
	}

	return short
}

// A TallyItem represents a player and all the votes they have against them.
type TallyItem struct {
	Candidate player.Player   `json:"candidate"`
	Votes     []player.Player `json:"votes"`
}

// Item creates a new TallyItem from a player and a list of players (who are voting for
// the player).
//
// TODO this doesn't seem like a great name?
func Item(c player.Player, v []player.Player) *TallyItem {
	return &TallyItem{
		Candidate: c,
		Votes:     v,
	}
}
