package checkers

import (
	"errors"
)

const (
	//TurnsBeforeDraw is the number of turns the engine
	//wil consider before resulting in a draw if a piece
	//capture has not occurred in that alloted time.
	TurnsBeforeDraw = 100
)

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
	combo     position
	turnTimer int
	oTurn     bool
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
	return Game{
		board:     b,
		combo:     position{-1, -1},
		turnTimer: 0,
		oTurn:     false,
	}
}

func (g *Game) deferAction() {
	g.turnTimer++
	g.oTurn = !g.oTurn
}

//GetActions returns a list of moves that can be made
//by the current player.
func (g Game) GetActions() []Move {
	if g.turnTimer == TurnsBeforeDraw {
		return nil
	}

	//If we're continuing a combo, only return capture
	//action from this position
	captures := make([]Move, 0)
	if g.combo.inBounds() {
		g.captureCheck(
			g.combo,
			g.getVertMovesFromPos(g.combo),
			g.board[g.combo.i][g.combo.j],
			&captures,
		)
		return captures
	}

	moves := make([]Move, 0)
	for i := 0; i < ROWS; i++ {
		for j := 0; j < COLS; j++ {
			if g.board[i][j] == '_' {
				continue
			}
			if g.oTurn && (g.board[i][j] == 'x' || g.board[i][j] == 'X') ||
				!g.oTurn && (g.board[i][j] == 'o' || g.board[i][j] == 'O') {
				continue
			}

			p := position{i, j}
			vertMoves := g.getVertMovesFromPos(p)
			g.captureCheck(
				p, vertMoves, g.board[i][j], &captures,
			)
			if len(captures) == 0 {
				g.checkForAdjacentVacantSpots(
					position{i, j}, vertMoves, &moves,
				)
			}
		}
	}

	ret := moves
	if len(captures) > 0 {
		ret = captures
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
	g.combo = position{-1, -1}

	//Move starting piece
	if m.end.i != m.start.i || m.end.j != m.start.j {
		g.board[m.end.i][m.end.j] = g.board[m.start.i][m.start.j]
		g.board[m.start.i][m.start.j] = '_'
	}

	zeroPosition := position{}
	if m.capturedPiece != zeroPosition {
		g.turnTimer = 0
		g.board[m.capturedPiece.i][m.capturedPiece.j] = '_'

		//Defer action to opponent if we can't continue a combo
		//Otherwise, save the end position to continue the combo
		//next turn
		if !g.canCaptureFromPos(m.end) {
			g.deferAction()
		} else {
			g.combo = m.end
		}
	} else {
		g.deferAction()
	}

	//Upgrade
	if m.end.i == ROWS-1 && g.board[m.end.i][m.end.j] == 'x' {
		g.board[m.end.i][m.end.j] = 'X'
	} else if m.end.i == 0 && g.board[m.end.i][m.end.j] == 'o' {
		g.board[m.end.i][m.end.j] = 'O'
	}

	return g, nil
}

func (g Game) pieceCounts() (xCount, oCount int) {
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
		}
	}
	return
}

func (g Game) canMove() bool {
	for i := 0; i < ROWS; i++ {
		for j := 0; j < COLS; j++ {
			if g.board[i][j] == '_' {
				continue
			}
			if g.oTurn && (g.board[i][j] == 'x' || g.board[i][j] == 'X') ||
				!g.oTurn && (g.board[i][j] == 'o' || g.board[i][j] == 'O') {
				continue
			}

			if g.canMoveFromPos(position{i, j}) {
				return true
			}
		}
	}
	return false
}

//IsTerminalState returns whether the game is finished or not.
func (g Game) IsTerminalState() bool {
	if g.turnTimer == TurnsBeforeDraw {
		return true
	}

	//Count the number of o's and x's on the field
	//If there are at least 1 of each, the game might not
	//be finished yet. Otherwise, the game is over.
	xCount, oCount := g.pieceCounts()
	if xCount == 0 || oCount == 0 {
		return true
	}

	//Can the current player make a move?
	//If not, the other player wins.
	return !g.canMove()
}

//Winner returns the winner's ascii value.
//
//If the game results in a draw, this method returns '_' as the winner.
//
//Returns ErrGameNotOver if the game is not over.
//
//Returns ErrInvalidGameState if the game is in an invalid game state.
func (g Game) Winner() (byte, error) {
	if !g.IsTerminalState() {
		return '_', ErrGameNotOver
	}

	//Draw
	if g.turnTimer == TurnsBeforeDraw {
		return '_', nil
	}

	xCount, oCount := g.pieceCounts()
	if xCount == 0 {
		return 'o', nil
	} else if oCount == 0 {
		return 'x', nil
	}

	//If the current player cannot move, then they lost
	if len(g.GetActions()) == 0 {
		if g.oTurn {
			return 'o', nil
		}
		return 'x', nil
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
