package Game

import (
	. "projet-ai/Types"
	"sync"
)

var playerPool sync.Pool

func init() {
	playerPool.New = func() interface{} {
		return &Player{}
	}
}

type Move struct {
	HoleIndex int
	MoveColor MoveColor
}

type Player struct {
	Seeds       int
	PlayerIndex int
	IsAi        bool
}

func (p *Player) Copy() *Player {
	playerCopy := playerPool.Get().(*Player)
	*playerCopy = *p
	return playerCopy
}

func PutPlayer(p *Player) {
	playerPool.Put(p)
}
