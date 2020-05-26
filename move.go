package checkers

import (
	"fmt"
	"sort"
	"strings"
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
	start, end     position
	capturedPieces string
}

//String returns a string representation of the move
//
//The string returned is of the form "s -> e {p1|p2|...|pn}" such that...
//
//s is the starting position of the piece being moved
//
//e is the ending position of the piece being moved,
//
//pn is a piece captured by the move
func (m Move) String() string {
	return fmt.Sprintf("{%s -> %s {%s}}", m.start, m.end, m.capturedPieces)
}

func (m Move) inBounds() bool {
	return m.start.inBounds() && m.end.inBounds()
}

func (m *Move) addCapturedPiece(p position) {
	if len(m.capturedPieces) == 0 {
		m.capturedPieces = p.String()
	} else {
		pieces := strings.Split(m.capturedPieces, "|")
		pieces = append(pieces, p.String())
		sort.Strings(pieces)
		m.capturedPieces = strings.Join(pieces, "|")
	}
}

func (m Move) getCapturedPieces() (pieces []position) {
	serializedPieces := strings.Split(m.capturedPieces, "|")

	if len(serializedPieces) == 1 && serializedPieces[0] == "" {
		return
	}

	pieces = make([]position, len(serializedPieces))
	for i, s := range serializedPieces {
		var y, x int
		fmt.Sscanf(s, "%d-%d", &y, &x)
		pieces[i] = position{y, x}
	}
	return
}
