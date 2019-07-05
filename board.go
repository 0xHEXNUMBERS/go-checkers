package main

import (
	"fmt"
	"errors"
)

const (
	ROWS = 8
	COLS = 4
)

var (
	ERR_INTERFACE_FAIL = errors.New("Failed to convert input interface to a Move")
)

type Board [ROWS][COLS]byte

type Game struct {
	Board
	oTurn bool
}

type position struct {
	y, x int
}

func (p position) String() string {
	return fmt.Sprintf("%d - %d", p.y, p.x)
}

type Move struct {
	start, end position
	capturedPieces []position
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

func NewGame() Game {
	var board Board
	for i := 0; i < ROWS; i++ {
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
	return Game{Board: board}
}

func NewGameCapture() Game {
	var board Board
	for i := 0; i < ROWS; i++ {
		for j := 0; j < COLS; j++ {
			board[i][j] = '_'
		}
	}
	board[2][2] = 'X'
	board[3][2] = 'o'
	board[3][3] = 'o'
	return Game{Board: board}
}

func NewGameCombo() Game {
	var board Board
	for i := 0; i < ROWS; i++ {
		for j := 0; j < COLS; j++ {
			board[i][j] = '_'
		}
	}
	board[2][2] = 'X'
	board[3][1] = 'o'
	board[3][2] = 'o'
	board[1][1] = 'o'
	board[1][2] = 'o'
	return Game{Board: board}
}

func NewGameCombo2() Game {
	var board Board
	for i := 0; i < ROWS; i++ {
		for j := 0; j < COLS; j++ {
			board[i][j] = '_'
		}
	}
	board[0][0] = 'x'
	board[1][1] = 'o'
	board[3][2] = 'o'
	board[5][2] = 'o'
	return Game{Board: board}
}

func NewGameCombo3() Game {
	var board Board
	for i := 0; i < ROWS; i++ {
		for j := 0; j < COLS; j++ {
			board[i][j] = '_'
		}
	}
	board[0][3] = 'x'
	board[1][3] = 'o'
	board[3][2] = 'o'
	board[5][2] = 'o'
	return Game{Board: board}
}

func NewGameUpgrade() Game {
	var board Board
	for i := 0; i < ROWS; i++ {
		for j := 0; j < COLS; j++ {
			board[i][j] = '_'
		}
	}
	board[1][0] = 'o'
	return Game{board, true}
}

func rowParity(i int) bool {
	return i % 2 == 0
}

func (b Board) String() string {
	str := ""
	for i := 0; i < ROWS; i++ {
		str += "["
		for j := 0; j < COLS; j++ {
			if i % 2 == 1 {
				str += fmt.Sprintf("%c|_|", b[i][j])
			} else {
				str += fmt.Sprintf("_|%c|", b[i][j])
			}
		}
		str = str[:len(str)-1] + "]\n"
	}
	return str
}

func inBounds(i, j int) bool {
	return i >= 0 && i < ROWS && j >= 0 && j < COLS
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
							end: position{i+vert+vert, j-1},
							capturedPieces: make([]position, 1),
						}
						move.capturedPieces[0] = position{i+vert, j-horiz}
						moves = append(moves, move)
					} else {
						for index, _ := range combos {
							combos[index].capturedPieces = append(
								combos[index].capturedPieces,
								position{i+vert, j-horiz},
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
							end: position{i+vert+vert, j+1},
							capturedPieces: make([]position, 1),
						}
						move.capturedPieces[0] = position{i+vert, j+(1-horiz)}
						moves = append(moves, move)
					} else {
						for index, _ := range combos {
							combos[index].capturedPieces = append(
								combos[index].capturedPieces,
								position{i+vert, j+(1-horiz)},
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

func (b Board) checkDirections(i, j int, verticalMoves []int, player byte) []interface{} {
	var moves []interface{} = nil

	//Get all combo moves
	combos := b.comboCheck(i, j, i, j, verticalMoves, player)
	for _, m := range combos {
		moves = append(moves, m)
	}

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
					end: position{i+vert, j-horiz},
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
					end: position{i+vert, j+(1-horiz)},
				}
				moves = append(moves, move)
			}
		}
	}
	return moves
}

func (b Board) getMovesFromPos(i, j int) []interface{} {
	verticalMoves := make([]int, 0)
	if b[i][j] == 'X' || b[i][j] == 'O' {
		verticalMoves = append(verticalMoves, 1, -1)
	} else if b[i][j] == 'x' {
		verticalMoves = append(verticalMoves, 1)
	} else if b[i][j] == 'o' {
		verticalMoves = append(verticalMoves, -1)
	}

	moves := b.checkDirections(i, j, verticalMoves, b[i][j])
	return moves
}

func (g Game) GetActions() []interface{} {
	var moves []interface{} = make([]interface{}, 0)
	for i := 0; i < ROWS; i++ {
		for j := 0; j < COLS; j++ {
			if g.Board[i][j] == '_' {
				continue
			}
			if g.oTurn && (g.Board[i][j] == 'x' || g.Board[i][j] == 'X') ||
			  !g.oTurn && (g.Board[i][j] == 'o' || g.Board[i][j] == 'O') {
				continue
			}
			moves = append(moves, g.getMovesFromPos(i, j)...)
		}
	}

	//Weed out actions that lead to the same result
	uniqueActions := make(map[Game]bool)
	var ret []interface{} = nil

	for _, m := range moves {
		//Here, ApplyAction() can't error out
		newGameState, _ := g.ApplyAction(m)
		_, ok := uniqueActions[newGameState]
		if !ok {
			ret = append(ret, m)
			uniqueActions[newGameState] = true
		}
	}

	return ret
}

func (g Game) ApplyAction(move interface{}) (Game, error) {
	m, ok := move.(Move)
	if !ok {
		return g, ERR_INTERFACE_FAIL
	}

	//Move starting piece
	if m.end.y != m.start.y || m.end.x != m.start.x {
		g.Board[m.end.y][m.end.x] = g.Board[m.start.y][m.start.x]
		g.Board[m.start.y][m.start.x] = '_'
	}

	//Remove captured pieces
	for _, p := range m.capturedPieces {
		g.Board[p.y][p.x] = '_'
	}

	//Upgrade
	if m.end.y == ROWS - 1 && g.Board[m.end.y][m.end.x] == 'x' {
		g.Board[m.end.y][m.end.x] = 'X'
	} else if m.end.y == 0 && g.Board[m.end.y][m.end.x] == 'o' {
		g.Board[m.end.y][m.end.x] = 'O'
	}

	//Switch turns
	g.oTurn = !g.oTurn

	return g, nil
}

func main() {
	game := NewGameCombo3()
	fmt.Println(game)
	moves := game.GetActions()
	for _, m := range moves {
		fmt.Println(m)
		game2, err := game.ApplyAction(m)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(game2)
	}
}