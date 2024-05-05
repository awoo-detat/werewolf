package player

import (
	"testing"

	"github.com/awoo-detat/werewolf/role"

	"github.com/stretchr/testify/assert"
)

func TestAttributeView(t *testing.T) {
	p := NewPlayer()
	a := role.MaxEvilAttribute
	v := NewAttributeView(p, a, true, 0)

	assert.Equal(t, a.String(), v.For())
}

func TestRoleView(t *testing.T) {
	p := NewPlayer()
	r := role.Seer()
	v := NewRoleView(p, r, 0)

	assert.Equal(t, r.Name, v.For())
}
