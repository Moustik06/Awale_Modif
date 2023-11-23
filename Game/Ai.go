package Game

import (
	"fmt"
	"projet-ai/utils"
	"sort"
	"sync"
)

type Ia struct {
	Player
}

type IaPlayer interface {
	ChooseMove(holes *Holes, opponent *Player) Move
}

var movesPool sync.Pool
var counter = 0

func init() {
	movesPool.New = func() interface{} {
		return make([]Move, 0, 24)
	}
}

func (ia *Ia) SortMovesByEvaluation(possiblesMoves []Move, holes *Holes, current *Player, opponent *Player) {
	moveScores := make(map[Move]float64)

	// Calculer et stocker les scores pour chaque coup
	for _, move := range possiblesMoves {
		moveScores[move] = ia.Evaluate(holes, current, opponent, false)
	}

	// Trier les coups possibles en fonction de leur score
	sort.Slice(possiblesMoves, func(i, j int) bool {
		scoreA := moveScores[possiblesMoves[i]]
		scoreB := moveScores[possiblesMoves[j]]
		return scoreA > scoreB // Triez les coups par ordre décroissant de score
	})
}

func (ia *Ia) ChooseMove(holes *Holes, current *Player, opponent *Player) Move {
	bestMove := Move{}
	bestValue := -1000.00
	depth := 6

	if Step > 60 {
		depth = 7
	}
	if Step > 90 {
		depth = 8
	}

	if Step > 150 {
		depth = 9
	}
	possiblesMoves := movesPool.Get().([]Move)
	defer movesPool.Put(possiblesMoves)
	setPossiblesMoves(holes, ia.PlayerIndex, &possiblesMoves)
	ia.SortMovesByEvaluation(possiblesMoves, holes, current, opponent)

	type Result struct {
		move  Move
		value float64
	}

	resultCh := make(chan Result)
	wg := sync.WaitGroup{}
	wg.Add(len(possiblesMoves))

	for _, move := range possiblesMoves {
		go func(move Move) {
			defer wg.Done()

			newHoles := holes.Copy()
			newCurrent := current.Copy()
			newOpponent := opponent.Copy()

			lastSeededHole := Sow(newHoles, &move)
			newCurrent.Seeds += Capture(newHoles, lastSeededHole)

			gameResult := CheckEndGame(holes, newCurrent, newOpponent, false)
			if newCurrent.PlayerIndex == 1 {
				switch gameResult {
				case utils.Resume:
					// println("Resume")
					break
				case utils.P1Wins:
					// println("P1Wins")
					resultCh <- Result{move, 1000}
					break
				case utils.P2Wins:
					// println("P2Wins")
					bestValue = -1000
					break
				case utils.Draw:
					bestValue = 0
					break
				}
			} else if newCurrent.PlayerIndex == 2 {
				switch gameResult {
				case utils.Resume:
					// println("Resume")
					break
				case utils.P1Wins:
					// println("P1Wins")
					bestValue = -1000
					break
				case utils.P2Wins:
					// println("P2Wins")
					resultCh <- Result{move, 1000}
					break
				case utils.Draw:
					bestValue = 0
					break
				}
			}
			var score float64
			score = ia.MiniMax(newHoles, newCurrent, newOpponent, depth-1, false, -1000, 1000, newCurrent)

			PutPlayer(newCurrent)
			PutPlayer(newOpponent)
			PutHolesCopy(newHoles)

			resultCh <- Result{move, score}
		}(move)
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	for result := range resultCh {
		if result.value > bestValue {
			bestValue = result.value
			bestMove = result.move
		}
		if bestValue == -1000.00 && bestMove.HoleIndex+1 == 1 && bestMove.MoveColor == utils.B {
			bestMove = result.move
			bestValue = result.value
		}
	}

	ia.printMoves(&bestMove)
	println("Score du Move : ", bestValue)
	return bestMove
}

func NewIa(playerIndex int) Ia {
	return Ia{Player{0, playerIndex, true}}
}

func setPossiblesMoves(holes *Holes, playerIndex int, possiblesMoves *[]Move) {
	for i := playerIndex - 1; i < 16; i += 2 {
		if !holes.IsEmpty(i, utils.Red) {
			*possiblesMoves = append(*possiblesMoves, Move{i, utils.R})
		}
		if !holes.IsEmpty(i, utils.Blue) {
			*possiblesMoves = append(*possiblesMoves, Move{i, utils.B})
		}
		if !holes.IsEmpty(i, utils.Transparent) {
			*possiblesMoves = append(*possiblesMoves, Move{i, utils.TR})
			*possiblesMoves = append(*possiblesMoves, Move{i, utils.TB})
		}
	}
}

func (ia *Ia) printMoves(move *Move) {
	color := ""
	switch move.MoveColor {
	case utils.R:
		color = "RED"
	case utils.B:
		color = "BLUE"
	case utils.TR:
		color = "TRANSPARENT RED"
	case utils.TB:
		color = "TRANSPARENT BLUE"
	}
	fmt.Println("Hole : ", move.HoleIndex+1, " Color : ", color)
}

func (ia *Ia) MiniMax(holes *Holes, current *Player, opponent *Player, depth int, isMaximizing bool, alpha float64, beta float64, joueurOrigin *Player) float64 {
	counter += 1
	gameResult := CheckEndGame(holes, current, opponent, !isMaximizing)
	if joueurOrigin.PlayerIndex == 1 {
		switch gameResult {
		case utils.Resume:
			// println("Resume")
			break
		case utils.P1Wins:
			// println("P1Wins")
			return 1000
		case utils.P2Wins:
			// println("P2Wins")
			return -1000
		case utils.Draw:
			return 0
		}
	} else if joueurOrigin.PlayerIndex == 2 {
		switch gameResult {
		case utils.Resume:
			// println("Resume")
			break
		case utils.P1Wins:
			// println("P1Wins")
			return -1000
		case utils.P2Wins:
			// println("P2Wins")
			return 1000
		case utils.Draw:
			return 0
		}
	}

	if depth == 0 {
		return ia.Evaluate(holes, current, opponent, isMaximizing)
	}

	possiblesMoves := movesPool.Get().([]Move)

	defer movesPool.Put(possiblesMoves)

	var bestValue float64
	var playerIndex int
	var mustBreak bool

	if isMaximizing {
		bestValue = -1000
		playerIndex = current.PlayerIndex
	} else {
		bestValue = 1000
		playerIndex = opponent.PlayerIndex
	}

	setPossiblesMoves(holes, playerIndex, &possiblesMoves)
	ia.SortMovesByEvaluation(possiblesMoves, holes, current, opponent)

	for _, move := range possiblesMoves {
		newHoles := holes.Copy()
		newCurrent := current.Copy()
		newOpponent := opponent.Copy()

		lastSeededHole := Sow(newHoles, &move)
		if isMaximizing {
			newCurrent.Seeds += Capture(newHoles, lastSeededHole)
		} else {
			newOpponent.Seeds += Capture(newHoles, lastSeededHole)
		}

		score := ia.MiniMax(newHoles, newCurrent, newOpponent, depth-1, !isMaximizing, alpha, beta, joueurOrigin)

		if isMaximizing {
			bestValue = utils.Max(score, bestValue)
			alpha = utils.Max(alpha, bestValue)
			mustBreak = beta <= alpha
		} else {
			bestValue = utils.Min(score, bestValue)
			beta = utils.Min(beta, bestValue)
			mustBreak = beta <= alpha
		}

		if mustBreak {
			break
		}
		PutPlayer(newCurrent)
		PutPlayer(newOpponent)
		PutHolesCopy(newHoles)

	}
	return bestValue
}

const (
	WEIGHT_SCORE               = 9.52
	WEIGHT_BLOCK               = 1.08
	WEIGHT_STARVATION          = 901.50
	WEIGHT_EMPTY_OPPONET_HOLES = 8.65
	WEIGHT_EMPTY_CURRENT_HOLES = 6.61
	WEIGHT_CAPTURE             = 2.04
)

func (ia *Ia) Evaluate(holes *Holes, current *Player, opponent *Player, isMaximizing bool) float64 {
	// Base score : différence entre les graines des joueurs
	score := 0.00
	score = float64(current.Seeds-opponent.Seeds) * WEIGHT_SCORE

	if IsPlayerStarved(holes, opponent.PlayerIndex) {
		score += WEIGHT_STARVATION
	}

	if IsPlayerStarved(holes, current.PlayerIndex) {
		score -= WEIGHT_STARVATION
	}

	// Bonus pour chaque trou de l'adversaire qui est vide
	emptyHolesOpponent := 0.00
	for i := opponent.PlayerIndex - 1; i < 16; i += 2 {
		if holes.Sum(i) == 0 {
			emptyHolesOpponent += WEIGHT_BLOCK
		}
	}
	score += emptyHolesOpponent * WEIGHT_EMPTY_OPPONET_HOLES

	// // Malus pour chaque trou du joueur courant qui est vide
	emptyHolesCurrent := 0.00
	for i := current.PlayerIndex - 1; i < 16; i += 2 {
		if holes.Sum(i) == 0 {
			emptyHolesCurrent++
		}
	}
	score -= emptyHolesCurrent * WEIGHT_EMPTY_CURRENT_HOLES

	// Bonus pour chaque trou de l'adversaire qui a 2 ou 3 graines (potentiel de capture au prochain tour)
	capturableHolesOpponent := 0.00
	for i := opponent.PlayerIndex - 1; i < 16; i += 2 {
		seeds := holes.Sum(i)
		if seeds == 2 || seeds == 3 {
			capturableHolesOpponent++
		}
	}
	score += capturableHolesOpponent * WEIGHT_CAPTURE

	return score
}
