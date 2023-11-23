package utils

import . "projet-ai/Types"

const (
	AnsiColorRed         = "\x1b[31m"
	AnsiColorBlue        = "\x1b[34m"
	AnsiColorTransparent = "\x1b[37m"
	AnsiColorReset       = "\x1b[0m"
)

const (
	Blue Color = iota
	Red
	Transparent
)

const (
	B MoveColor = iota
	R
	TB
	TR
)

const (
	Resume GameResult = iota
	P1Wins
	P2Wins
	Draw
)
