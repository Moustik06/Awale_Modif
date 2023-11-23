package Game

import (
	"fmt"
	. "projet-ai/Types"
	"projet-ai/utils"
	"sync"
)

type IHoles interface {
	Print()
	Pop(index int, color Color) int
	Sum(index int) int
	Empty(index int)
	IsEmpty(index int, color Color) bool
}

var holesPool sync.Pool

func init() {
	holesPool.New = func() interface{} {
		return &Holes{}
	}
}

type Holes struct {
	Holes   [16][3]int
	NbSeeds int
}

func (h *Holes) Print() {
	fmt.Println()
	for i := 0; i < 8; i++ {
		fmt.Printf("%2d(%s%dB%s|%s%dR%s|%s%dT%s)=%d\t ", i+1, utils.AnsiColorBlue, h.Holes[i][0], utils.AnsiColorReset, utils.AnsiColorRed, h.Holes[i][1], utils.AnsiColorReset, utils.AnsiColorTransparent, h.Holes[i][2], utils.AnsiColorReset, h.Sum(i))
	}
	fmt.Println()
	for i := 15; i > 7; i-- {
		fmt.Printf("%2d(%s%dB%s|%s%dR%s|%s%dT%s)=%d\t ", i+1, utils.AnsiColorBlue, h.Holes[i][0], utils.AnsiColorReset, utils.AnsiColorRed, h.Holes[i][1], utils.AnsiColorReset, utils.AnsiColorTransparent, h.Holes[i][2], utils.AnsiColorReset, h.Sum(i))
	}
	fmt.Println("\n")
}

func NewHoles(blueSeeds int, redSeeds int, transparentSeeds int) *Holes {
	holes := &Holes{NbSeeds: 80}
	for i := 0; i < 16; i++ {
		holes.Holes[i][utils.Blue] = blueSeeds
		holes.Holes[i][utils.Red] = redSeeds
		holes.Holes[i][utils.Transparent] = transparentSeeds
	}
	return holes
}

func (h *Holes) Pop(index int, color Color) int {
	seeds := h.Holes[index][color]
	h.Holes[index][color] = 0
	h.NbSeeds -= seeds

	return seeds
}
func (h *Holes) Add(index int, seeds int, color Color) {
	h.Holes[index][color] += seeds
	h.NbSeeds += seeds
}
func (h *Holes) Sum(index int) int {
	sum := 0
	for _, i := range h.Holes[index] {
		sum += i
	}
	return sum
}

func (h *Holes) Empty(index int) {

	for i := 0; i < 3; i++ {
		h.Pop(index, Color(i))
	}
}

func (h *Holes) IsEmpty(index int, color Color) bool {
	if h.Holes[index][color] == 0 {
		return true
	}
	return false
}

func (h *Holes) Copy() *Holes {
	holes := holesPool.Get().(*Holes)
	*holes = *h
	return holes
}

func PutHolesCopy(h *Holes) {
	holesPool.Put(h)
}
