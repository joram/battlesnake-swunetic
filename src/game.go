package main

import "fmt"

func NewGame(numSnakes int) *Game {
	initialMoveRequest := MoveRequest{} // TODO: build this out
	initialGameState := NewGameState(initialMoveRequest)
	return &Game{
		currentGameState: &initialGameState,
	}
}

func (game *Game) Run() []HeuristicSnake {
	for {
		game.Print()
		game.currentGameState = game.currentGameState.NextGameState()
		if game.currentGameState.state != "running" {
			break
		}
	}
	game.Print()
	return game.currentGameState.winners
}

func (game *Game) Print() {
	fmt.Printf("game is %v on turn %v\n", game.currentGameState.state, game.currentGameState.Turn)
	for _, snake := range game.currentGameState.HeuristicSnakes {
		println("winner snake had the weights:")
		for _, w := range snake.WeightedHeuristics {
			fmt.Printf("%v: %v", w.Name, w.weight)
		}
	}
}
