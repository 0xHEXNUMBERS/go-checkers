package checkers

import (
	"fmt"
)

const (
	ROWS = 8
	COLS = 4
)

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

func (b Board) isOppositePlayer(i, j int, player byte) bool {
	if player == 'x' || player == 'X' {
		return b[i][j] == 'o' || b[i][j] == 'O'
	} else {
		return b[i][j] == 'x' || b[i][j] == 'X'
	}
}

func (b *Board) capturePiece(
	startI, startJ, fromI, fromJ, pieceI, pieceJ, toI, toJ int,
	verticalMoves []int, player byte) []Move {
	var moves []Move = nil

	//Simulate capturing the piece
	tmpMoving := b[fromI][fromJ]
	b[fromI][fromJ] = '_'
	b[toI][toJ] = tmpMoving
	tmpCapture := b[pieceI][pieceJ]
	b[pieceI][pieceJ] = '_'

	//Check if we can continue capturing from here
	combos := b.captureCheck(toI, toJ, startI, startJ, verticalMoves, player)

	//Restore State
	b[pieceI][pieceJ] = tmpCapture
	b[fromI][fromJ] = tmpMoving
	b[toI][toJ] = '_'

	if combos == nil {
		move := Move{
			start: position{startI, startJ},
			end:   position{toI, toJ},
		}
		move.addCapturedPiece(
			position{pieceI, pieceJ},
		)
		moves = append(moves, move)
	} else {
		for index, _ := range combos {
			combos[index].addCapturedPiece(
				position{pieceI, pieceJ},
			)
		}
		moves = append(moves, combos...)
	}

	return moves
}

func (b *Board) captureCheck(i, j, si, sj int, verticalMoves []int, player byte) []Move {
	var moves []Move = nil

	if i != si || j != sj {
		moves = []Move{
			Move{
				start: position{si, sj},
				end:   position{i, j},
			},
		}
	}

	horiz := 1
	if rowParity(i) {
		horiz = 0
	}

	for _, vert := range verticalMoves {
		//Moving left
		if inBounds(i+vert, j-horiz) {
			if b.isOppositePlayer(i+vert, j-horiz, player) {
				if inBounds(i+vert+vert, j-1) && b[i+vert+vert][j-1] == '_' {
					combos := b.capturePiece(
						si, sj, i, j,
						i+vert, j-horiz,
						i+vert+vert, j-1,
						verticalMoves, player,
					)
					moves = append(moves, combos...)
				}
			}
		}
		//Moving right
		if inBounds(i+vert, j+(1-horiz)) {
			if b.isOppositePlayer(i+vert, j+(1-horiz), player) {
				if inBounds(i+vert+vert, j+1) && b[i+vert+vert][j+1] == '_' {
					combos := b.capturePiece(
						si, sj, i, j,
						i+vert, j+(1-horiz),
						i+vert+vert, j+1,
						verticalMoves, player,
					)
					moves = append(moves, combos...)
				}
			}
		}
	}
	return moves
}

func (b Board) isSpotVacant(i, j int) bool {
	if !inBounds(i, j) {
		return false
	}

	return b[i][j] == '_'
}

func (b Board) checkForAdjacentVacantpots(i, j int, verticalMoves []int) []Move {
	var moves []Move = nil

	horiz := 1
	if rowParity(i) {
		horiz = 0
	}

	for _, vert := range verticalMoves {
		//Moving left
		if b.isSpotVacant(i+vert, j-horiz) {
			move := Move{
				start: position{i, j},
				end:   position{i + vert, j - horiz},
			}
			moves = append(moves, move)
		}

		//Moving right
		if b.isSpotVacant(i+vert, j+(1-horiz)) {
			move := Move{
				start: position{i, j},
				end:   position{i + vert, j + (1 - horiz)},
			}
			moves = append(moves, move)
		}
	}
	return moves
}

func (b Board) getMovesFromPos(i, j int) []Move {
	verticalMoves := make([]int, 0)
	if b[i][j] == 'X' || b[i][j] == 'O' {
		verticalMoves = append(verticalMoves, 1, -1)
	} else if b[i][j] == 'x' {
		verticalMoves = append(verticalMoves, 1)
	} else if b[i][j] == 'o' {
		verticalMoves = append(verticalMoves, -1)
	}

	moves := b.captureCheck(i, j, i, j, verticalMoves, b[i][j])
	moves = append(moves, b.checkForAdjacentVacantpots(i, j, verticalMoves)...)
	return moves
}
