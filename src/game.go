package main

import (
	"fmt"
	"math/rand"
)

func NewGame(numSnakes int, foodFrequency int) *Game {
	initialMoveRequest := MoveRequest{
		Food:   []Point{},
		GameId: "the one and only game atm",
		Height: 20,
		Width:  20,
		Turn:   0,
		Snakes: []Snake{},
		You:    "",
	}

	for i := 0; i < numSnakes; i += 1 {
		snake := Snake{
			Id:    fmt.Sprintf("Snake-%v", i),
			Name:  fmt.Sprintf("Snake-%v", i),
			Taunt: "poop",
			Coords: []Point{
				Point{X: i, Y: i},
				Point{X: i, Y: i},
				Point{X: i, Y: i},
			},
			HealthPoints: 100,
		}
		initialMoveRequest.Snakes = append(initialMoveRequest.Snakes, snake)
	}

	initialGameState := NewGameState(initialMoveRequest)
	return &Game{
		currentGameState: &initialGameState,
		foodFrequency:    foodFrequency,
	}
}

func (game *Game) Run() []HeuristicSnake {
	for {
		game.Print()
		game.currentGameState = game.currentGameState.NextGameState()
		if game.SpawnFood() {
			println("spawning food")
			game.currentGameState.SpawnFood()
		}
		if game.currentGameState.state != "running" {
			break
		}
	}
	game.Print()
	return game.currentGameState.winners
}
func (game *Game) SpawnFood() bool {
	return rand.Int()%game.foodFrequency == 1
}

func (game *Game) Print() {
	if game.currentGameState.state != "running" {
		fmt.Printf("Game over on turn %v\n", game.currentGameState.Turn)
		for _, snake := range game.currentGameState.winners {
			winnerDetails := fmt.Sprintf("WINNER[%v] %v:\t", game.currentGameState.Turn, snake.Id)
			for _, w := range snake.WeightedHeuristics {
				winnerDetails += fmt.Sprintf("%v:%v ", w.Name, w.Weight)
			}
			println(winnerDetails)
		}

		for _, snake := range game.currentGameState.losers {
			winnerDetails := fmt.Sprintf("LOSER[%v] %v:\t", snake.DiedOnTurn, snake.Id)
			for _, w := range snake.WeightedHeuristics {
				winnerDetails += fmt.Sprintf("%v:%v ", w.Name, w.Weight)
			}
			println(winnerDetails)
		}

	}
}
