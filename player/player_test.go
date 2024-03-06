package player

import (
	"testing"

	"github.com/awoo-detat/werewolf/role"

	"github.com/stretchr/testify/assert"
)

func TestInitialization(t *testing.T) {
	p := NewPlayer()

	assert.NotEmpty(t, p.ID)
	assert.Empty(t, p.Name)
	assert.Nil(t, p.Role)
}

func TestSetName(t *testing.T) {
	p := NewPlayer()
	name := "Test Player"

	p.SetName(name)

	assert.Equal(t, name, p.Name)
}

func TestSetRole(t *testing.T) {
	p := NewPlayer()
	r := role.Villager()

	p.SetRole(r)

	assert.Equal(t, p.Role, r)
}
