package main

import "math/rand"

func RandomHeuristic(gameState *GameState) WeightedDirections {

	mySnake := gameState.MySnake()
	if mySnake == nil {
		return []WeightedDirection{{Direction: NOOP, Weight: 0}}
	}

	head := mySnake.Coords[0]
	allDirections := []string{UP, DOWN, LEFT, RIGHT}

	validDirections := []string{}
	for _, direction := range allDirections {
		directionOfMovement := directionVector(direction)
		possibleNewHead := head.Add(directionOfMovement)
		if !gameState.IsPossiblySolid(&possibleNewHead, mySnake.Id) {
			validDirections = append(validDirections, direction)
		}
	}

	if len(validDirections) == 0 {
		return []WeightedDirection{{Direction: NOOP, Weight: 0}}
	}

	i := rand.Int() % len(validDirections)
	return []WeightedDirection{{Direction: validDirections[i], Weight: 50}}
}
