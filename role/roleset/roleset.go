package roleset

import (
	"github.com/awoo-detat/werewolf/role"
)

type Roleset struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Roles       []*role.Role `json:"roles"`
}

func (rs *Roleset) String() string {
	return rs.Name
}

func List() map[string]*Roleset {
	return sets
}

type RolesetMap map[string]*Roleset

var sets = RolesetMap{}

func registerRoleset(roleset *Roleset) {
	sets[roleset.Name] = roleset
}
