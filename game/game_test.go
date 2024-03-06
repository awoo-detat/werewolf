package game

import (
	"testing"

	"github.com/awoo-detat/werewolf/player"

	"github.com/stretchr/testify/assert"
)

func TestInitialization(t *testing.T) {
	l := player.NewPlayer()

	g := NewGame(l)

	assert.NotEmpty(t, g.ID)
	assert.Equal(t, l, g.Leader)
	assert.Len(t, g.Players, 1)
}

func TestAddingPlayers(t *testing.T) {
	l := player.NewPlayer()
	g := NewGame(l)
	p := player.NewPlayer()

	g.AddPlayer(p)

	assert.Len(t, g.Players, 2)
}
