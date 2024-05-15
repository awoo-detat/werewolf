package vote

import (
	"time"

	"github.com/awoo-detat/werewolf/player"
)

type Vote struct {
	Candidate *player.Player `json:"candidate"`
	Voter     *player.Player `json:"voter"`
	Timestamp time.Time      `json:"timestamp"`
}

func New(fp *player.FingerPoint) *Vote {
	return &Vote{
		Candidate: fp.To,
		Voter:     fp.From,
		Timestamp: time.Now(),
	}
}
