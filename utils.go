package main

func rowParity(i int) bool {
	return i%2 == 0
}

func inBounds(i, j int) bool {
	return i >= 0 && i < ROWS && j >= 0 && j < COLS
}
