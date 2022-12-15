package common

import "fmt"

type Direction int

const (
	UP Direction = iota + 1
	RIGHT
	DOWN
	LEFT
	UNKNOWN
	DEADEND
)

type Position struct {
	X, Y int
}

func (p Position) String() string {
	return fmt.Sprintf("(%v,%v)", p.X, p.Y)
}

func (d Direction) String() string {
	switch d {
	case UP:
		return "UP"
	case DOWN:
		return "DOWN"
	case RIGHT:
		return "RIGHT"
	case LEFT:
		return "LEFT"
	default:
		return "N/A"
	}
}

func (d Direction) StringForMap() string {
	switch d {
	case UP:
		return "^"
	case DOWN:
		return "v"
	case RIGHT:
		return ">"
	case LEFT:
		return "<"
	default:
		return "."
	}
}
