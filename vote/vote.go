package vote

import (
	"time"

	"github.com/awoo-detat/werewolf/player"
)

type Vote struct {
	Candidate *player.Player
	Voter     *player.Player
	Timestamp time.Time
}

func New(from *player.Player, to *player.Player) *Vote {
	return &Vote{
		Candidate: to,
		Voter:     from,
		Timestamp: time.Now(),
	}
}
