package checkers

import "errors"

var (
	ERR_GAME_NOT_OVER      = errors.New("Game is not finished")
	ERR_INVALID_GAME_STATE = errors.New("Invalid game state")
	ERR_MOVE_NOT_IN_BOUNDS = errors.New("Move is not in bounds")
)

type Game struct {
	Board
	oTurn bool
}

func NewGame() Game {
	var board Board
	for i := 0; i < ROWS; i++ {
		for j := 0; j < COLS; j++ {
			board[i][j] = '_'
		}
	}
	for i := 0; i < (ROWS/2)-1; i++ {
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

func (g Game) GetActions() []Move {
	var moves []Move = nil
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
	var ret []Move = nil

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

func (g Game) ApplyAction(m Move) (Game, error) {
	if !m.inBounds() {
		return Game{}, ERR_MOVE_NOT_IN_BOUNDS
	}

	//Move starting piece
	if m.end.y != m.start.y || m.end.x != m.start.x {
		g.Board[m.end.y][m.end.x] = g.Board[m.start.y][m.start.x]
		g.Board[m.start.y][m.start.x] = '_'
	}

	removePieces := m.getCapturedPieces()
	for _, p := range removePieces {
		g.Board[p.y][p.x] = '_'
	}

	//Upgrade
	if m.end.y == ROWS-1 && g.Board[m.end.y][m.end.x] == 'x' {
		g.Board[m.end.y][m.end.x] = 'X'
	} else if m.end.y == 0 && g.Board[m.end.y][m.end.x] == 'o' {
		g.Board[m.end.y][m.end.x] = 'O'
	}

	//Switch turns
	g.oTurn = !g.oTurn

	return g, nil
}

func (g Game) IsTerminalState() bool {
	//Count the number of o's and x's on the field
	//If there are at least 1 of each, the game isn't
	//finished yet. Otherwise, the game is over
	var oCount int
	var xCount int

	for i := 0; i < ROWS; i++ {
		for j := 0; j < COLS; j++ {
			if g.Board[i][j] == '_' {
				continue
			}

			if g.Board[i][j] == 'o' || g.Board[i][j] == 'O' {
				oCount++
			} else {
				xCount++
			}

			if oCount > 0 && xCount > 0 {
				return false
			}
		}
	}

	return true
}

func (g Game) WinningPlayers() ([]byte, error) {
	if !g.IsTerminalState() {
		return nil, ERR_GAME_NOT_OVER
	}

	for i := 0; i < ROWS; i++ {
		for j := 0; j < COLS; j++ {
			if g.Board[i][j] == '_' {
				continue
			}

			if g.Board[i][j] == 'o' || g.Board[i][j] == 'O' {
				return []byte{'o'}, nil
			} else {
				return []byte{'x'}, nil
			}
		}
	}
	return nil, ERR_INVALID_GAME_STATE
}
