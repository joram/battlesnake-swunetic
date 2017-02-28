package main

import (
	"fmt"
	snek "github.com/sendwithus/snek/src"
)

func NewSnekSnake() SnakeAI {
	return SnekSnake{}
}

func (snake SnekSnake) SetDiedOnTurn(turn int) {
	fmt.Printf("snek died on %v\n", turn)
	snake.DiedOnTurn = turn
}

func (snake SnekSnake) GetDiedOnTurn() int {
	return snake.DiedOnTurn
}

func (snake SnekSnake) Mutate(weight int) {
}

func (snake SnekSnake) GetId() string {
	return "snek"
}

func (snake SnekSnake) GetWeights() map[string]int {
	return map[string]int{}
}

func (snake SnekSnake) Move(gameState *GameState) string {
	return gameState.GetSnekMoveRequest().GenerateMove()
}

func (gameState *GameState) GetSnekMoveRequest() snek.MoveRequest {
	mr := snek.MoveRequest{
		Food:   [][]int{},
		GameId: "fake",
		Height: gameState.Height,
		Width:  gameState.Width,
		Turn:   gameState.Turn,
		Snakes: []snek.Snake{},
		You:    "snek",
	}

	// populate food
	for _, food := range gameState.Food {
		mr.Food = append(mr.Food, []int{food.X, food.Y})
	}

	// populate snakes
	for _, snake := range gameState.Snakes {
		snekSnake := snek.Snake{
			Coords:       [][]int{},
			HealthPoints: snake.HealthPoints,
			Id:           snake.Id,
			Name:         snake.Name,
			Taunt:        snake.Taunt,
		}
		for _, coord := range snake.Coords {
			snekSnake.Coords = append(snekSnake.Coords, []int{coord.X, coord.Y})
		}
		mr.Snakes = append(mr.Snakes, snekSnake)
	}
	return mr
}
