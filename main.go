package main

import "fmt"

func main() {
	game := NewGameCapture()
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
		fmt.Println("Terminal:", game2.IsTerminalState())
		winners, err := game2.WinningPlayers()
		fmt.Println("Winner:", winners)
		fmt.Println(err)
	}
}
