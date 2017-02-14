package main

import (
	"math/rand"
	"time"
)

// NOTES:
// maybe split into multiple files if this gets too big
// these funcs will take in board state when the other branch gets in



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
			if !gameState.IsSolid(possibleNewHead) {
				return direction
			}
			break
		}
	}

	// random other direction
	availableDirections := []string{}
	for _, possibleDirection := range allDirections {
		possibleNewHead := head.Add(directionOfMovement)
		if !gameState.IsSolid(possibleNewHead) {
			availableDirections = append(availableDirections, possibleDirection)
		}
	}

	if len(availableDirections) > 0 {
		rand.Seed(time.Now().Unix())
		n := rand.Int() % len(availableDirections)
		return availableDirections[n]
	}

	// SO DEAD!
	return UP
}