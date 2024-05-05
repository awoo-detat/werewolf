package role

type PlayerType int

const (
	Good    PlayerType = iota
	Evil               = iota
	Neutral            = iota
)

type Attribute int

const (
	MaxEvilAttribute Attribute = 1 << iota
	AuxEvilAttribute
	SeerAttribute
	TinkerAttribute
)

func (a Attribute) String() string {
	switch a {
	case MaxEvilAttribute:
		return "Max Evil"
	case AuxEvilAttribute:
		return "Aux Evil"
	case SeerAttribute:
		return "Seer"
	case TinkerAttribute:
		return "Tinker"
	}
	return ""
}

type Action int

const (
	viewForMax Action = 1 << iota
	nightKill
	viewForSeer
	viewForAux
	randomN0Clear
	knowsMaxes
)

type Role struct {
	Name           string     `json:"name"`
	Description    string     `json:"description"`
	Team           PlayerType `json:"team"`
	Parity         int        `json:"-"`
	VoteMultiplier int        `json:"-"`
	Health         int        `json:"-"`
	Alive          bool       `json:"alive"`
	Actions        Action     `json:"night_action"`
	Attributes     Attribute  `json:"-"`
}

func (r *Role) String() string {
	return r.Name
}

// IsMaxEvil returns whether or not a player is a max evil (ie a Werewolf). Another player viewing for max evil should use ViewForMaxEvil().
func (r *Role) IsMaxEvil() bool {
	return r.Attributes&MaxEvilAttribute > 0
}

// ViewForMaxEvil allows Seers to view if a role is max evil. It differs from IsMaxEvil because the Tinker can invert the result.
func (r *Role) ViewForMaxEvil() bool {
	if r.Attributes&TinkerAttribute > 0 {
		return !r.IsMaxEvil()
	}
	return r.IsMaxEvil()
}

// IsAuxEvil returns whether or not a player is aux evil (ie a Cultist).
func (r *Role) IsAuxEvil() bool {
	return r.Attributes&AuxEvilAttribute > 0
}

// ViewForAuxEvil allows Seers to view if a role is aux evil. It differs from IsAuxEvil because the tinker
// can invert the result.
func (r *Role) ViewForAuxEvil() bool {
	if r.Attributes&TinkerAttribute > 0 {
		return !r.IsAuxEvil()
	}
	return r.IsAuxEvil()
}

// IsSeer returns whether or not a player is a seer.
func (r *Role) IsSeer() bool {
	return r.Attributes&SeerAttribute > 0
}

// ViewForSeer allows sorcerers to view if a role is a seer. It differs from IsSeer because the
// tinker can invert the result.
func (r *Role) ViewForSeer() bool {
	if r.Attributes&TinkerAttribute > 0 {
		return !r.IsSeer()
	}
	return r.IsSeer()
}

func (r *Role) CanViewForMax() bool {
	return r.Actions&viewForMax > 0
}

func (r *Role) CanNightKill() bool {
	return r.Actions&nightKill > 0
}

func (r *Role) CanViewForSeer() bool {
	return r.Actions&viewForSeer > 0
}

func (r *Role) CanViewForAux() bool {
	return r.Actions&viewForAux > 0
}

func (r *Role) HasRandomN0Clear() bool {
	return r.Actions&randomN0Clear > 0
}

func (r *Role) KnowsMaxes() bool {
	return r.Actions&knowsMaxes > 0
}

// SetTinker makes a role a Tinker: all views will be the inverse of the truth
func (r *Role) SetTinker() {
	r.Attributes = r.Attributes | TinkerAttribute
}

// Kill attempts to kill the player. If they had more than 1 health (ie
// were "tough") then they will remain alive.
// It returns whether or not the kill was successful.
func (r *Role) Kill() bool {
	// maybe this should error if you try to kill a dead person?
	r.Health--
	if r.Health <= 0 {
		r.Alive = false
	}
	return !r.Alive
}
