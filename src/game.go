package main

import "fmt"

func NewGame(numSnakes int) *Game {
	initialMoveRequest := MoveRequest{
		Food: []Point{},
		GameId: "the one and only game atm",
		Height: 20,
		Width: 20,
		Turn:0,
		Snakes: []Snake{},
		You: "",
	}

	for i := 0; i < numSnakes; i += 1 {
		snake := Snake{
			Id: fmt.Sprintf("Snake-%v", i),
			Name: fmt.Sprintf("Snake-%v", i),
			Taunt: "poop",
			Coords: []Point{
				Point{X:i, Y:i},
				Point{X:i, Y:i},
				Point{X:i, Y:i},
			},
			HealthPoints: 100,
		}
		initialMoveRequest.Snakes = append(initialMoveRequest.Snakes, snake)
	}

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