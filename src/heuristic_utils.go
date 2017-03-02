package main

import "fmt"

func MoveTo(gameState *GameState, goals []*Point) WeightedDirections {
	WeightUp := 0
	WeightDown := 0
	WeightLeft := 0
	WeightRight := 0
	snake := gameState.MySnake()
	if snake == nil {
		return []WeightedDirection{{Direction: NOOP, Weight: 0}}
	}

	head := snake.Coords[0]

	for _, goal := range goals {
		aStar, set := gameState.aStar[snake.Id]
		if !set {
			for _, snakeId := range gameState.aStar {
				fmt.Printf("\t%v", snakeId)
			}
			println("couldn't find aStar for ", snake.Id)
			continue
		}
		path := aStar.pathTo(goal)

		direction := path[0].Subtract(head)
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

func BoardControl(gameState *GameState, start *Point) float64 {
	if gameState.IsPossiblySolid(start, "") {
		return float64(0)
	}
	toVisit := []*Point{start}
	haveVisited := []*Point{}
	for len(toVisit) > 0 {
		p := toVisit[len(toVisit)-1]
		toVisit = toVisit[:len(toVisit)-1]
		haveVisited = append(haveVisited, p)
		for _, neighbour := range p.Neighbours() {
			if !gameState.IsPossiblySolid(neighbour, "") {
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
