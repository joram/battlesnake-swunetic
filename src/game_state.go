package main

import "fmt"

func NewGameState(moveRequest MoveRequest) GameState {

	heuristicSnakes := []HeuristicSnake{}
	for _, snake := range moveRequest.Snakes {
		heuristicSnakes = append(heuristicSnakes, NewHeuristicSnake(snake.Id))
	}

	return GameState{
		HeuristicSnakes: heuristicSnakes,
		Snakes:          moveRequest.Snakes,
		Width:           moveRequest.Width,
		Height:          moveRequest.Height,
		Turn:            moveRequest.Turn,
		Food:            moveRequest.Food,
		winners:         []HeuristicSnake{},
		state:           "running",
	}
}

func (gameState *GameState) Run() []HeuristicSnake {
	gameState.state = "running"
	for {
		gameState = gameState.NextGameState()
		if gameState.state != "running" {
			fmt.Printf("game ended: %v on turn %v\n", gameState.state, gameState.Turn)
			for _, snake := range gameState.HeuristicSnakes {
				weights := []int{}
				for _, w := range snake.WeightedHeuristics {
					weights = append(weights, w.weight)
				}
				fmt.Printf("winner: %v\n", weights)
			}
			break
		}
	}
	return gameState.winners
}

func (gameState *GameState) NextGameState() *GameState {
	nextGameState := GameState{
		Turn:            gameState.Turn + 1,
		HeuristicSnakes: gameState.HeuristicSnakes,
		Snakes:          gameState.Snakes,
		Width:           gameState.Width,
		Height:          gameState.Height,
		Food:            gameState.Food,
		winners:         gameState.winners,
		state:           gameState.state,
	}

	// get all moves
	moveDirections := map[string]string{}
	for _, snake := range gameState.HeuristicSnakes {
		moveDirections[snake.Id] = snake.Move(gameState)
	}

	// extend all snakes
	newHeads := map[string]Point{}
	for _, heuristicSnake := range nextGameState.HeuristicSnakes {
		snake := nextGameState.GetSnake(heuristicSnake.Id)
		direction := moveDirections[snake.Id]
		newHeads[snake.Id] = snake.Extend(direction)
	}

	// eat or shrink

	// collision
	for snakeId, newHead := range newHeads {
		if nextGameState.IsSolid(newHead) {
			nextGameState.KillSnake(snakeId)
		}
	}

	// check win/draw states

	return &nextGameState
}

func (gameState *GameState) IsSolid(point Point) bool {
	if point.X < 0 || point.X >= gameState.Width {
		return true
	}
	if point.Y < 0 || point.Y >= gameState.Height {
		return true
	}

	// TODO: take in to account tale shrinks (when no food was eaten)
	for _, snake := range gameState.Snakes {
		for _, coord := range snake.Coords {
			if coord.X == point.X && coord.Y == point.Y {
				return true
			}
		}
	}
	return false
}

func (gameState *GameState) KillSnake(snakeId string) {
	newSnakes := []Snake{}
	for _, snake := range gameState.Snakes {
		if snake.Id != snakeId {
			newSnakes = append(newSnakes, snake)
		}
	}
	gameState.Snakes = newSnakes
}

func (gameState *GameState) GetSnake(snakeId string) *Snake {
	for i, snake := range gameState.Snakes {
		if snake.Id == snakeId {
			return &gameState.Snakes[i]
		}
	}
	return nil
}
