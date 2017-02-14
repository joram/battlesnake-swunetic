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
		game.currentGameState = game.currentGameState.NextGameState()
		if game.currentGameState.state != "running" {
			game.Print()
			break
		}
		game.Print()
	}
	return game.currentGameState.winners
}

func (game *Game) Print() {
	fmt.Printf("game is %v on turn %v\n", game.currentGameState.state, game.currentGameState.Turn)
	for _, snake := range game.currentGameState.HeuristicSnakes {
		weights := []int{}
		for _, w := range snake.WeightedHeuristics {
			weights = append(weights, w.weight)
		}
		fmt.Printf("winner weight: %v\n", weights)
	}
}
