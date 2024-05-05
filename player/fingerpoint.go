package player

// A FingerPoint represents one player "pointing" at another in
// a real-life game of Werewolf. It can be used for both votes
// and night actions.
type FingerPoint struct {
	From *Player
	To   *Player
}
