package game

import (
	"testing"

	"github.com/awoo-detat/werewolf/player"
	"github.com/awoo-detat/werewolf/role"
	"github.com/awoo-detat/werewolf/role/roleset"

	"github.com/stretchr/testify/assert"
)

func TestInitialization(t *testing.T) {
	assert := assert.New(t)
	l := player.NewPlayer()

	g := NewGame(l)

	assert.NotEmpty(g.ID)
	assert.Equal(l, g.Leader)
	assert.Len(g.Players, 1)
	assert.Nil(g.Roleset)
	assert.Equal(Setup, g.State())
	assert.Equal(0, g.Phase)
}

func TestAddingPlayers(t *testing.T) {
	l := player.NewPlayer()
	g := NewGame(l)
	p := player.NewPlayer()

	g.AddPlayer(p)

	assert.Len(t, g.Players, 2)
}

func TestSetRoleset(t *testing.T) {
	assert := assert.New(t)
	l := player.NewPlayer()
	g := NewGame(l)
	rs := roleset.VanillaFiver()

	err := g.ChooseRoleset("Vanilla Fiver")

	assert.Equal(rs, g.Roleset)
	assert.Nil(err)
}

func TestRolesetMustExist(t *testing.T) {
	assert := assert.New(t)
	l := player.NewPlayer()
	g := NewGame(l)

	err := g.ChooseRoleset("dkjjfkfwegfwegy")

	assert.Nil(g.Roleset)
	assert.Error(err)
}

func TestCanChangeRolesetWhenInSetup(t *testing.T) {
	assert := assert.New(t)
	l := player.NewPlayer()
	g := NewGame(l)
	assert.Equal(Setup, g.State())
	err := g.ChooseRoleset("Vanilla Fiver")
	assert.Nil(err)

	err = g.ChooseRoleset("Basic Niner")
	assert.Nil(err)
	err = g.ChooseRoleset("Fast Fiver")
	assert.Nil(err)

}

func TestGameCannotHaveStartedWhenChoosingRoleset(t *testing.T) {
	assert := assert.New(t)
	p1 := player.NewPlayer()
	p2 := player.NewPlayer()
	p3 := player.NewPlayer()
	p4 := player.NewPlayer()
	p5 := player.NewPlayer()
	g := NewGame(p1)
	assert.Equal(Setup, g.State())
	g.AddPlayer(p2)
	g.AddPlayer(p3)
	g.AddPlayer(p4)
	g.AddPlayer(p5)
	err := g.ChooseRoleset("Vanilla Fiver")
	assert.Nil(err)
	err = g.StartGame()
	assert.Nil(err)

	err = g.ChooseRoleset("Fast Fiver")

	assert.Error(err)
}

func TestCannotStartWhenNotInSetup(t *testing.T) {
	assert := assert.New(t)
	p1 := player.NewPlayer()
	p2 := player.NewPlayer()
	p3 := player.NewPlayer()
	p4 := player.NewPlayer()
	p5 := player.NewPlayer()
	g := NewGame(p1)
	assert.Equal(Setup, g.State())
	g.AddPlayer(p2)
	g.AddPlayer(p3)
	g.AddPlayer(p4)
	g.AddPlayer(p5)
	err := g.ChooseRoleset("Vanilla Fiver")
	assert.Nil(err)
	err = g.StartGame()
	assert.Nil(err)

	err = g.StartGame()

	assert.Error(err)
}

func TestCannotStartGameWithoutRoleset(t *testing.T) {
	l := player.NewPlayer()
	g := NewGame(l)

	err := g.StartGame()

	assert.Error(t, err)
}

func TestCannotStartGameWithMismatchedPlayerCountAndRoleset(t *testing.T) {
	assert := assert.New(t)
	l := player.NewPlayer()
	g := NewGame(l)
	err := g.ChooseRoleset("Vanilla Fiver")
	assert.Nil(err)

	err = g.StartGame()

	assert.Error(err) // TODO match on the error?
}

func TestRolesAreAssignedAtGameStart(t *testing.T) {
	assert := assert.New(t)
	p1 := player.NewPlayer()
	p2 := player.NewPlayer()
	p3 := player.NewPlayer()
	p4 := player.NewPlayer()
	p5 := player.NewPlayer()

	g := NewGame(p1)
	assert.Equal(Setup, g.State())

	g.AddPlayer(p2)
	g.AddPlayer(p3)
	g.AddPlayer(p4)
	g.AddPlayer(p5)

	err := g.ChooseRoleset("Vanilla Fiver")
	assert.Nil(err)

	assert.Nil(p1.Role)
	assert.Nil(p2.Role)
	assert.Nil(p3.Role)
	assert.Nil(p4.Role)
	assert.Nil(p5.Role)

	err = g.StartGame()
	assert.Nil(err)
	assert.Equal(Running, g.State())

	assignedRoles := make([]*role.Role, 0)
	for _, p := range g.Players {
		assert.NotNil(p.Role)
		assignedRoles = append(assignedRoles, p.Role)
	}
	assert.ElementsMatch(g.Roleset.Roles, assignedRoles)
}

// really more of an integration test, maybe!
func TestFiver(t *testing.T) {
	assert := assert.New(t)
	p1 := player.NewPlayer()
	p2 := player.NewPlayer()
	p3 := player.NewPlayer()
	p4 := player.NewPlayer()
	p5 := player.NewPlayer()

	g := NewGame(p1)
	assert.Equal(Setup, g.State())

	g.AddPlayer(p2)
	g.AddPlayer(p3)
	g.AddPlayer(p4)
	g.AddPlayer(p5)

	err := g.ChooseRoleset("Vanilla Fiver")
	assert.Nil(err)

	err = g.StartGame()
	assert.Nil(err)
	assert.Equal(Running, g.State())
}
