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
	err = g.Start()
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
	err = g.Start()
	assert.Nil(err)

	err = g.Start()

	assert.Error(err)
}

func TestCannotStartGameWithoutRoleset(t *testing.T) {
	l := player.NewPlayer()
	g := NewGame(l)

	err := g.Start()

	assert.Error(t, err)
}

func TestCannotStartGameWithMismatchedPlayerCountAndRoleset(t *testing.T) {
	assert := assert.New(t)
	l := player.NewPlayer()
	g := NewGame(l)
	err := g.ChooseRoleset("Vanilla Fiver")
	assert.Nil(err)

	err = g.Start()

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

	err = g.Start()
	assert.Nil(err)
	assert.Equal(Running, g.State())

	assignedRoles := []*role.Role{}
	for _, p := range g.Players {
		assert.NotNil(p.Role)
		assignedRoles = append(assignedRoles, p.Role)
	}
	assert.ElementsMatch(g.Roleset.Roles, assignedRoles)
}

// really more of an integration test, maybe!
func TestGame(t *testing.T) {
	assert := assert.New(t)

	g := NewGame(player.NewPlayer())
	assert.Equal(Setup, g.State())

	t.Run("Choose roleset and signup", func(t *testing.T) {

		for i := 0; i < 10; i++ {
			g.AddPlayer(player.NewPlayer())
		}

		err := g.ChooseRoleset("Imperfect Eleven")
		assert.Nil(err)

		err = g.Start()
		assert.Nil(err)
		assert.Equal(Running, g.State())
	})

	var wolf1, wolf2, sorcerer, hunter, seer, v1, v2, v3, v4, v5, v6 *player.Player

	// assign each role to a known variable
	for _, p := range g.Players {
		if p.Role.IsMaxEvil() {
			if wolf1 == nil {
				wolf1 = p
			} else {
				wolf2 = p
			}
		} else if p.Role.IsAuxEvil() {
			sorcerer = p
		} else if p.Role.Parity == 2 {
			hunter = p
		} else if p.Role.IsSeer() {
			seer = p
		} else {
			if v1 == nil {
				v1 = p
			} else if v2 == nil {
				v2 = p
			} else if v3 == nil {
				v3 = p
			} else if v4 == nil {
				v4 = p
			} else if v5 == nil {
				v5 = p
			} else {
				v6 = p
			}
		}
	}

	t.Run("N0", func(t *testing.T) {
		// wolf 1 knows wolf 2
		assert.Len(wolf1.Views, 1)
		assert.Contains(wolf1.Views, &player.View{
			Player:    wolf2,
			Attribute: role.MaxEvilAttribute,
			Hit:       true,
			GamePhase: 0,
		})

		// wolf 2 knows wolf 1
		assert.Len(wolf2.Views, 1)
		assert.Contains(wolf2.Views, &player.View{
			Player:    wolf1,
			Attribute: role.MaxEvilAttribute,
			Hit:       true,
			GamePhase: 0,
		})

		// sorc has a random clear, does NOT know the wolves
		assert.Len(sorcerer.Views, 1)
		sorcClear := sorcerer.Views[0]
		assert.Equal(role.SeerAttribute, sorcClear.Attribute)
		assert.False(sorcClear.Hit)
		assert.NotEqual(sorcClear.Player, seer)
		assert.Equal(0, sorcClear.GamePhase)

		// hunter knows nothing
		assert.Empty(hunter.Views)

		// seer has a random clear
		assert.Len(seer.Views, 1)
		seerClear := seer.Views[0]
		assert.Equal(role.MaxEvilAttribute, seerClear.Attribute)
		assert.False(seerClear.Hit)
		assert.NotEqual(seerClear.Player, wolf1)
		assert.NotEqual(seerClear.Player, wolf2)
		assert.Equal(0, seerClear.GamePhase)

		// villagers know nothing
		assert.Empty(v1.Views)
		assert.Empty(v2.Views)
		assert.Empty(v3.Views)
		assert.Empty(v4.Views)
		assert.Empty(v5.Views)
		assert.Empty(v6.Views)

		assert.Equal(0, g.Phase)
		assert.Equal(Running, g.State())
		assert.True(g.IsDay())
		assert.False(g.IsNight())
	})

	t.Run("D1, villager dies", func(t *testing.T) {
		assert.Equal(2, len(g.AliveMaxEvils()))
		assert.True(g.IsDay())
		assert.Equal(0, g.Phase)

		assert.Nil(g.Vote(wolf1, v1)) // 1/6 needed
		assert.Nil(g.Vote(wolf2, v2))
		assert.Nil(g.Vote(sorcerer, v1)) // 2/6
		assert.Nil(g.Vote(hunter, seer))
		assert.Nil(g.Vote(seer, v3))
		assert.Nil(g.Vote(v1, wolf1))
		assert.Nil(g.Vote(v2, v1)) // 3/6
		assert.Nil(g.Vote(v3, hunter))
		assert.Nil(g.Vote(v4, v5))
		assert.Nil(g.Vote(v5, v1)) // 4/6
		assert.Nil(g.Vote(v6, v5))

		assert.True(g.IsDay())
		assert.Equal(0, g.Phase)

		assert.Nil(g.Vote(seer, v1)) // 5/6
		assert.True(g.IsDay())
		assert.Equal(0, g.Phase)

		assert.Nil(g.Vote(hunter, v1)) // 6/6
		assert.True(g.IsNight())
		assert.Equal(1, g.Phase)
	})

	t.Run("N1, villager eaten", func(t *testing.T) {
		// can't vote at night
	})

	t.Run("D2, villager dies", func(t *testing.T) {
	})

	t.Run("N2, double hits, villager dies", func(t *testing.T) {
	})

	t.Run("D3, wolf dies", func(t *testing.T) {
	})

	t.Run("N3, seer dies", func(t *testing.T) {
	})

	t.Run("D4, sorc dies", func(t *testing.T) {
	})

	t.Run("N4, villager dies", func(t *testing.T) {
	})

	t.Run("D5, villager dies, good wins", func(t *testing.T) {
	})
}
