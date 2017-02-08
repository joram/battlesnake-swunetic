package main

import "fmt"

type Game struct {
	snakes       []HeuristicSnake
	moveCount    int
	currentBoard [][]BoardCell
	state        string
	winners      []HeuristicSnake
}

func NewGame(numberOfSnakes int, weights [][]int) Game {
	// when creating a new game, we need to start up a certain number of snakes, each one
	// needs slightly different weights to decide which one is better and try and work
	// towards optimal weighting.
	snakes := [numberOfSnakes]HeuristicSnake{}
	for _, weights := range weights {
		snakes = append(snakes, NewHeuristicSnake(weights))
	}

	return Game{
		snakes:       snakes,
		moveCount:    0,
		currentBoard: [20][20]BoardCell{},
		winners:      []HeuristicSnake{},
	}
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
	board := game.currentBoard
	for _, snake := range game.snakes {
		direction := snake.Move(&board)
		game.MakeMove(&snake, direction)
	}
	game.moveCount += 1
}

func (game *Game) MakeMove(snake *HeuristicSnake, direction string) {
	// TODO: update game.currentBoard based on direction
	// Possible side effects:
	//	- sensible move, board's updated
	//	- death
	//	- win state=win
	//	- draw state=draw
}
