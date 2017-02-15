package main

import (
	"fmt"
	"math/rand"
)

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

func (gameState *GameState) MySnake() *Snake {
	return gameState.GetSnake(gameState.You)
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
		gameState.You = snake.Id
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
	foodEaten := []Point{}
	for snakeId, newHead := range newHeads {
		snake := gameState.GetSnake(snakeId)
		if gameState.FoodAt(&newHead) {
			foodEaten = append(foodEaten, newHead)
			snake.HealthPoints = 100
		} else {
			snake.HealthPoints -= 1
			if snake.HealthPoints <= 0 {
				nextGameState.KillSnake(snakeId)
			}
			snake.Coords = snake.Coords[len(snake.Coords)-1:]

		}
	}

	// remove food
	for _, eatenFood := range foodEaten {
		newFoodList := []Point{}
		for _, food := range gameState.Food {
			if !(food.X == eatenFood.X && food.Y == eatenFood.Y) {
				newFoodList = append(newFoodList, food)
			}
		}
		gameState.Food = newFoodList
	}

	// collision
	for snakeId, newHead := range newHeads {
		if gameState.IsSolid(newHead, snakeId) {
			nextGameState.KillSnake(snakeId)
		}
	}

	// check win/draw states
	numSnakes := len(nextGameState.Snakes)
	if numSnakes == 1 {
		nextGameState.state = "Won"
		nextGameState.winners = []HeuristicSnake{
			*nextGameState.GetHeuristicSnake(nextGameState.Snakes[0].Id),
		}
	}
	if numSnakes == 0 {
		nextGameState.state = "Draw"
		nextGameState.winners = []HeuristicSnake{}
		for _, snake := range gameState.Snakes {
			heuristicSnake := gameState.GetHeuristicSnake(snake.Id)
			nextGameState.winners = append(nextGameState.winners, *heuristicSnake)
		}
	}

	return &nextGameState
}

func (gameState *GameState) SpawnFood() {
	emptyPoints := []Point{}
	for x := 0; x < gameState.Width; x += 1 {
		for y := 0; y < gameState.Height; y += 1 {
			p := Point{X: x, Y: y}
			if !gameState.IsSolid(p, "") {
				emptyPoints = append(emptyPoints, p)
			}
		}
	}

	i := rand.Int() % len(emptyPoints)
	gameState.Food = append(gameState.Food, emptyPoints[i])
}

func (gameState *GameState) FoodAt(p *Point) bool {
	for _, food := range gameState.Food {
		if food.X == p.X && food.Y == p.Y {
			return true
		}
	}
	return false
}

func (gameState *GameState) IsSolid(point Point, ignoreSnakeHead string) bool {
	if point.X < 0 || point.X >= gameState.Width {
		return true
	}
	if point.Y < 0 || point.Y >= gameState.Height {
		return true
	}

	// TODO: take in to account tale shrinks (when no food was eaten)
	for _, snake := range gameState.Snakes {
		for i, coord := range snake.Coords {
			if coord.X == point.X && coord.Y == point.Y {
				if !(i == 0 && snake.Id == ignoreSnakeHead) {
					return true
				}
			}
		}
	}
	return false
}

func (gameState *GameState) KillSnake(snakeId string) {
	fmt.Printf("Killing snake %v\n", snakeId)
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

func (gameState *GameState) GetHeuristicSnake(snakeId string) *HeuristicSnake {
	for i, snake := range gameState.HeuristicSnakes {
		if snake.Id == snakeId {
			return &gameState.HeuristicSnakes[i]
		}
	}
	return nil
}
