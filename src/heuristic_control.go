package main

import (
	"sync"
)

func BoardControlHeuristic(gameState *GameState) WeightedDirections {
	mySnake := gameState.MySnake()
	if mySnake == nil {
		return []WeightedDirection{{Direction: NOOP, Weight: 0}}
	}

	head := mySnake.Coords[0]
	wg := sync.WaitGroup{}
	lock := sync.Mutex{}
	control := map[string]float64{
		UP:    0,
		DOWN:  0,
		LEFT:  0,
		RIGHT: 0,
	}
	for dir, _ := range control {
		wg.Add(1)
		go func(dir string) {
			newHead := head.Add(directionVector(dir))
			c := BoardControl(gameState, &newHead)
			lock.Lock()
			control[dir] = c
			lock.Unlock()
			wg.Done()
		}(dir)
	}
	wg.Wait()

	maxControl := float64(-1)
	minControl := float64(-1)
	for _, val := range control {
		if maxControl == -1 || val > maxControl {
			maxControl = val
		}
		if minControl == -1 || val < minControl {
			minControl = val
		}
	}
	if maxControl == minControl {
		return []WeightedDirection{}
	}

	weightedDirections := []WeightedDirection{
		{Weight: int(control[LEFT] / maxControl * 100), Direction: LEFT},
		{Weight: int(control[RIGHT] / maxControl * 100), Direction: RIGHT},
		{Weight: int(control[UP] / maxControl * 100), Direction: UP},
		{Weight: int(control[DOWN] / maxControl * 100), Direction: DOWN},
	}

	return weightedDirections
}
