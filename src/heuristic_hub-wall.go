package main

func HuggWallsHeuristic(gameState *GameState) WeightedDirections {
	me := gameState.MySnake()
	head := me.Head()
	if head == nil {
		return []WeightedDirection{}
	}

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
