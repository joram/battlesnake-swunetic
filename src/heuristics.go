package main

import (
	"math"
	"math/rand"
)

// NOTE: maybe split into multiple files if this gets too big

func NearestFoodHeuristic(gameState *GameState) WeightedDirections {
	WeightUp := 0
	WeightDown := 0
	WeightLeft := 0
	WeightRight := 0
	snake := gameState.MySnake()
	head := snake.Coords[0]
	foods := []*Point{}
	for _, food := range gameState.Food {
		foods = append(foods, &food)
	}
	pathCalc := NewPathCalculation(&head, foods, gameState)
	pathCalc.Run()
	paths := pathCalc.Paths()
	for _, path := range paths {
		direction := path[len(path)-1].Subtract(head)
		if direction == directionVector(UP) {
			WeightUp += 25
		}
		if direction == directionVector(DOWN) {
			WeightDown += 25
		}
		if direction == directionVector(LEFT) {
			WeightLeft += 25
		}
		if direction == directionVector(RIGHT) {
			WeightRight += 25
		}
	}
	return []WeightedDirection{
		{Direction: UP, Weight: WeightUp},
		{Direction: DOWN, Weight: WeightDown},
		{Direction: LEFT, Weight: WeightLeft},
		{Direction: RIGHT, Weight: WeightRight},
	}
}

func GoStraightHeuristic(gameState *GameState) WeightedDirections {

	mySnake := gameState.MySnake()

	if len(mySnake.Coords) <= 1 {
		return []WeightedDirection{{Direction: NOOP, Weight: 0}}
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
				return []WeightedDirection{{Direction: direction, Weight: 50}}
			}
		}
	}
	return []WeightedDirection{{Direction: NOOP, Weight: 0}}
}

func RandomHeuristic(gameState *GameState) WeightedDirections {

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
		return []WeightedDirection{{Direction: NOOP, Weight: 0}}
	}

	i := rand.Int() % len(validDirections)
	return []WeightedDirection{{Direction: validDirections[i], Weight: 50}}
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
