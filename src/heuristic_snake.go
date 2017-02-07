package main

import (
	"sync"
)

type MoveHeuristic func(request *MoveRequest) string


type WeightedHeuristic struct {
	weight int
	move string
	moveHeuristic MoveHeuristic
}

func (weightedHeuristic *WeightedHeuristic) Calculate(request *MoveRequest){
	// TODO
}

type HeuristicSnake struct {
	weightedHeuristics []WeightedHeuristic
}

func NewHeuristicSnake(weights []int) HeuristicSnake {
	snake := HeuristicSnake{
		weightedHeuristics: []WeightedHeuristic{},
	}

        heuristics := []MoveHeuristic{
	//	 this is where we list all heuristics we've written
	}
	for i, weight := range weights {
		snake.weightedHeuristics = append(snake.weightedHeuristics, WeightedHeuristic{
			weight:weight,
			moveHeuristic:heuristics[i],
		})
	}
	return snake
}


func (snake *HeuristicSnake) Move(request *MoveRequest) string {
	var heuristicWaitGroup sync.WaitGroup
	heuristicWaitGroup.Add(len(snake.weightedHeuristics))

	// do heuristics
	for _, weightedHeuristic := range snake.weightedHeuristics {
		go func(h *WeightedHeuristic){
			h.Calculate(request)
			heuristicWaitGroup.Done()
		}(&weightedHeuristic)
	}
	heuristicWaitGroup.Wait()

	// calc weights of moves
	weights := map[string]int{
		"u": 0,
		"d": 0,
		"l": 0,
		"r": 0,
	}
	for _, weightedHeuristic := range snake.weightedHeuristics {
		weights[weightedHeuristic.move] += weightedHeuristic.weight
	}

	// pick heaviest weighted move
	bestDirection := "u"
	bestWeight := weights["u"]
	for direction, weight := range weights {
		if weight > bestWeight {
			bestDirection = direction
		}
	}

	return  bestDirection
}
