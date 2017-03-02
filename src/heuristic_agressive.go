package main

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

func StraightLineAgressionHeuristic(gameState *GameState) WeightedDirections {
	Weights := map[string]int{
		UP:    0,
		DOWN:  0,
		LEFT:  0,
		RIGHT: 0,
	}
	me := gameState.MySnake()

	for _, them := range gameState.OtherSnakes() {
		currHead := *me.Head()
		for _, directionStr := range []string{UP, DOWN, LEFT, RIGHT} {
			direction := directionVector(directionStr)
			path := []Point{}
			canHitSolid := false
			for gameState.IsEmpty(&currHead) || len(path) == 0 {
				currHead = currHead.Add(direction)
				if !gameState.IsEmpty(&currHead) && len(path) > 0 {
					canHitSolid = true
					break
				}
				if gameState.aStar[me.Id].turnsTo[currHead] < gameState.aStar[them.Id].turnsTo[currHead] {
					path = append(path, currHead)
				} else {
					canHitSolid = false
					break
				}
			}

			if canHitSolid {
				futureGameState := *gameState
				futureMe := futureGameState.MySnake()
				futureMe.Coords = append(futureMe.Coords, path...)
				theirAStar := NewAStar(&futureGameState, them.Head())
				if theirAStar.canVisitCount < gameState.aStar[them.Id].canVisitCount/2 {
					Weights[directionStr] += 1
				}

			}
		}
	}

	maxWeight := Weights[UP]
	for _, directionStr := range []string{DOWN, LEFT, RIGHT} {
		if Weights[directionStr] > maxWeight {
			maxWeight = Weights[directionStr]
		}
	}

	options := []WeightedDirection{}
	if maxWeight > 0 {
		options = []WeightedDirection{
			{UP, Weights[UP] * 100 / maxWeight},
			{DOWN, Weights[DOWN] * 100 / maxWeight},
			{LEFT, Weights[LEFT] * 100 / maxWeight},
			{RIGHT, Weights[RIGHT] * 100 / maxWeight},
		}
	}
	return options
}
