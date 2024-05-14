package gamechannel

type ActivityType int

const (
	Join ActivityType = iota
	Vote
)

type Activity struct {
	Type ActivityType
}
