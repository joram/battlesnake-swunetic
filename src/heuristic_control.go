package main

import "math"

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
	minControl := math.Min(controlLeft, math.Min(controlRight, math.Min(controlUp, controlDown)))
	if maxControl == minControl {
		return []WeightedDirection{}
	}

	weightedDirections := []WeightedDirection{
		{Weight: int(controlLeft / maxControl * 100), Direction: LEFT},
		{Weight: int(controlRight / maxControl * 100), Direction: RIGHT},
		{Weight: int(controlUp / maxControl * 100), Direction: UP},
		{Weight: int(controlDown / maxControl * 100), Direction: DOWN},
	}

	return weightedDirections
}
