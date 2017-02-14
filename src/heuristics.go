package main

import "math/rand"

// NOTE: maybe split into multiple files if this gets too big

func GoStraightHeuristic(gameState *GameState) string {

	mySnake := gameState.MySnake()
	head := mySnake.Coords[0]
	neck := mySnake.Coords[1]
	directionOfMovement := Point{
		X: head.X - neck.X,
		Y: head.Y - neck.Y,
	}
	allDirections := []string{UP, DOWN, LEFT, RIGHT}

	// try nd go straight
	for _, direction := range allDirections {
		if directionOfMovement == directionVector(direction) {
			possibleNewHead := head.Add(directionOfMovement)
			if !gameState.IsSolid(possibleNewHead, mySnake.Id) {
				return direction
			}
		}
	}
	return NOOP
}

func RandomHeuristic(gameState *GameState) string {

	mySnake := gameState.MySnake()
	head := mySnake.Coords[0]
	allDirections := []string{UP, DOWN, LEFT, RIGHT}

	validDirections := []string{}
	for _, direction := range allDirections {
		directionOfMovement := directionVector(direction)
		possibleNewHead := head.Add(directionOfMovement)
		if !gameState.IsSolid(possibleNewHead, mySnake.Id) {
			validDirections = append(validDirections, direction)
		}
	}

	if len(validDirections) == 0 {
		return NOOP
	}

	i := rand.Int() % len(validDirections)
	return validDirections[i]
}
