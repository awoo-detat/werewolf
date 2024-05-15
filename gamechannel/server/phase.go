package server

import (
	"fmt"
)

type GamePhase string

const (
	Day   GamePhase = "day"
	Night           = "Night"
)

type Phase struct {
	Phase GamePhase `json:"phase"`
	Count int       `json:"count"`
}

func (p *Phase) String() string {
	return fmt.Sprintf("%s %v", p.Phase, p.Count)
}
