package main

import (
	"fmt"
	"math/rand"
	"time"
)

func NewGame(name string, snakeNames []string, foodFrequency int) *Game {
	initialMoveRequest := MoveRequest{
		Food:   [][]int{},
		GameId: "the one and only game atm",
		Height: 20,
		Width:  20,
		Turn:   0,
		Snakes: []MoveRequestSnake{},
		You:    "",
	}

	for _, snakeName := range snakeNames {
		x := rand.Intn(initialMoveRequest.Width)
		y := rand.Intn(initialMoveRequest.Height)
		snake := MoveRequestSnake{
			Id:   snakeName,
			Name: snakeName,
			Coords: [][]int{
				{x, y},
				{x, y},
				{x, y},
			},
			HealthPoints: 100,
		}
		initialMoveRequest.Snakes = append(initialMoveRequest.Snakes, snake)
	}

	initialGameState := NewGameState(initialMoveRequest)
	for i := 0; i < foodFrequency; i++ {
		initialGameState.SpawnFood()
	}

	return &Game{
		currentGameState: &initialGameState,
		foodFrequency:    foodFrequency,
		name:             name,
	}
}

func (game *Game) Run() []SnakeAI {
	start := time.Now()
	for {
		game.currentGameState = game.currentGameState.NextGameState()
		if game.currentGameState.state != "running" {
			break
		}
	}
	game.duration = time.Since(start)
	return game.currentGameState.winners
}

func (game *Game) Print() {
	if game.currentGameState.state != "running" {
		fmt.Printf("Game name:%v  \tturn:%v  \twinners:%v\tduration:%v\n",
			game.name,
			game.currentGameState.Turn,
			len(game.currentGameState.winners),
			game.duration,
		)
	}
}
