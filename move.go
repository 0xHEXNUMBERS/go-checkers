package main

import (
	"fmt"
	"sort"
	"strings"
)

type position struct {
	y, x int
}

func (p position) String() string {
	return fmt.Sprintf("%d-%d", p.y, p.x)
}

type Move struct {
	start, end     position
	capturedPieces string
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

func (m Move) String() string {
	return fmt.Sprintf("{%s -> %s {%s}}", m.start, m.end, m.capturedPieces)
}

func (m Move) getCapturedPieces() (pieces []position) {
	serializedPieces := strings.Split(m.capturedPieces, "|")

	if len(serializedPieces) == 1 && serializedPieces[0] == "" {
		return
	}

	for _, s := range serializedPieces {
		var y, x int
		fmt.Sscanf(s, "%d-%d", &y, &x)
		pieces = append(pieces, position{y, x})
	}
	return
}
