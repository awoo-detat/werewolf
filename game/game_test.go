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
	l := player.NewPlayer(player.NewMockCommunicator())

	g := NewGame(l)

	assert.NotEmpty(g.ID)
	assert.Equal(l, g.Leader)
	assert.Len(g.Players, 1)
	assert.Nil(g.Roleset)
	assert.Equal(Setup, g.State())
	assert.Equal(0, g.Phase)
}

func TestAddingPlayers(t *testing.T) {
	l := player.NewPlayer(player.NewMockCommunicator())
	g := NewGame(l)
	p := player.NewPlayer(player.NewMockCommunicator())

	g.AddPlayer(p)

	assert.Len(t, g.Players, 2)
}

func TestSetRoleset(t *testing.T) {
	assert := assert.New(t)
	l := player.NewPlayer(player.NewMockCommunicator())
	g := NewGame(l)
	rs := roleset.VanillaFiver()

	err := g.ChooseRoleset("Vanilla Fiver")

	assert.Equal(rs, g.Roleset)
	assert.Nil(err)
}

func TestRolesetMustExist(t *testing.T) {
	assert := assert.New(t)
	l := player.NewPlayer(player.NewMockCommunicator())
	g := NewGame(l)

	err := g.ChooseRoleset("dkjjfkfwegfwegy")

	assert.Nil(g.Roleset)
	assert.Error(err)
}

func TestCanChangeRolesetWhenInSetup(t *testing.T) {
	assert := assert.New(t)
	l := player.NewPlayer(player.NewMockCommunicator())
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
	p1 := player.NewPlayer(player.NewMockCommunicator())
	p2 := player.NewPlayer(player.NewMockCommunicator())
	p3 := player.NewPlayer(player.NewMockCommunicator())
	p4 := player.NewPlayer(player.NewMockCommunicator())
	p5 := player.NewPlayer(player.NewMockCommunicator())
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
	p1 := player.NewPlayer(player.NewMockCommunicator())
	p2 := player.NewPlayer(player.NewMockCommunicator())
	p3 := player.NewPlayer(player.NewMockCommunicator())
	p4 := player.NewPlayer(player.NewMockCommunicator())
	p5 := player.NewPlayer(player.NewMockCommunicator())
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
	l := player.NewPlayer(player.NewMockCommunicator())
	g := NewGame(l)

	err := g.Start()

	assert.Error(t, err)
}

func TestCannotStartGameWithMismatchedPlayerCountAndRoleset(t *testing.T) {
	assert := assert.New(t)
	l := player.NewPlayer(player.NewMockCommunicator())
	g := NewGame(l)
	err := g.ChooseRoleset("Vanilla Fiver")
	assert.Nil(err)

	err = g.Start()

	assert.Error(err) // TODO match on the error?
}

func TestRolesAreAssignedAtGameStart(t *testing.T) {
	assert := assert.New(t)
	p1 := player.NewPlayer(player.NewMockCommunicator())
	p2 := player.NewPlayer(player.NewMockCommunicator())
	p3 := player.NewPlayer(player.NewMockCommunicator())
	p4 := player.NewPlayer(player.NewMockCommunicator())
	p5 := player.NewPlayer(player.NewMockCommunicator())

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

	g := NewGame(player.NewPlayer(player.NewMockCommunicator()))
	assert.Equal(Setup, g.State())

	t.Run("Choose roleset and signup", func(t *testing.T) {

		for i := 0; i < 10; i++ {
			g.AddPlayer(player.NewPlayer(player.NewMockCommunicator()))
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
				wolf1.SetName("Wolf 1")
			} else {
				wolf2 = p
				wolf2.SetName("Wolf 2")
			}
		} else if p.Role.IsAuxEvil() {
			sorcerer = p
			sorcerer.SetName("Sorcerer")
		} else if p.Role.Parity == 2 {
			hunter = p
			hunter.SetName("Hunter")
		} else if p.Role.IsSeer() {
			seer = p
			seer.SetName("Seer")
		} else {
			if v1 == nil {
				v1 = p
				v1.SetName("Villager 1")
			} else if v2 == nil {
				v2 = p
				v2.SetName("Villager 2")
			} else if v3 == nil {
				v3 = p
				v3.SetName("Villager 3")
			} else if v4 == nil {
				v4 = p
				v4.SetName("Villager 4")
			} else if v5 == nil {
				v5 = p
				v5.SetName("Villager 5")
			} else {
				v6 = p
				v6.SetName("Villager 6")
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

		// now it's day!
		assert.Equal(1, g.Phase)
		assert.Equal(Running, g.State())
		assert.True(g.IsDay())
		assert.False(g.IsNight())
	})

	t.Run("D1, villager dies", func(t *testing.T) {
		assert.Equal(2, len(g.AliveMaxEvils()))
		assert.True(g.IsDay())
		assert.Equal(1, g.Phase)
		// seer can do a view already!
		assert.Nil(g.SetNightAction(&player.FingerPoint{From: seer, To: v1}))

		assert.Nil(g.Vote(&player.FingerPoint{From: wolf1, To: v1})) // 1/6 needed
		assert.Nil(g.Vote(&player.FingerPoint{From: wolf2, To: v2}))
		assert.Nil(g.Vote(&player.FingerPoint{From: sorcerer, To: v1})) // 2/6
		assert.Nil(g.Vote(&player.FingerPoint{From: hunter, To: seer}))
		assert.Nil(g.Vote(&player.FingerPoint{From: seer, To: v3}))
		assert.Nil(g.Vote(&player.FingerPoint{From: v1, To: wolf1}))
		assert.Nil(g.Vote(&player.FingerPoint{From: v2, To: v1})) // 3/6
		assert.Nil(g.Vote(&player.FingerPoint{From: v3, To: hunter}))
		assert.Nil(g.Vote(&player.FingerPoint{From: v4, To: v5}))
		assert.Nil(g.Vote(&player.FingerPoint{From: v5, To: v1})) // 4/6
		assert.Nil(g.Vote(&player.FingerPoint{From: v6, To: v5}))

		assert.True(g.IsDay())
		assert.Equal(1, g.Phase)

		assert.Nil(g.Vote(&player.FingerPoint{From: seer, To: v1})) // 5/6
		assert.True(g.IsDay())
		assert.Equal(1, g.Phase)

		assert.Nil(g.Vote(&player.FingerPoint{From: hunter, To: v1})) // 6/6
		assert.True(g.IsNight())
		assert.Equal(2, g.Phase)
		assert.False(v1.Role.Alive)

		for _, p := range g.Players {
			v := p.Views[len(p.Views)-1]
			assert.Equal(v1, v.Player)
			assert.Equal(v1.Role, v.Role)
			assert.True(v.Hit)
			assert.Equal(g.Phase-1, v.GamePhase)
		}
	})

	t.Run("N1, villager eaten", func(t *testing.T) {
		// can't vote at night
		assert.Error(g.Vote(&player.FingerPoint{From: wolf1, To: v2}))

		assert.Nil(g.SetNightAction(&player.FingerPoint{From: wolf1, To: v2}))
		assert.Nil(g.SetNightAction(&player.FingerPoint{From: wolf2, To: v3})) // this one "takes"
		// seer changing their view since v1 is super dead
		assert.Nil(g.SetNightAction(&player.FingerPoint{From: seer, To: v2}))
		assert.Error(g.SetNightAction(&player.FingerPoint{From: v1, To: wolf1}))
		assert.Error(g.SetNightAction(&player.FingerPoint{From: sorcerer, To: v1}))
		// sure why not have fun
		assert.Nil(g.SetNightAction(&player.FingerPoint{From: v3, To: v2}))
		assert.True(g.IsNight())
		assert.Equal(2, g.Phase)

		assert.Nil(g.SetNightAction(&player.FingerPoint{From: sorcerer, To: hunter}))
		assert.True(g.IsDay())
		assert.Equal(3, g.Phase)

		assert.False(v3.Role.Alive)
		for _, p := range g.Players {
			v := p.Views[len(p.Views)-1]
			assert.Equal(v3, v.Player)
			assert.Equal(v3.Role, v.Role)
			assert.True(v.Hit)
			assert.Equal(g.Phase-1, v.GamePhase)
		}

		// no new views
		assert.Len(wolf1.Views, 3)
		assert.Len(wolf2.Views, 3)
		assert.Len(hunter.Views, 2)
		assert.Len(v1.Views, 2)
		assert.Len(v2.Views, 2)
		assert.Len(v3.Views, 2)
		assert.Len(v4.Views, 2)
		assert.Len(v5.Views, 2)
		assert.Len(v6.Views, 2)

		// sorc has a new view
		assert.Len(sorcerer.Views, 4)
		sorcView := sorcerer.Views[2]
		assert.Equal(role.SeerAttribute, sorcView.Attribute)
		assert.Equal(hunter, sorcView.Player)
		assert.False(sorcView.Hit)
		assert.Equal(2, sorcView.GamePhase)

		// seer has a new view
		assert.Len(seer.Views, 4)
		seerView := seer.Views[2]
		assert.Equal(role.MaxEvilAttribute, seerView.Attribute)
		assert.Equal(v2, seerView.Player)
		assert.False(seerView.Hit)
		assert.Equal(2, seerView.GamePhase)

	})

	t.Run("D2, villager dies", func(t *testing.T) {
		// new day!
		assert.Empty(g.Tally.List[0].Votes)
		assert.Empty(g.nightActions)
		assert.Equal(2, len(g.AliveMaxEvils()))
		assert.True(g.IsDay())
		assert.Equal(3, g.Phase)
		assert.Equal(9, len(g.AlivePlayers))

		assert.Error(g.Vote(&player.FingerPoint{From: wolf1, To: v1}))
		assert.Error(g.Vote(&player.FingerPoint{From: v1, To: wolf1}))
		assert.Nil(g.Vote(&player.FingerPoint{From: wolf1, To: v2})) // 1/5 needed
		assert.Nil(g.Vote(&player.FingerPoint{From: wolf2, To: seer}))
		assert.Nil(g.Vote(&player.FingerPoint{From: sorcerer, To: v2})) // 2/5
		assert.Nil(g.Vote(&player.FingerPoint{From: hunter, To: seer}))
		assert.Nil(g.Vote(&player.FingerPoint{From: seer, To: v4}))
		assert.Nil(g.Vote(&player.FingerPoint{From: v2, To: hunter}))
		assert.Nil(g.Vote(&player.FingerPoint{From: v4, To: v5}))
		assert.Nil(g.Vote(&player.FingerPoint{From: v5, To: v2})) // 3/5
		assert.Nil(g.Vote(&player.FingerPoint{From: v6, To: v5}))
		assert.Nil(g.Vote(&player.FingerPoint{From: seer, To: v2}))   // 4/5
		assert.Nil(g.Vote(&player.FingerPoint{From: hunter, To: v2})) // 5/5

		assert.True(g.IsNight())
		assert.Equal(4, g.Phase)
		assert.False(v2.Role.Alive)
		for _, p := range g.Players {
			v := p.Views[len(p.Views)-1]
			assert.Equal(v2, v.Player)
			assert.Equal(v2.Role, v.Role)
			assert.True(v.Hit)
			assert.Equal(g.Phase-1, v.GamePhase)
		}
	})

	t.Run("N2, double hits, villager dies", func(t *testing.T) {
		assert.Nil(g.SetNightAction(&player.FingerPoint{From: wolf1, To: v4}))
		assert.Nil(g.SetNightAction(&player.FingerPoint{From: wolf2, To: v4}))
		assert.Nil(g.SetNightAction(&player.FingerPoint{From: seer, To: wolf1}))
		assert.Nil(g.SetNightAction(&player.FingerPoint{From: sorcerer, To: seer}))

		assert.True(g.IsDay())
		assert.Equal(5, g.Phase)

		assert.False(v4.Role.Alive)
		for _, p := range g.Players {
			v := p.Views[len(p.Views)-1]
			assert.Equal(v4, v.Player)
			assert.Equal(v4.Role, v.Role)
			assert.True(v.Hit)
			assert.Equal(g.Phase-1, v.GamePhase)
		}

		// no new views
		assert.Len(wolf1.Views, 5)
		assert.Len(wolf2.Views, 5)
		assert.Len(hunter.Views, 4)
		assert.Len(v1.Views, 4)
		assert.Len(v2.Views, 4)
		assert.Len(v3.Views, 4)
		assert.Len(v4.Views, 4)
		assert.Len(v5.Views, 4)
		assert.Len(v6.Views, 4)

		// sorc has a new view
		assert.Len(sorcerer.Views, 7)
		sorcView := sorcerer.Views[5]
		assert.Equal(role.SeerAttribute, sorcView.Attribute)
		assert.Equal(seer, sorcView.Player)
		assert.True(sorcView.Hit)
		assert.Equal(4, sorcView.GamePhase)

		// seer has a new view
		assert.Len(seer.Views, 7)
		seerView := seer.Views[5]
		assert.Equal(role.MaxEvilAttribute, seerView.Attribute)
		assert.Equal(wolf1, seerView.Player)
		assert.True(seerView.Hit)
		assert.Equal(4, seerView.GamePhase)
	})

	t.Run("D3, wolf dies", func(t *testing.T) {
		assert.Empty(g.Tally.List[0].Votes)
		assert.Empty(g.nightActions)
		assert.Equal(2, len(g.AliveMaxEvils()))
		assert.True(g.IsDay())
		assert.Equal(5, g.Phase)
		assert.Equal(7, len(g.AlivePlayers))

		assert.Nil(g.Vote(&player.FingerPoint{From: wolf1, To: seer}))
		assert.Nil(g.Vote(&player.FingerPoint{From: wolf2, To: wolf1})) // 1/4
		assert.Nil(g.Vote(&player.FingerPoint{From: sorcerer, To: hunter}))
		assert.Nil(g.Vote(&player.FingerPoint{From: hunter, To: wolf1})) // 2/4
		assert.Nil(g.Vote(&player.FingerPoint{From: seer, To: wolf1}))   // 3/4

		assert.True(g.IsDay())
		assert.Equal(5, g.Phase)
		assert.Nil(g.Vote(&player.FingerPoint{From: v5, To: wolf1}))
		// didn't even need v6

		assert.True(g.IsNight())
		assert.Equal(6, g.Phase)
		assert.False(wolf1.Role.Alive)
		for _, p := range g.Players {
			v := p.Views[len(p.Views)-1]
			assert.Equal(wolf1, v.Player)
			assert.Equal(wolf1.Role, v.Role)
			assert.True(v.Hit)
			assert.Equal(g.Phase-1, v.GamePhase)
		}
	})

	t.Run("N3, seer dies", func(t *testing.T) {
		assert.Nil(g.SetNightAction(&player.FingerPoint{From: wolf2, To: seer}))
		assert.Nil(g.SetNightAction(&player.FingerPoint{From: seer, To: sorcerer}))
		assert.True(g.IsNight())
		assert.Nil(g.SetNightAction(&player.FingerPoint{From: sorcerer, To: seer}))

		assert.True(g.IsDay())
		assert.Equal(7, g.Phase)

		assert.False(seer.Role.Alive)
		for _, p := range g.Players {
			v := p.Views[len(p.Views)-1]
			assert.Equal(seer, v.Player)
			assert.Equal(seer.Role, v.Role)
			assert.True(v.Hit)
			assert.Equal(g.Phase-1, v.GamePhase)
		}

		// no new views
		assert.Len(wolf1.Views, 7)
		assert.Len(wolf2.Views, 7)
		assert.Len(hunter.Views, 6)
		assert.Len(v1.Views, 6)
		assert.Len(v2.Views, 6)
		assert.Len(v3.Views, 6)
		assert.Len(v4.Views, 6)
		assert.Len(v5.Views, 6)
		assert.Len(v6.Views, 6)

		// sorc has a new view
		assert.Len(sorcerer.Views, 10)
		sorcView := sorcerer.Views[8]
		assert.Equal(role.SeerAttribute, sorcView.Attribute)
		assert.Equal(seer, sorcView.Player)
		assert.True(sorcView.Hit)
		assert.Equal(6, sorcView.GamePhase)

		// seer has a new view, even though they're dead
		assert.Len(seer.Views, 10)
		seerView := seer.Views[8]
		assert.Equal(role.MaxEvilAttribute, seerView.Attribute)
		assert.Equal(sorcerer, seerView.Player)
		assert.False(seerView.Hit)
		assert.Equal(6, seerView.GamePhase)

	})

	t.Run("D4, sorc dies", func(t *testing.T) {
		assert.Empty(g.Tally.List[0].Votes)
		assert.Empty(g.nightActions)
		assert.Equal(1, len(g.AliveMaxEvils()))
		assert.True(g.IsDay())
		assert.Equal(7, g.Phase)
		assert.Equal(5, len(g.AlivePlayers))

		assert.Nil(g.Vote(&player.FingerPoint{From: wolf2, To: hunter}))
		assert.Nil(g.Vote(&player.FingerPoint{From: sorcerer, To: hunter}))
		assert.Nil(g.Vote(&player.FingerPoint{From: hunter, To: sorcerer})) // 1/3
		assert.Nil(g.Vote(&player.FingerPoint{From: v5, To: sorcerer}))     // 2/3
		assert.True(g.IsDay())
		assert.Equal(7, g.Phase)
		assert.Nil(g.Vote(&player.FingerPoint{From: v6, To: sorcerer})) // 3/3

		assert.True(g.IsNight())
		assert.Equal(8, g.Phase)
		assert.False(sorcerer.Role.Alive)
		for _, p := range g.Players {
			v := p.Views[len(p.Views)-1]
			assert.Equal(sorcerer, v.Player)
			assert.Equal(sorcerer.Role, v.Role)
			assert.True(v.Hit)
			assert.Equal(g.Phase-1, v.GamePhase)
		}
	})

	t.Run("N4, villager dies", func(t *testing.T) {
		assert.Nil(g.SetNightAction(&player.FingerPoint{From: wolf2, To: v5}))

		assert.True(g.IsDay())
		assert.Equal(9, g.Phase)

		assert.False(v5.Role.Alive)
		for _, p := range g.Players {
			v := p.Views[len(p.Views)-1]
			assert.Equal(v5, v.Player)
			assert.Equal(v5.Role, v.Role)
			assert.True(v.Hit)
			assert.Equal(g.Phase-1, v.GamePhase)
		}

		// no new views
		assert.Len(wolf1.Views, 9)
		assert.Len(wolf2.Views, 9)
		assert.Len(hunter.Views, 8)
		assert.Len(v1.Views, 8)
		assert.Len(v2.Views, 8)
		assert.Len(v3.Views, 8)
		assert.Len(v4.Views, 8)
		assert.Len(v5.Views, 8)
		assert.Len(v6.Views, 8)
		assert.Len(sorcerer.Views, 12)
		assert.Len(seer.Views, 12)

	})

	t.Run("D5, villager dies, good wins by hunter victory", func(t *testing.T) {
		assert.Empty(g.Tally.List[0].Votes)
		assert.Empty(g.nightActions)
		assert.Equal(1, len(g.AliveMaxEvils()))
		assert.True(g.IsDay())
		assert.Equal(9, g.Phase)
		assert.Equal(3, len(g.AlivePlayers))

		assert.Nil(g.Vote(&player.FingerPoint{From: wolf2, To: v6})) // 1/2 needed
		assert.Nil(g.Vote(&player.FingerPoint{From: v6, To: hunter}))
		assert.True(g.IsDay())
		assert.Equal(9, g.Phase)
		assert.Nil(g.Vote(&player.FingerPoint{From: hunter, To: v6})) // 2/2

		assert.False(v6.Role.Alive)
		for _, p := range g.Players {
			v := p.Views[len(p.Views)-1]
			assert.Equal(v6, v.Player)
			assert.Equal(v6.Role, v.Role)
			assert.True(v.Hit)
			assert.Equal(g.Phase, v.GamePhase)
		}
		assert.True(wolf2.Role.Alive)
		assert.Len(g.AliveMaxEvils(), 1, "didn't need to kill all wolves")
		assert.Equal(Finished, g.State())
		assert.Equal(role.Good, g.Winner)
	})
}
