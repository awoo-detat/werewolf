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
	assert.Equal(t, Setup, g.State())
	assert.Equal(t, 0, g.Phase)
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

	assert.Equal(t, rs, g.Roleset)
	assert.Nil(t, err)
}

func TestRolesetMustExist(t *testing.T) {
	l := player.NewPlayer()
	g := NewGame(l)

	err := g.ChooseRoleset("dkjjfkfwegfwegy")

	assert.Nil(t, g.Roleset)
	assert.Error(t, err)
}

func TestCanChangeRolesetWhenInSetup(t *testing.T) {
	l := player.NewPlayer()
	g := NewGame(l)
	assert.Equal(t, Setup, g.State())
	err := g.ChooseRoleset("Vanilla Fiver")
	assert.Nil(t, err)

	err = g.ChooseRoleset("Basic Niner")
	assert.Nil(t, err)
	err = g.ChooseRoleset("Fast Fiver")
	assert.Nil(t, err)

}

func TestGameCannotHaveStartedWhenChoosingRoleset(t *testing.T) {
	p1 := player.NewPlayer()
	p2 := player.NewPlayer()
	p3 := player.NewPlayer()
	p4 := player.NewPlayer()
	p5 := player.NewPlayer()
	g := NewGame(p1)
	assert.Equal(t, Setup, g.State())
	g.AddPlayer(p2)
	g.AddPlayer(p3)
	g.AddPlayer(p4)
	g.AddPlayer(p5)
	err := g.ChooseRoleset("Vanilla Fiver")
	assert.Nil(t, err)
	err = g.StartGame()
	assert.Nil(t, err)

	err = g.ChooseRoleset("Fast Fiver")

	assert.Error(t, err)
}

func TestCannotStartWhenNotInSetup(t *testing.T) {
	p1 := player.NewPlayer()
	p2 := player.NewPlayer()
	p3 := player.NewPlayer()
	p4 := player.NewPlayer()
	p5 := player.NewPlayer()
	g := NewGame(p1)
	assert.Equal(t, Setup, g.State())
	g.AddPlayer(p2)
	g.AddPlayer(p3)
	g.AddPlayer(p4)
	g.AddPlayer(p5)
	err := g.ChooseRoleset("Vanilla Fiver")
	assert.Nil(t, err)
	err = g.StartGame()
	assert.Nil(t, err)

	err = g.StartGame()

	assert.Error(t, err)
}

func TestCannotStartGameWithoutRoleset(t *testing.T) {
	l := player.NewPlayer()
	g := NewGame(l)

	err := g.StartGame()

	assert.Error(t, err)
}

func TestCannotStartGameWithMismatchedPlayerCountAndRoleset(t *testing.T) {
	l := player.NewPlayer()
	g := NewGame(l)
	err := g.ChooseRoleset("Vanilla Fiver")
	assert.Nil(t, err)

	err = g.StartGame()

	assert.Error(t, err) // TODO match on the error?
}

// really more of an integration test, maybe!
func TestFiver(t *testing.T) {
	p1 := player.NewPlayer()
	p2 := player.NewPlayer()
	p3 := player.NewPlayer()
	p4 := player.NewPlayer()
	p5 := player.NewPlayer()

	g := NewGame(p1)
	assert.Equal(t, Setup, g.State())

	g.AddPlayer(p2)
	g.AddPlayer(p3)
	g.AddPlayer(p4)
	g.AddPlayer(p5)

	err := g.ChooseRoleset("Vanilla Fiver")
	assert.Nil(t, err)

	err = g.StartGame()
	assert.Nil(t, err)
	assert.Equal(t, Running, g.State())
}
