package main

import (
	"fmt"
	"math/rand"
	"time"
)

func NewGame(name string, snakeNames []string, amountOfFood, width, height int) *Game {
	initialMoveRequest := MoveRequest{
		Food:   [][]int{},
		GameId: "the one and only game atm",
		Height: height,
		Width:  width,
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
	for i := 0; i < amountOfFood; i++ {
		initialGameState.SpawnFood()
	}

	return &Game{
		currentGameState: &initialGameState,
		foodFrequency:    amountOfFood,
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
