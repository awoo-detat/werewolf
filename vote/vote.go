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

func New(fp *player.FingerPoint) *Vote {
	return &Vote{
		Candidate: fp.To,
		Voter:     fp.From,
		Timestamp: time.Now(),
	}
}
