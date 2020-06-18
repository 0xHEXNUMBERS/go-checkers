package checkers

import (
	"fmt"
)

type position struct {
	i, j int
}

func (p position) String() string {
	return fmt.Sprintf("%d-%d", p.i, p.j)
}

func (p position) inBounds() bool {
	return inBounds(p.i, p.j)
}

//Move represents an action in checkers
type Move struct {
	start, end    position
	capturedPiece position
}

//String returns a string representation of the move
//
//The string returned is of the form "s -> e {p}" such that...
//
//s is the starting position of the piece being moved
//
//e is the ending position of the piece being moved,
//
//p is the piece captured by the move
func (m Move) String() string {
	return fmt.Sprintf("{%s -> %s {%s}}", m.start, m.end, m.capturedPiece)
}

func (m Move) inBounds() bool {
	return m.start.inBounds() && m.end.inBounds() && m.capturedPiece.inBounds()
}
