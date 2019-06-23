package main

import (
	"fmt"
)

const (
	ROWS = 8
	COLS = 4
)

type Board [][]byte

type Move struct {
	iStart, jStart int
	iEnd, jEnd     int
}

func (m Move) String() string {
	return fmt.Sprintln(m.iStart, "-", m.jStart, "->", m.iEnd, "-", m.jEnd)
}

func NewGame() Board {
	var board Board
	for i := 0; i < ROWS; i++ {
		board = append(board, make([]byte, COLS))
		for j := 0; j < COLS; j++ {
			board[i][j] = '_'
		}
	}
	for i := 0; i < (ROWS / 2) - 1; i++ {
		for j := 0; j < COLS; j++ {
			board[i][j] = 'x'
		}
	}
	for i := (ROWS / 2) + 1; i < ROWS; i++ {
		for j := 0; j < COLS; j++ {
			board[i][j] = 'o'
		}
	}
	return board
}

func (board Board) String() string {
	str := ""
	for i := 0; i < ROWS; i++ {
		str += "["
		for j := 0; j < COLS; j++ {
			if i % 2 == 1 {
				str += fmt.Sprintf("%c|_|", board[i][j])
			} else {
				str += fmt.Sprintf("_|%c|", board[i][j])
			}
		}
		str = str[:len(str)-1] + "]\n"
	}
	return str
}

func inBounds(i, j int) bool {
	return i >= 0 && i < ROWS && j >= 0 && j < COLS
}

func (board Board) getMovesFromPos(i, j int) []interface{} {
	horiz := 1
	if i % 2 == 1 {
		horiz = -1
	}

	verticalMoves := make([]int, 0)
	if board[i][j] == 'X' || board[i][j] == 'O' {
		verticalMoves = append(verticalMoves, 1, -1)
	} else if board[i][j] == 'x' {
		verticalMoves = append(verticalMoves, 1)
	} else if board[i][j] == 'o' {
		verticalMoves = append(verticalMoves, -1)
	}

	moves := make([]interface{}, 0)
	for _, vert := range verticalMoves {
		if inBounds(i+vert, j) && board[i+vert][j] == '_' {
			moves = append(moves, Move{i, j, i+vert, j})
		}
		if inBounds(i+vert, j+horiz) && board[i+vert][j+horiz] == '_' {
			moves = append(moves, Move{i, j, i+vert, j+horiz})
		}
	}
	return moves
}

func (board Board) GetActions() []interface{} {
	var moves []interface{} = make([]interface{}, 0)
	for i := 0; i < ROWS; i++ {
		for j := 0; j < COLS; j++ {
			if board[i][j] == '_' {
				continue
			}
			moves = append(moves, board.getMovesFromPos(i, j)...)
		}
	}
	return moves
}

func main() {
	game := NewGame()
	fmt.Println(game)
	moves := game.GetActions()
	fmt.Println(moves)
}