package gamechannel

// GameChannel provides a way for the players to communicate to the game. It is a one-way communication.
type GameChannel chan *Activity
