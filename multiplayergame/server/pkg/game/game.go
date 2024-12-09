package game

import (
	"sync"
)

// Game is used to hold references to all Players Registered, and Broadcasting etc
type Game struct {
	Pool    *Pool
	Bullets sync.Map
}

// NewGame is used to initalize all the values inside the Game
func NewGame() *Game {
	g := &Game{
		Pool: NewPool(),
	}

	return g
}

func (g *Game) Start() {
	go g.Pool.Start()
}
