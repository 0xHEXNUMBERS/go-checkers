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

//Board is a representation of the board state
type Board [ROWS][COLS]byte

func (b Board) String() string {
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

func (b Board) isOppositePlayer(p position, player byte) bool {
	i, j := p.i, p.j
	if player == 'x' || player == 'X' {
		return b[i][j] == 'o' || b[i][j] == 'O'
	}
	return b[i][j] == 'x' || b[i][j] == 'X'
}

func (b *Board) capturePiece(
	start, from, piece, to position,
	verticalMoves []int, player byte) []Move {

	fromI, fromJ := from.i, from.j
	pieceI, pieceJ := piece.i, piece.j
	toI, toJ := to.i, to.j
	var moves []Move = nil

	//Simulate capturing the piece
	tmpMoving := b[fromI][fromJ]
	b[fromI][fromJ] = '_'
	b[toI][toJ] = tmpMoving
	tmpCapture := b[pieceI][pieceJ]
	b[pieceI][pieceJ] = '_'

	//Check if we can continue capturing from here
	combos := b.captureCheck(
		start, to,
		verticalMoves, player,
	)

	//Restore State
	b[pieceI][pieceJ] = tmpCapture
	b[fromI][fromJ] = tmpMoving
	b[toI][toJ] = '_'

	if combos == nil {
		move := Move{
			start: start,
			end:   to,
		}
		move.addCapturedPiece(piece)
		moves = []Move{move}
	} else {
		for index := range combos {
			combos[index].addCapturedPiece(piece)
		}
		moves = combos
	}

	return moves
}

func (b Board) isVacant(p position) bool {
	if !inBounds(p.i, p.j) {
		return false
	}

	return b[p.i][p.j] == '_'
}

func (b *Board) captureCheck(start, to position,
	verticalMoves []int, player byte) []Move {
	var moves []Move = nil

	if to != start {
		moves = []Move{
			Move{
				start: start,
				end:   to,
			},
		}
	}

	i1, j1 := to.i, to.j

	horiz := 1
	if rowParity(i1) {
		horiz = 0
	}

	for _, vert := range verticalMoves {
		//Moving left
		left := position{i1 + vert, j1 - horiz}
		if left.inBounds() {
			if b.isOppositePlayer(left, player) {
				posAfterMove := position{i1 + vert + vert, j1 - 1}
				if posAfterMove.inBounds() && b.isVacant(posAfterMove) {
					combos := b.capturePiece(
						start, to,
						left, posAfterMove,
						verticalMoves, player,
					)
					moves = append(moves, combos...)
				}
			}
		}

		//Moving right
		right := position{i1 + vert, j1 + (1 - horiz)}
		if right.inBounds() {
			if b.isOppositePlayer(right, player) {
				posAfterMove := position{i1 + vert + vert, j1 + 1}
				if posAfterMove.inBounds() && b.isVacant(posAfterMove) {
					combos := b.capturePiece(
						start, to,
						right, posAfterMove,
						verticalMoves, player,
					)
					moves = append(moves, combos...)
				}
			}
		}
	}
	return moves
}

func (b Board) checkForAdjacentVacantSpots(p position, verticalMoves []int) []Move {
	i, j := p.i, p.j
	var moves []Move = nil

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
			moves = append(moves, move)
		}

		//Moving right
		right := position{i + vert, j + (1 - horiz)}
		if b.isVacant(right) {
			move := Move{
				start: p,
				end:   right,
			}
			moves = append(moves, move)
		}
	}
	return moves
}

func (b Board) getMovesFromPos(p position) []Move {
	i, j := p.i, p.j
	verticalMoves := make([]int, 0)

	if b[i][j] == 'X' || b[i][j] == 'O' {
		verticalMoves = append(verticalMoves, 1, -1)
	} else if b[i][j] == 'x' {
		verticalMoves = append(verticalMoves, 1)
	} else if b[i][j] == 'o' {
		verticalMoves = append(verticalMoves, -1)
	}

	moves := b.captureCheck(
		position{i, j}, position{i, j},
		verticalMoves, b[i][j],
	)
	moves = append(
		moves,
		b.checkForAdjacentVacantSpots(
			position{i, j}, verticalMoves)...,
	)
	return moves
}
