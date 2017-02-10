package main

import "fmt"

type Game struct {
	snakes       []HeuristicSnake
	moveCount    int
	currentBoard [][]BoardCell
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

func generateInitialBoard(width int, height int) [][]BoardCell {
	board := [][]BoardCell{}
	for y := 0; y < width; y++ {
		row := []BoardCell{}
		for x := 0; x < height; x++ {
			cell := BoardCell{}
			row = append(row, cell)
		}
		board = append(board, row)
	}
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
	for _, snake := range game.snakes {
		direction := snake.Move(game.buildMoveRequest())
		game.MakeMove(&snake, direction)
	}
	game.moveCount += 1
}

func (game *Game) buildMoveRequest() *MoveRequest {
	request := MoveRequest{}
	return &request
}

func (game *Game) MakeMove(snake *HeuristicSnake, direction string) {
	// TODO: update game.currentBoard based on direction
	// Possible side effects:
	//	- sensible move, board's updated
	//	- death
	//	- win state=win
	//	- draw state=draw
}
