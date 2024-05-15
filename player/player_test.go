package player

import (
	"testing"

	"github.com/awoo-detat/werewolf/role"

	"github.com/stretchr/testify/assert"
)

func TestInitialization(t *testing.T) {
	assert := assert.New(t)
	p := NewPlayer(NewMockCommunicator())

	assert.NotEmpty(p.ID)
	assert.NotEmpty(p.Name)
	assert.Empty(p.Views)
	assert.Nil(p.Role)
	assert.Equal(p.Name, p.String())
}

func TestSetName(t *testing.T) {
	p := NewPlayer(NewMockCommunicator())
	name := "Test Player"

	p.SetName(name)

	assert.Equal(t, name, p.Name)
	assert.Equal(t, name, p.String())
}

func TestSetRole(t *testing.T) {
	p := NewPlayer(NewMockCommunicator())
	r := role.Villager()

	p.SetRole(r)

	assert.Equal(t, p.Role, r)
}
