package checkers

import (
	"fmt"
	"reflect"
	"testing"
)

func testNonTerminalAction(g Game, m Move) error {
	gameCont, err := g.ApplyAction(m)
	if err != nil {
		return fmt.Errorf("Could not apply non-terminal action: %s", err)
	}

	if gameCont.IsTerminalState() {
		return fmt.Errorf("Non-terminal action resulted in a terminal state")
	}

	_, err = gameCont.WinningPlayers()
	if err == nil {
		return fmt.Errorf("Non-terminal action resulted in a state with winning players: %s", err)
	}
	return nil
}

func testTerminalAction(g Game, m Move, winner byte) error {
	gameFinish, err := g.ApplyAction(m)
	if err != nil {
		return fmt.Errorf("Could not apply winning capture action: %s", err)
	}

	if !gameFinish.IsTerminalState() {
		return fmt.Errorf("Winning capture action did not result in a terminal state")
	}

	player, err := gameFinish.WinningPlayers()
	if err != nil {
		return fmt.Errorf("Winning capture action resulted in an invalid state: %s", err)
	}

	if player[0] != winner {
		return fmt.Errorf("Winning capture action did not result in the correct player winning: want '%c', got '%c'",
			winner,
			player[0],
		)
	}

	return nil
}

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

func TestGameCapture(t *testing.T) {
	game := NewGameCapture()
	actionsGot := game.GetActions()

	actionsWant := []Move{
		Move{
			start:          position{3, 2},
			end:            position{1, 3},
			capturedPieces: "2-2",
		},
		Move{
			start:          position{3, 3},
			end:            position{1, 2},
			capturedPieces: "2-2",
		},
		Move{
			start: position{3, 2},
			end:   position{2, 1},
		},
		Move{
			start: position{3, 3},
			end:   position{2, 3},
		},
	}

	if len(actionsGot) != len(actionsWant) {
		t.Errorf(
			"Didn't get expected number of actions: got %d, wanted %d",
			len(actionsGot),
			len(actionsWant),
		)
	}

	if !containSameMoves(actionsGot, actionsWant) {
		t.Error("Actions collected are not the same as the actions we wanted")
	}

	//Loop through terminal actions
	for i := 0; i < 2; i++ {
		if err := testTerminalAction(game, actionsWant[i], 'o'); err != nil {
			t.Error(err)
		}
	}

	for i := 2; i < len(actionsWant); i++ {
		if err := testNonTerminalAction(game, actionsWant[i]); err != nil {
			t.Error(err)
		}
	}
}

func NewGameUpgrade() Game {
	var board Board
	for i := 0; i < ROWS; i++ {
		for j := 0; j < COLS; j++ {
			board[i][j] = '_'
		}
	}
	board[ROWS-2][COLS-1] = 'x'
	board[1][0] = 'o'
	return Game{board, true}
}

func TestGameUpgrade(t *testing.T) {
	game := NewGameUpgrade()
	oActionGot := game.GetActions()[0]

	oActionWant := Move{
		start: position{1, 0},
		end:   position{0, 0},
	}

	if oActionGot != oActionWant {
		t.Errorf("Did not get the wanted upgraded 'o' action: got: %s, want: %s",
			oActionGot,
			oActionWant,
		)
	}

	game2, err := game.ApplyAction(oActionWant)
	if err != nil {
		t.Errorf("Could not perform action to upgrade 'o' piece: %s", err)
	}

	if game2.Board[0][0] != 'O' {
		t.Errorf("'o' piece did not get upgraded")
	}

	xActionGot := game2.GetActions()[0]

	xActionWant := Move{
		start: position{ROWS - 2, COLS - 1},
		end:   position{ROWS - 1, COLS - 1},
	}

	if xActionGot != xActionWant {
		t.Errorf("Did not get the wanted upgraded 'x' action: got: %s, want: %s",
			xActionGot,
			xActionWant,
		)
	}

	game3, err := game2.ApplyAction(xActionWant)
	if err != nil {
		t.Errorf("Could not perform action to upgrade 'x' piece: %s", err)
	}

	if game3.Board[ROWS-1][COLS-1] != 'X' {
		t.Errorf("'x' piece did not get upgraded")
	}
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

func TestGameCombo(t *testing.T) {
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

	if !containSameMoves(actionsGot, actionsWant) {
		t.Error("Actions collected are not the same as the actions we wanted")
	}

	if err := testTerminalAction(game, actionsWant[0], 'x'); err != nil {
		t.Error(err)
	}

	for i := 1; i < len(actionsWant); i++ {
		if err := testNonTerminalAction(game, actionsWant[i]); err != nil {
			t.Error(err)
		}
	}
}
