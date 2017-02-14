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
	return NOOP
}
