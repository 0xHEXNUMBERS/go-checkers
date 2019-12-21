package main

import (
	"fmt"
)

type position struct {
	y, x int
}

func (p position) String() string {
	return fmt.Sprintf("%d-%d", p.y, p.x)
}

type capturedPieces string

func (c *capturedPieces) addPiece(p position) {
	if len(*c) == 0 {
		*c = (capturedPieces)(p.String())
	} else {
		*c += (capturedPieces)("|" + p.String())
	}
}

type Move struct {
	start, end position
	capturedPieces
}

func (m Move) String() string {
	return fmt.Sprintf("{%s -> %s {%s}}", m.start, m.end, m.capturedPieces)
}
