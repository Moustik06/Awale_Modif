package Game

import (
	"bufio"
	"fmt"
	"os"
	. "projet-ai/Types"
	. "projet-ai/utils"
	"strconv"
	"strings"
	"unicode"
)

type Human struct {
	Player
}

func (h *Human) ChooseMove(holes *Holes) Move {
	var index int
	var line string

	fmt.Printf("Player %d:\n", h.PlayerIndex)
	print("Input move index and color: ")
	reader := bufio.NewReader(os.Stdin)

	for {
		line, _ = reader.ReadString('\n')
		line = strings.Trim(line, "\n\r")
		index, line = SplitAfterLastInt(line)
		println("Player choose move: ", index, line)
		if h.IsMoveValid(index-1, line, holes) {
			break
		}
	}

	// convert string into MoveColor
	var color MoveColor
	switch line {
	case "B":
		color = B
	case "R":
		color = R
	case "TB":
		color = TB
	default:
		color = TR
	}
	move := Move{index - 1, color}
	return move

}

func (p *Player) IsMoveValid(holeIndex int, seedColor string, holes *Holes) bool {
	// checks inputs
	if (holeIndex < 0 || holeIndex > 15) || (seedColor != "B" && seedColor != "R" && seedColor != "TB" && seedColor != "TR") {
		return false
	}

	// checks if hole belongs to the player
	if (p.PlayerIndex%2 == 0 && holeIndex%2 == 0) || (p.PlayerIndex%2 != 0 && holeIndex%2 != 0) {
		return false
	}

	var color Color
	if seedColor == "B" {
		color = Blue
	} else if seedColor == "R" {
		color = Red
	} else {
		color = Transparent
	}

	return !holes.IsEmpty(holeIndex, color)
}

func NewHuman(idx int) *Human {
	println("NewHuman")
	h := &Human{Player: Player{0, idx, false}}
	return h
}

func SplitAfterLastInt(s string) (int, string) {
	lastIntIndex := 0
	for i, r := range s {
		if unicode.IsDigit(r) {
			lastIntIndex = i
		}
	}

	index, _ := strconv.Atoi(s[:lastIntIndex+1])
	color := s[lastIntIndex+1:]

	return index, color
}
