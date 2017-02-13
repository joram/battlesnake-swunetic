package main

import "fmt"

type Game struct {
	snakes       []HeuristicSnake
	moveCount    int
	currentBoard MoveRequest
	state        string
	winners      []HeuristicSnake
}

func NewGame(numberOfSnakes int) Game {
	// when creating a new game, we need to start up a certain number of snakes, each one
	// needs slightly different weights to decide which one is better and try and work
	// towards optimal weighting.
	snakes := []HeuristicSnake{}
	for i := 0; i < numberOfSnakes; i++ {
		snake := NewHeuristicSnake([]int{})
		snakes = append(snakes, snake)
	}

	return Game{
		snakes:       snakes,
		moveCount:    0,
		currentBoard: generateInitialBoard(20, 20),
		winners:      []HeuristicSnake{},
	}
}

func generateInitialBoard(width int, height int) MoveRequest {
	board := MoveRequest{}
	return board
}

func (game *Game) Run() []HeuristicSnake {
	game.state = "running"
	for {
		game.Tick()
		if game.state != "running" {
			fmt.Printf("game ended: %v on move %v\n", game.state, game.moveCount)
			for _, snake := range game.snakes {
				weights := []int{}
				for _, w := range snake.weightedHeuristics {
					weights = append(weights, w.weight)
				}
				fmt.Printf("winner: %v\n", weights)
			}
			break
		}
	}
	return game.winners
}

func (game *Game) Tick() {

	// get all moves
	moveDirections := map[string]string{}
	for _, snake := range game.snakes {
		moveDirections[snake.id] = snake.Move(game.buildMoveRequest())
	}

	// extend all snakes
	newHeads := map[string]Point{}
	for _, heuristicSnake := range game.snakes {
		snake := game.currentBoard.GetSnake(heuristicSnake.id)
		direction := moveDirections[snake.Id]
		newHeads[snake.Id] = snake.Extend(direction)
	}

	// eat or shrink


	// collision
	for snakeId, newHead := range newHeads {
		if game.currentBoard.IsSolid(newHead) {
			game.currentBoard.KillSnake(snakeId)
		}
	}

	// check win/draw states

	game.moveCount += 1
}

func (game *Game) buildMoveRequest() *MoveRequest {
	request := MoveRequest{}
	return &request
}
