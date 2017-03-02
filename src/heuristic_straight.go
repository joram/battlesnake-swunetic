package main

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
			if !gameState.IsPossiblySolid(&possibleNewHead, mySnake.Id) {
				return []WeightedDirection{{Direction: direction, Weight: 50}}
			}
		}
	}
	return []WeightedDirection{{Direction: NOOP, Weight: 0}}
}
