package main

import (
	"fmt"
	_ "net/http/pprof"
	. "projet-ai/Game"

	"github.com/bytedance/gopkg/util/gctuner"
)

func main() {
	gcTuning()
	player1 := NewIa(1)
	player2 := NewIa(2)
	game := NewGame(&player1.Player, &player2.Player)

	game.Run()

	game.Holes.Print()
	fmt.Println("Seeds : ", game.Holes.NbSeeds)
	fmt.Println("Player 1 seeds : ", player1.Seeds)
	fmt.Println("Player 2 seeds : ", player2.Seeds)
}

func gcTuning() {
	var limit float64 = 10 * 1024 * 1024 * 1024
	threshold := uint64(limit * 0.7)
	gctuner.Tuning(threshold)
}
