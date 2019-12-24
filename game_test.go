package checkers

import (
	"fmt"
	"reflect"
	"testing"
)

func mappifyMoves(actions []Move) map[Move]bool {
	actionsCollected := make(map[Move]bool)
	for _, a := range actions {
		actionsCollected[a] = true
	}
	return actionsCollected
}

func containSameMoves(a, b []Move) bool {
	return reflect.DeepEqual(
		mappifyMoves(a), mappifyMoves(b),
	)
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
	return Game{Board: board, oTurn: true}
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

func TestGetComboActions(t *testing.T) {
	game := NewGameCombo()
	actionsGot := game.GetActions()

	startPos := position{2, 2}

	actionsWant := []Move{
		Move{
			start:          startPos,
			end:            position{2, 2},
			capturedPieces: "1-1|1-2|3-1|3-2",
		},
		Move{
			start:          startPos,
			end:            position{4, 1},
			capturedPieces: "3-2",
		},
		Move{
			start:          startPos,
			end:            position{2, 0},
			capturedPieces: "3-1|3-2",
		},
		Move{
			start:          startPos,
			end:            position{0, 1},
			capturedPieces: "1-1|3-1|3-2",
		},
		Move{
			start:          startPos,
			end:            position{0, 1},
			capturedPieces: "1-2",
		},
		Move{
			start:          startPos,
			end:            position{2, 0},
			capturedPieces: "1-1|1-2",
		},
		Move{
			start:          startPos,
			end:            position{4, 1},
			capturedPieces: "1-1|1-2|3-1",
		},
		Move{
			start:          startPos,
			end:            position{3, 3},
			capturedPieces: "",
		},
		Move{
			start:          startPos,
			end:            position{1, 3},
			capturedPieces: "",
		},
	}

	if len(actionsGot) != len(actionsWant) {
		t.Errorf(
			"Didn't get expected number of actions: got %d, wanted %d",
			len(actionsGot),
			len(actionsWant),
		)
	}

	fmt.Println(actionsGot)
	fmt.Println(actionsWant)

	if !containSameMoves(actionsGot, actionsWant) {
		t.Error("Actions collected are not the same as the actions we wanted")
	}

	gameFinish, err := game.ApplyAction(actionsWant[0])
	if err != nil {
		t.Errorf("Could not apply combo winning action: %s", err)
	}

	if !gameFinish.IsTerminalState() {
		t.Errorf("Winning combo action did not result in a terminal state")
	}

	player, err := gameFinish.WinningPlayers()
	if err != nil {
		t.Errorf("Winning combo action resulted in an invalid state: %s", err)
	}

	if player[0] != 'x' {
		t.Errorf("Winning combo action did not result in the correct player winning: want 'x', got '%c'",
			player[0],
		)
	}
}
