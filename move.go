package main

import (
	"fmt"
)

type position struct {
	y, x int
}

func (p position) String() string {
	return fmt.Sprintf("%d - %d", p.y, p.x)
}

type capturedPieces []position

func (c capturedPieces) Len() int {
	return len(c)
}

func (c capturedPieces) Less(i, j int) bool {
	if c[i].y < c[j].y {
		return true
	}
	return c[i].x < c[j].x
}

func (c capturedPieces) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

type Move struct {
	start, end position
	capturedPieces
}

func (m Move) String() string {
	str := fmt.Sprintf("{%s -> %s", m.start, m.end)
	if m.capturedPieces != nil {
		str += " | "
		for _, p := range m.capturedPieces {
			str += fmt.Sprintf("%s ", p)
		}
	}
	str += "}"
	return str
}
