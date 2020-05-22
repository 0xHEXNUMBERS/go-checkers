package checkers

import "errors"

var (
	//ErrGameNotOver error
	ErrGameNotOver = errors.New("Game is not finished")

	//ErrInvalidGameState error
	ErrInvalidGameState = errors.New("Invalid game state")

	//ErrMoveNotInBounds error
	ErrMoveNotInBounds = errors.New("Move is not in bounds")
)

//Game is the base struct that holds game state information.
type Game struct {
	board
	oTurn bool
}

//NewGame returns a new valid game of checkers.
func NewGame() Game {
	var b board
	for i := 0; i < ROWS; i++ {
		for j := 0; j < COLS; j++ {
			b[i][j] = '_'
		}
	}
	for i := 0; i < (ROWS/2)-1; i++ {
		for j := 0; j < COLS; j++ {
			b[i][j] = 'x'
		}
	}
	for i := (ROWS / 2) + 1; i < ROWS; i++ {
		for j := 0; j < COLS; j++ {
			b[i][j] = 'o'
		}
	}
	return Game{board: b}
}

//GetActions returns a list of moves that can be made
//by the current player.
func (g Game) GetActions() []Move {
	var moves []Move = nil
	for i := 0; i < ROWS; i++ {
		for j := 0; j < COLS; j++ {
			if g.board[i][j] == '_' {
				continue
			}
			if g.oTurn && (g.board[i][j] == 'x' || g.board[i][j] == 'X') ||
				!g.oTurn && (g.board[i][j] == 'o' || g.board[i][j] == 'O') {
				continue
			}
			moves = append(moves, g.getMovesFromPos(position{i, j})...)
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

//ApplyAction takes a Move and applies the action to the current game state.
//
//Returns the new game state and ErrMoveNotInBounds
//if the Move m is invalid.
func (g Game) ApplyAction(m Move) (Game, error) {
	if !m.inBounds() {
		return Game{}, ErrMoveNotInBounds
	}

	//Move starting piece
	if m.end.i != m.start.i || m.end.j != m.start.j {
		g.board[m.end.i][m.end.j] = g.board[m.start.i][m.start.j]
		g.board[m.start.i][m.start.j] = '_'
	}

	removePieces := m.getCapturedPieces()
	for _, p := range removePieces {
		g.board[p.i][p.j] = '_'
	}

	//Upgrade
	if m.end.i == ROWS-1 && g.board[m.end.i][m.end.j] == 'x' {
		g.board[m.end.i][m.end.j] = 'X'
	} else if m.end.i == 0 && g.board[m.end.i][m.end.j] == 'o' {
		g.board[m.end.i][m.end.j] = 'O'
	}

	//Switch turns
	g.oTurn = !g.oTurn

	return g, nil
}

//IsTerminalState returns whether the game is finished or not.
func (g Game) IsTerminalState() bool {
	//Count the number of o's and x's on the field
	//If there are at least 1 of each, the game isn't
	//finished yet. Otherwise, the game is over
	var oCount int
	var xCount int

	for i := 0; i < ROWS; i++ {
		for j := 0; j < COLS; j++ {
			if g.board[i][j] == '_' {
				continue
			}

			if g.board[i][j] == 'o' || g.board[i][j] == 'O' {
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

//Winner returns the winner's ascii value.
//
//Returns ErrGameNotOver if the game is not over.
//
//Returns ErrInvalidGameState if the game is in an invalid game state.
func (g Game) Winner() (byte, error) {
	if !g.IsTerminalState() {
		return '_', ErrGameNotOver
	}

	//Loop through the entire game board
	//and search for the winner
	for i := 0; i < ROWS; i++ {
		for j := 0; j < COLS; j++ {
			if g.board[i][j] == '_' {
				continue
			}

			if g.board[i][j] == 'o' || g.board[i][j] == 'O' {
				return 'o', nil
			}
			return 'x', nil
		}
	}
	return '_', ErrInvalidGameState
}

//Player returns the ascii value of the player
//that is currently deciding a move.
//
//Player returns 'o' if player o is making a move.
//Otherwise, Player return 'x'.
func (g Game) Player() byte {
	if g.oTurn {
		return 'o'
	}
	return 'x'
}
