package Game

import (
	"fmt"
	"log"
	. "projet-ai/Types"
	"projet-ai/utils"
)

type Game struct {
	Holes       *Holes
	Players     [2]*Player
	PlayerIndex int
}

func NewGame(player1 *Player, player2 *Player) *Game {
	game := &Game{
		Holes:   NewHoles(2, 2, 1),
		Players: [2]*Player{player1, player2}, PlayerIndex: 0}

	return game
}

var ch = make(chan Move)
var ReceivedMove = make(chan Move)
var Step = 1

const MqttToggle = false

func (g *Game) Run() {

	var move Move

	g.Holes.Print()
	if MqttToggle {
		go InitMqtt(ch)
	}
	for g.Holes.NbSeeds > 0 {
		println("Step : ", Step)
		opponentIndex := 0
		if g.PlayerIndex == 0 {
			opponentIndex = 1
		} else {
			opponentIndex = 0
		}

		if g.Players[g.PlayerIndex].IsAi {
			ai := NewIa(g.Players[g.PlayerIndex].PlayerIndex)
			move = ai.ChooseMove(g.Holes, g.Players[g.PlayerIndex], g.Players[opponentIndex])
			log.Println(counter)
			counter = 0

			if MqttToggle {
				ch <- move
			}
		} else {

			if MqttToggle {
				move = <-ReceivedMove
			} else {
				human := NewHuman(g.Players[g.PlayerIndex].PlayerIndex)
				move = human.ChooseMove(g.Holes)
			}
		}

		lastSeededHole := Sow(g.Holes, &move)

		seedsCaptured := Capture(g.Holes, lastSeededHole)
		g.Players[g.PlayerIndex].Seeds += seedsCaptured

		fmt.Println("Seeds : ", g.Holes.NbSeeds)
		fmt.Println("Player 1 seeds : ", g.Players[0].Seeds)
		fmt.Println("Player 2 seeds : ", g.Players[1].Seeds)

		isMaximising := g.PlayerIndex == 0
		switch CheckEndGame(g.Holes, g.Players[0], g.Players[1], isMaximising) {
		case utils.Resume:
			break
		case utils.P1Wins:
			fmt.Println("Player 1 wins")
			return
		case utils.P2Wins:
			fmt.Println("Player 2 wins")
			return
		case utils.Draw:
			fmt.Println("Draw")
			return
		default:
			break
		}
		g.Holes.Print()

		if g.PlayerIndex == 1 {
			g.PlayerIndex = 0
		} else {
			g.PlayerIndex = 1
		}
		Step++
	}

}
func Sow(holes *Holes, move *Move) int {
	var color Color
	if move.MoveColor != utils.TR {
		color = Color(move.MoveColor)
	} else {
		color = utils.Transparent
	}

	seeds := holes.Pop(move.HoleIndex, color)
	index := move.HoleIndex
	initialIndex := move.HoleIndex

	index = (index + 1) % 16
	holes.Add(index, 1, color)
	seeds--

	for seeds > 0 {
		if move.MoveColor == MoveColor(utils.B) || move.MoveColor == MoveColor(utils.TB) {
			index = (index + 2) % 16
		} else {
			index = (index + 1) % 16
		}
		if index == initialIndex {
			index = (index + 1) % 16
		}

		holes.Add(index, 1, color)
		seeds--

	}
	return index
}

func Capture(holes *Holes, holeIndex int) int {
	sum := 0
	seeds := holes.Sum(holeIndex)
	for seeds == 2 || seeds == 3 {
		sum += seeds
		holes.Empty(holeIndex)

		if holeIndex == 0 {
			holeIndex = 15
		} else {
			holeIndex--
		}
		seeds = holes.Sum(holeIndex)
	}
	return sum
}
func CheckEndGame(holes *Holes, player1 *Player, player2 *Player, isMaximising bool) GameResult {
	if player1.Seeds >= 41 {
		return utils.P1Wins
	}
	if player2.Seeds >= 41 {
		return utils.P2Wins
	}
	if isMaximising {
		if IsPlayerStarved(holes, player2.PlayerIndex) {
			player1.Seeds += holes.NbSeeds
			holes.NbSeeds = 0
			return utils.P1Wins
		}
	} else {
		if IsPlayerStarved(holes, player1.PlayerIndex) {
			player2.Seeds += holes.NbSeeds
			holes.NbSeeds = 0

			return utils.P2Wins
		}
	}

	if holes.NbSeeds < 10 {
		if player1.Seeds > player2.Seeds {
			return utils.P1Wins
		} else if player1.Seeds < player2.Seeds {
			return utils.P2Wins
		} else {
			return utils.Draw
		}

	}
	return utils.Resume
}

func IsPlayerStarved(holes *Holes, playerIndex int) bool {
	for i := playerIndex - 1; i < 16; i += 2 {
		if holes.Sum(i) > 0 {
			return false
		}
	}
	return true
}
