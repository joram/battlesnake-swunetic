package main

import (
	"math"
	"math/rand"
)

// NOTE: maybe split into multiple files if this gets too big

func NearestFoodHeuristic(gameState *GameState) WeightedDirections {

	var closestFood *Vector
	var food Point

	snake := gameState.MySnake()
	head := snake.Coords[0]
	for _, p := range gameState.Food {
		test := getDistanceBetween(head, p)
		if closestFood == nil {
			closestFood = test
			food = p
		} else if test.Magnitude() < closestFood.Magnitude() {
			closestFood = test
			food = p
		}
	}

	if closestFood == nil {
		return []WeightedDirection{{Direction: NOOP, Weight: 0}}
	}

	if head.Left().isCloser(&head, &food) && !gameState.IsSolid(head.Add(directionVector(LEFT)), snake.Id) {
		return []WeightedDirection{{Direction: LEFT, Weight: 100 - snake.HealthPoints}}
	}
	if head.Right().isCloser(&head, &food) && !gameState.IsSolid(head.Add(directionVector(RIGHT)), snake.Id) {
		return []WeightedDirection{{Direction: RIGHT, Weight: 100 - snake.HealthPoints}}
	}
	if head.Up().isCloser(&head, &food) && !gameState.IsSolid(head.Add(directionVector(UP)), snake.Id) {
		return []WeightedDirection{{Direction: UP, Weight: 100 - snake.HealthPoints}}
	}
	if head.Down().isCloser(&head, &food) && !gameState.IsSolid(head.Add(directionVector(DOWN)), snake.Id) {
		return []WeightedDirection{{Direction: DOWN, Weight: 100 - snake.HealthPoints}}
	}

	return []WeightedDirection{{Direction: NOOP, Weight: 0}}
}

func GoStraightHeuristic(gameState *GameState) string {

	mySnake := gameState.MySnake()

	if len(mySnake.Coords) <= 1 {
		return NOOP
	}

	head := mySnake.Coords[0]
	neck := mySnake.Coords[1]
	directionOfMovement := Point{
		X: head.X - neck.X,
		Y: head.Y - neck.Y,
	}
	allDirections := []string{UP, DOWN, LEFT, RIGHT}

	// try and go straight
	for _, direction := range allDirections {
		if directionOfMovement.Equals(directionVector(direction)) {
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

func BoardControl(gameState *GameState, start *Point) float64 {
	if !gameState.IsEmpty(*start) {
		return float64(0)
	}
	toVisit := []*Point{start}
	haveVisited := []*Point{}
	for len(toVisit) > 0 {
		p := toVisit[len(toVisit)-1]
		toVisit = toVisit[:len(toVisit)-1]
		haveVisited = append(haveVisited, p)
		for _, neighbour := range p.Neighbours() {
			if gameState.IsEmpty(*neighbour) {
				shouldCheck := true

				for _, visitedPoint := range haveVisited {
					if visitedPoint.Equals(*neighbour) {
						shouldCheck = false
					}
				}

				for _, toVisitedPoint := range toVisit {
					if toVisitedPoint.Equals(*neighbour) {
						shouldCheck = false
					}
				}

				if shouldCheck {
					toVisit = append(toVisit, neighbour)
				}
			}
		}
	}
	canVisit := float64(len(haveVisited))
	return canVisit
}

func BoardControlHeuristic(gameState *GameState) WeightedDirections {
	mySnake := gameState.MySnake()
	head := mySnake.Coords[0]
	controlLeft := BoardControl(gameState, head.Left())
	controlRight := BoardControl(gameState, head.Right())
	controlUp := BoardControl(gameState, head.Up())
	controlDown := BoardControl(gameState, head.Down())
	maxControl := math.Max(controlLeft, math.Max(controlRight, math.Max(controlUp, controlDown)))

	weightedDirections := []WeightedDirection{
		{Weight: int(controlLeft / maxControl * 100), Direction: LEFT},
		{Weight: int(controlRight / maxControl * 100), Direction: RIGHT},
		{Weight: int(controlUp / maxControl * 100), Direction: UP},
		{Weight: int(controlDown / maxControl * 100), Direction: DOWN},
	}

	return weightedDirections
}
