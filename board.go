package checkers

import (
	"fmt"
)

const (
	//ROWS is the number of rows in a checkers board
	ROWS = 8

	//COLS is the number of cols in a checkers board
	//this variable is represented as half the amount of the columns
	//on a typical checkers board as half of the slots on the board are
	//unused. This implementation takes use of that to save on memory
	COLS = 4
)

type board [ROWS][COLS]byte

func (b board) String() string {
	str := ""
	for i := 0; i < ROWS; i++ {
		str += "["
		for j := 0; j < COLS; j++ {
			if i%2 == 1 {
				str += fmt.Sprintf("%c|_|", b[i][j])
			} else {
				str += fmt.Sprintf("_|%c|", b[i][j])
			}
		}
		str = str[:len(str)-1] + "]\n"
	}
	return str
}

func (b board) isOppositePlayer(p position, player byte) bool {
	if !p.inBounds() {
		return false
	}

	i, j := p.i, p.j
	if player == 'x' || player == 'X' {
		return b[i][j] == 'o' || b[i][j] == 'O'
	}
	return b[i][j] == 'x' || b[i][j] == 'X'
}

func (b board) isVacant(p position) bool {
	if !inBounds(p.i, p.j) {
		return false
	}

	return b[p.i][p.j] == '_'
}

func (b board) canCaptureLeft(p position, vert, horiz int, player byte) bool {
	left := position{p.i + vert, p.j - horiz}
	if b.isOppositePlayer(left, player) {
		posAfterMove := position{p.i + vert + vert, p.j - 1}
		if b.isVacant(posAfterMove) {
			return true
		}
	}
	return false
}

func (b board) canCaptureRight(p position, vert, horiz int, player byte) bool {
	right := position{p.i + vert, p.j + (1 - horiz)}
	if b.isOppositePlayer(right, player) {
		posAfterMove := position{p.i + vert + vert, p.j + 1}
		if b.isVacant(posAfterMove) {
			return true
		}
	}
	return false
}

func (b board) canCaptureFromPos(p position) bool {
	vertMoves := b.getVertMovesFromPos(p)
	player := b[p.i][p.j]
	horiz := 1
	if rowParity(p.i) {
		horiz = 0
	}

	for _, vert := range vertMoves {
		if b.canCaptureLeft(p, vert, horiz, player) {
			return true
		}

		if b.canCaptureRight(p, vert, horiz, player) {
			return true
		}
	}
	return false
}

func (b board) captureCheck(p position, vertMoves []int, player byte, moves *[]Move) {
	horiz := 1
	if rowParity(p.i) {
		horiz = 0
	}

	for _, vert := range vertMoves {
		//Moving left
		if b.canCaptureLeft(p, vert, horiz, player) {
			*moves = append(
				*moves,
				Move{
					start:         p,
					end:           position{p.i + vert + vert, p.j - 1},
					capturedPiece: position{p.i + vert, p.j - horiz},
				},
			)
		}

		//Moving right
		if b.canCaptureRight(p, vert, horiz, player) {
			*moves = append(
				*moves,
				Move{
					start:         p,
					end:           position{p.i + vert + vert, p.j + 1},
					capturedPiece: position{p.i + vert, p.j + (1 - horiz)},
				},
			)
		}
	}
}

func (b board) checkForAdjacentVacantSpots(p position, verticalMoves []int, moves *[]Move) {
	i, j := p.i, p.j

	horiz := 1
	if rowParity(i) {
		horiz = 0
	}

	for _, vert := range verticalMoves {
		//Moving left
		left := position{i + vert, j - horiz}
		if b.isVacant(left) {
			move := Move{
				start: p,
				end:   left,
			}
			*moves = append(*moves, move)
		}

		//Moving right
		right := position{i + vert, j + (1 - horiz)}
		if b.isVacant(right) {
			move := Move{
				start: p,
				end:   right,
			}
			*moves = append(*moves, move)
		}
	}
}

func (b board) getVertMovesFromPos(p position) []int {
	i, j := p.i, p.j
	var verticalMoves []int

	if b[i][j] == 'X' || b[i][j] == 'O' {
		verticalMoves = []int{1, -1}
	} else if b[i][j] == 'x' {
		verticalMoves = []int{1}
	} else if b[i][j] == 'o' {
		verticalMoves = []int{-1}
	}
	return verticalMoves
}

func (b board) canMoveFromPos(p position) bool {
	i, j := p.i, p.j

	horiz := 1
	if rowParity(i) {
		horiz = 0
	}

	for _, vert := range b.getVertMovesFromPos(p) {
		left := position{i + vert, j - horiz}
		if b.isVacant(left) {
			return true
		} else if b.isOppositePlayer(left, b[i][j]) {
			posAfterMove := position{i + vert + vert, j - 1}
			if b.isVacant(posAfterMove) {
				return true
			}
		}

		right := position{i + vert, j + (1 - horiz)}
		if b.isVacant(right) {
			return true
		} else if b.isOppositePlayer(right, b[i][j]) {
			posAfterMove := position{i + vert + vert, j + 1}
			if b.isVacant(posAfterMove) {
				return true
			}
		}
	}

	return false
}
