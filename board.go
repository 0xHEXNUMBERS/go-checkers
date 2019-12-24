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

func (b *Board) comboCheck(i, j, si, sj int, verticalMoves []int, player byte) []Move {
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
		//Moveing left
		if inBounds(i+vert, j-horiz) {
			if b.isOppositePlayer(i+vert, j-horiz, player) {
				if inBounds(i+vert+vert, j-1) && b[i+vert+vert][j-1] == '_' {
					//Don't double count caputred piece
					tmpMoving := b[i][j]
					b[i][j] = '_'
					b[i+vert+vert][j-1] = tmpMoving
					tmpCapture := b[i+vert][j-horiz]
					b[i+vert][j-horiz] = '_'
					combos := b.comboCheck(i+vert+vert, j-1, si, sj, verticalMoves, player)
					//Restore State
					b[i+vert][j-horiz] = tmpCapture
					b[i][j] = tmpMoving
					b[i+vert+vert][j-1] = '_'

					if combos == nil {
						move := Move{
							start: position{si, sj},
							end:   position{i + vert + vert, j - 1},
						}
						move.addCapturedPiece(
							position{i + vert, j - horiz},
						)
						moves = append(moves, move)
					} else {
						for index, _ := range combos {
							combos[index].addCapturedPiece(
								position{i + vert, j - horiz},
							)
						}
						moves = append(moves, combos...)
					}
				}
			}
		}
		//Moving right
		if inBounds(i+vert, j+(1-horiz)) {
			if b.isOppositePlayer(i+vert, j+(1-horiz), player) {
				if inBounds(i+vert+vert, j+1) && b[i+vert+vert][j+1] == '_' {
					//Don't double count caputred and moving pieces
					tmpMoving := b[i][j]
					b[i][j] = '_'
					b[i+vert+vert][j+1] = tmpMoving
					tmpCapture := b[i+vert][j+(1-horiz)]
					b[i+vert][j+(1-horiz)] = '_'
					combos := b.comboCheck(i+vert+vert, j+1, si, sj, verticalMoves, player)
					//Restore State
					b[i+vert][j+(1-horiz)] = tmpCapture
					b[i][j] = tmpMoving
					b[i+vert+vert][j+1] = '_'

					if combos == nil {
						move := Move{
							start: position{si, sj},
							end:   position{i + vert + vert, j + 1},
						}
						move.addCapturedPiece(
							position{i + vert, j + (1 - horiz)},
						)
						moves = append(moves, move)
					} else {
						for index, _ := range combos {
							combos[index].addCapturedPiece(
								position{i + vert, j + (1 - horiz)},
							)
						}
						moves = append(moves, combos...)
					}
				}
			}
		}
	}
	return moves
}

func (b Board) checkDirections(i, j int, verticalMoves []int, player byte) []Move {
	var moves []Move = nil

	horiz := 1
	if rowParity(i) {
		horiz = 0
	}

	for _, vert := range verticalMoves {
		//Moving left
		if inBounds(i+vert, j-horiz) {
			if b[i+vert][j-horiz] == '_' {
				//If we're jumping to an empty spot,
				//Add the possible move
				move := Move{
					start: position{i, j},
					end:   position{i + vert, j - horiz},
				}
				moves = append(moves, move)
			}
		}
		//Moving right
		if inBounds(i+vert, j+(1-horiz)) {
			if b[i+vert][j+(1-horiz)] == '_' {
				//If we're jumping to an empty spot,
				//Add the possible move
				move := Move{
					start: position{i, j},
					end:   position{i + vert, j + (1 - horiz)},
				}
				moves = append(moves, move)
			}
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

	moves := b.comboCheck(i, j, i, j, verticalMoves, b[i][j])
	moves = append(moves, b.checkDirections(i, j, verticalMoves, b[i][j])...)
	return moves
}
