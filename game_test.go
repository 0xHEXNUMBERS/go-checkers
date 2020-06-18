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

	_, err = gameCont.Winner()
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

	player, err := gameFinish.Winner()
	if err != nil {
		return fmt.Errorf("Winning capture action resulted in an invalid state: %s", err)
	}

	if player != winner {
		return fmt.Errorf("Winning capture action did not result in the correct player winning: want '%c', got '%c'",
			winner,
			player,
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
	var b board
	for i := 0; i < ROWS; i++ {
		for j := 0; j < COLS; j++ {
			b[i][j] = '_'
		}
	}
	b[2][2] = 'X'
	b[3][2] = 'o'
	b[3][3] = 'o'
	return Game{
		board:     b,
		combo:     position{-1, -1},
		turnTimer: 0,
		oTurn:     true,
	}
}

func TestGameCapture(t *testing.T) {
	game := NewGameCapture()
	actionsGot := game.GetActions()

	actionsWant := []Move{
		{
			start:         position{3, 2},
			end:           position{1, 3},
			capturedPiece: position{2, 2},
		},
		{
			start:         position{3, 3},
			end:           position{1, 2},
			capturedPiece: position{2, 2},
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
	var b board
	for i := 0; i < ROWS; i++ {
		for j := 0; j < COLS; j++ {
			b[i][j] = '_'
		}
	}
	b[ROWS-2][COLS-1] = 'x'
	b[1][0] = 'o'
	return Game{
		board:     b,
		combo:     position{-1, -1},
		turnTimer: 0,
		oTurn:     true,
	}
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

	if game2.board[0][0] != 'O' {
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

	if game3.board[ROWS-1][COLS-1] != 'X' {
		t.Errorf("'x' piece did not get upgraded")
	}
}

func NewGameCombo() Game {
	var b board
	for i := 0; i < ROWS; i++ {
		for j := 0; j < COLS; j++ {
			b[i][j] = '_'
		}
	}
	b[2][2] = 'X'
	b[3][1] = 'o'
	b[3][2] = 'o'
	b[1][1] = 'o'
	b[1][2] = 'o'
	return Game{
		board:     b,
		combo:     position{-1, -1},
		turnTimer: 0,
		oTurn:     false,
	}
}

func TestGameCombo(t *testing.T) {
	game := NewGameCombo()
	game, _ = game.ApplyAction(Move{
		start:         position{2, 2},
		end:           position{4, 1},
		capturedPiece: position{3, 2},
	})
	game, _ = game.ApplyAction(game.GetActions()[0])
	game, _ = game.ApplyAction(game.GetActions()[0])

	//At this point, there should only be 1 piece left, and this
	//last action will continue the combo started at the beginning
	//of this test to remove this last piece.
	if err := testTerminalAction(game, game.GetActions()[0], 'x'); err != nil {
		t.Error(err)
	}
}
