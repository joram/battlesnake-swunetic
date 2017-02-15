package main

import "math/rand"

// NOTE: maybe split into multiple files if this gets too big

func NearestFoodHeuristic(gameState *GameState) string {

	var closestFood *Vector

	snake := gameState.MySnake()
	head := snake.Coords[0]
	for _, p := range gameState.Food {
		test := getDistanceBetween(head, p)
		if closestFood == nil {
			closestFood = test
		} else if test.Magnitude() < closestFood.Magnitude() {
			closestFood = test
		}
	}

	if closestFood.X < head.X && !gameState.IsSolid(head.Add(directionVector(LEFT)), snake.Id) {
		return LEFT
	} else if closestFood.X > head.X && !gameState.IsSolid(head.Add(directionVector(RIGHT)), snake.Id) {
		return RIGHT
	} else if closestFood.Y < head.Y && !gameState.IsSolid(head.Add(directionVector(UP)), snake.Id) {
		return UP
	} else if closestFood.Y > head.Y && !gameState.IsSolid(head.Add(directionVector(DOWN)), snake.Id) {
		return DOWN
	}
	return NOOP
}

func GoStraightHeuristic(gameState *GameState) string {

	mySnake := gameState.MySnake()
	head := mySnake.Coords[0]
	neck := mySnake.Coords[1]
	directionOfMovement := Point{
		X: head.X - neck.X,
		Y: head.Y - neck.Y,
	}
	allDirections := []string{UP, DOWN, LEFT, RIGHT}

	// try and go straight
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
