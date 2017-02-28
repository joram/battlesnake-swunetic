package main

import (
	"fmt"
	"math"
	"math/rand"
)

// NOTE: maybe split into multiple files if this gets too big

func AvoidCollisionsWithBiggerSnakes(gameState *GameState) WeightedDirections {
	me := gameState.MySnake()

	pointsToAvoid := []*Point{}
	for _, otherSnake := range gameState.OtherSnakes() {
		if otherSnake.Length() >= me.Length() {
			pointsToAvoid = append(pointsToAvoid, otherSnake.Coords[0].Neighbours()...)
		}
	}

	weightedDirections := WeightedDirections{}

	option := me.Coords[0].Up()
	for _, pointToAvoid := range pointsToAvoid {
		if pointToAvoid.Equals(*option) {
			weightedDirections = append(weightedDirections, WeightedDirection{
				Direction: UP,
				Weight:    -100,
			})
			break
		}
	}

	option = me.Coords[0].Down()
	for _, pointToAvoid := range pointsToAvoid {
		if pointToAvoid.Equals(*option) {
			weightedDirections = append(weightedDirections, WeightedDirection{
				Direction: DOWN,
				Weight:    -100,
			})
			break
		}
	}

	option = me.Coords[0].Left()
	for _, pointToAvoid := range pointsToAvoid {
		if pointToAvoid.Equals(*option) {
			weightedDirections = append(weightedDirections, WeightedDirection{
				Direction: LEFT,
				Weight:    -100,
			})
			break
		}
	}

	option = me.Coords[0].Right()
	for _, pointToAvoid := range pointsToAvoid {
		if pointToAvoid.Equals(*option) {
			weightedDirections = append(weightedDirections, WeightedDirection{
				Direction: RIGHT,
				Weight:    -100,
			})
			break
		}
	}

	return weightedDirections
}

func AimeForCollisionsWithSmallerSnakes(gameState *GameState) WeightedDirections {
	me := gameState.MySnake()

	pointsToAimeFor := []*Point{}
	for _, otherSnake := range gameState.OtherSnakes() {
		if otherSnake.Length() < me.Length() {
			pointsToAimeFor = append(pointsToAimeFor, otherSnake.Coords[0].Neighbours()...)
		}
	}

	weightedDirections := WeightedDirections{}

	option := me.Coords[0].Up()
	for _, pointToAvoid := range pointsToAimeFor {
		if pointToAvoid.Equals(*option) {
			weightedDirections = append(weightedDirections, WeightedDirection{
				Direction: UP,
				Weight:    100,
			})
			break
		}
	}

	option = me.Coords[0].Down()
	for _, pointToAvoid := range pointsToAimeFor {
		if pointToAvoid.Equals(*option) {
			weightedDirections = append(weightedDirections, WeightedDirection{
				Direction: DOWN,
				Weight:    100,
			})
			break
		}
	}

	option = me.Coords[0].Left()
	for _, pointToAvoid := range pointsToAimeFor {
		if pointToAvoid.Equals(*option) {
			weightedDirections = append(weightedDirections, WeightedDirection{
				Direction: LEFT,
				Weight:    100,
			})
			break
		}
	}

	option = me.Coords[0].Right()
	for _, pointToAvoid := range pointsToAimeFor {
		if pointToAvoid.Equals(*option) {
			weightedDirections = append(weightedDirections, WeightedDirection{
				Direction: RIGHT,
				Weight:    100,
			})
		}
	}

	return weightedDirections
}

func CollisionHeuristic(gameState *GameState) WeightedDirections {
	smallHeads := []*Point{}
	me := gameState.MySnake()
	myLength := me.Length()
	for _, otherSnake := range gameState.OtherSnakes() {
		if otherSnake.Length() < myLength {
			smallHeads = append(smallHeads, &otherSnake.Coords[0])
		}
	}
	return MoveTo(gameState, smallHeads)
}

func HuggWallsHeuristic(gameState *GameState) WeightedDirections {
	me := gameState.MySnake()
	head := me.Coords[0]
	surroundingWallCountUp := gameState.CountSurroundingWalls(head.Up())
	surroundingWallCountDown := gameState.CountSurroundingWalls(head.Down())
	surroundingWallCountLeft := gameState.CountSurroundingWalls(head.Left())
	surroundingWallCountRight := gameState.CountSurroundingWalls(head.Right())
	total := surroundingWallCountUp + surroundingWallCountDown + surroundingWallCountLeft + surroundingWallCountRight

	return []WeightedDirection{
		{Direction: UP, Weight: 100 * surroundingWallCountUp / total},
		{Direction: DOWN, Weight: 100 * surroundingWallCountDown / total},
		{Direction: LEFT, Weight: 100 * surroundingWallCountLeft / total},
		{Direction: RIGHT, Weight: 100 * surroundingWallCountRight / total},
	}

}

func NearestFoodHeuristic(gameState *GameState) WeightedDirections {
	foods := []*Point{}
	for _, food := range gameState.Food {
		foods = append(foods, &food)
	}
	return MoveTo(gameState, foods)
}

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

func GoStraightHeuristic(gameState *GameState) WeightedDirections {

	mySnake := gameState.MySnake()
	if mySnake == nil {
		return []WeightedDirection{{Direction: NOOP, Weight: 0}}
	}

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
			if !gameState.IsSolid(&possibleNewHead, mySnake.Id) {
				return []WeightedDirection{{Direction: direction, Weight: 50}}
			}
		}
	}
	return []WeightedDirection{{Direction: NOOP, Weight: 0}}
}

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
		if !gameState.IsSolid(&possibleNewHead, mySnake.Id) {
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
	if !gameState.IsEmpty(start) {
		return float64(0)
	}
	toVisit := []*Point{start}
	haveVisited := []*Point{}
	for len(toVisit) > 0 {
		p := toVisit[len(toVisit)-1]
		toVisit = toVisit[:len(toVisit)-1]
		haveVisited = append(haveVisited, p)
		for _, neighbour := range p.Neighbours() {
			if gameState.IsEmpty(neighbour) {
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
	if mySnake == nil {
		return []WeightedDirection{{Direction: NOOP, Weight: 0}}
	}

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
