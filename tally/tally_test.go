package tally

import (
	"testing"

	"github.com/awoo-detat/werewolf/player"

	"github.com/stretchr/testify/assert"
)

func TestTally(t *testing.T) {
	assert := assert.New(t)

	dake := player.NewPlayer()
	dake.SetName("Dake")
	tommy := player.NewPlayer()
	tommy.SetName("Tommy")
	sigafoos := player.NewPlayer()
	sigafoos.SetName("Sigafoos")
	players := []*player.Player{dake, tommy, sigafoos}
	gt := New(players)

	t.Run("Creation", func(t *testing.T) {
		t.Run("Sorted by name", func(t *testing.T) {
			t.Parallel()
			assert.Len(gt.List, len(players))
			assert.Same(dake, gt.List[0].Player)
			assert.Same(sigafoos, gt.List[1].Player)
			assert.Same(tommy, gt.List[2].Player)
		})

		t.Run("List has no votes", func(t *testing.T) {
			t.Parallel()
			for _, ti := range gt.List {
				assert.NotNil(ti.Player)
				assert.Empty(ti.Votes)
			}
		})

		t.Run("Inverted has no votes", func(t *testing.T) {
			t.Parallel()
			assert.Len(gt.Inverted, len(players))
			for _, p := range players {
				assert.Nil(gt.Inverted[p])
			}
		})
	})

	t.Run("Vote adds to the tally", func(t *testing.T) {
		gt.Vote(sigafoos, tommy)

		leader := gt.List[0]
		assert.Equal(tommy, leader.Player, "Tommy has the most votes")
		assert.Len(leader.Votes, 1)
		assert.Equal(sigafoos, leader.Votes[0].Voter)
		assert.Equal(tommy, gt.Inverted[sigafoos].Candidate, "inverted tally has the vote")
	})

	t.Run("LHLV", func(t *testing.T) {
		gt.Vote(tommy, dake)

		leader := gt.List[0]
		assert.Equal(tommy, leader.Player, "Tommy is winning by LHLV")
		assert.Len(leader.Votes, 1)
		assert.Equal(sigafoos, leader.Votes[0].Voter)
		assert.Equal(tommy, gt.Inverted[sigafoos].Candidate, "inverted tally has the vote")

		second := gt.List[1]
		assert.Equal(dake, second.Player, "Dake has a vote")
		assert.Len(second.Votes, 1)
		assert.Equal(tommy, second.Votes[0].Voter)
		assert.Equal(dake, gt.Inverted[tommy].Candidate, "inverted tally has the vote")
	})

	t.Run("changing vote updates both tallies", func(t *testing.T) {
		gt.Vote(sigafoos, dake)

		leader := gt.List[0]
		assert.Equal(dake, leader.Player, "Dake now has two votes")
		assert.Len(leader.Votes, 2)
		assert.Equal(tommy, leader.Votes[0].Voter)
		assert.Equal(sigafoos, leader.Votes[1].Voter)

		assert.Equal(dake, gt.Inverted[sigafoos].Candidate, "inverted tally has been updated")

		assert.Empty(gt.List[1].Votes)
		assert.Empty(gt.List[2].Votes)
	})

	t.Run("Can unvote", func(t *testing.T) {
		gt.Unvote(sigafoos)

		leader := gt.List[0]
		assert.Equal(dake, leader.Player, "Dake still has one vote")
		assert.Len(leader.Votes, 1)
		assert.Equal(tommy, leader.Votes[0].Voter)

		assert.Nil(gt.Inverted[sigafoos], "inverted tally has been updated")

		assert.Empty(gt.List[1].Votes)
		assert.Empty(gt.List[2].Votes)
	})
}
