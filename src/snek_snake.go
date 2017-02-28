package main

import (
	"github.com/sendwithus/snek/src/"
)

func NewSnekSnake() SnakeAI {
	return SnekSnake{}
}

func (snake SnekSnake) SetDiedOnTurn(turn int) {
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
	mr := MoveRequest{}
	Mov
	return mr.GenerateMove()
}

func (gameState *GameState) GetSnekMoveRequest() MoveRequest {
	snek.src.MoveRequest{}
}
