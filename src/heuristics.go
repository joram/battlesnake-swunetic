package main

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
