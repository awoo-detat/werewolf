package game

import (
	"testing"

	"github.com/awoo-detat/werewolf/player"
	"github.com/awoo-detat/werewolf/role/roleset"

	"github.com/stretchr/testify/assert"
)

func TestInitialization(t *testing.T) {
	l := player.NewPlayer()

	g := NewGame(l)

	assert.NotEmpty(t, g.ID)
	assert.Equal(t, l, g.Leader)
	assert.Len(t, g.Players, 1)
	assert.Nil(t, g.Roleset)
	assert.Equal(t, g.State, Setup)
	assert.Equal(t, g.Phase, 0)
}

func TestAddingPlayers(t *testing.T) {
	l := player.NewPlayer()
	g := NewGame(l)
	p := player.NewPlayer()

	g.AddPlayer(p)

	assert.Len(t, g.Players, 2)
}

func TestSetRoleset(t *testing.T) {
	l := player.NewPlayer()
	g := NewGame(l)
	rs := roleset.VanillaFiver()

	err := g.ChooseRoleset("Vanilla Fiver")

	assert.Equal(t, g.Roleset, rs)
	assert.Nil(t, err)
}

func TestRolesetMustExist(t *testing.T) {
	l := player.NewPlayer()
	g := NewGame(l)
	rs := roleset.VanillaFiver()

	err := g.ChooseRoleset("dkjjfkfwegfwegy")

	assert.Equal(t, g.Roleset, rs)
	assert.Error(t, err)
}

func TestGameCannotHaveStartedWhenChoosingRoleset(t *testing.T) {
}

func TestCannotStartGameWithoutRoleset(t *testing.T) {
}

func TestCannotStartGameWithMismatchedPlayerCountAndRoleset(t *testing.T) {
}
