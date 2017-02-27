package main

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"
)

var heuristics = map[string]MoveHeuristic{
	"nearest-food": NearestFoodHeuristic,
	"straight":     GoStraightHeuristic,
	"random":       RandomHeuristic,
	"control":      BoardControlHeuristic,
	"agressive":    CollisionHeuristic,
	"attempt-kill": AimeForCollisionsWithSmallerSnakes,
	"avoid-death":  AvoidCollisionsWithBiggerSnakes,
	"hug-walls":    HuggWallsHeuristic,
}

func (h *WeightedHeuristic) Calculate(gameState *GameState) {
	start := time.Now()
	h.WeightedDirections = h.MoveFunc(gameState)

	if time.Since(start) > 500*time.Millisecond {
		fmt.Printf("Calculated %v in %v\n", h.Name, time.Since(start))
	}
}

func NewHeuristicSnake(id string) HeuristicSnake {
	snake := HeuristicSnake{
		Id:                 id,
		WeightedHeuristics: []*WeightedHeuristic{},
	}

	for name, heuristic := range heuristics {
		weightedHeuristic := &WeightedHeuristic{
			Weight:   int(getWeight(name)),
			MoveFunc: heuristic,
			Name:     name,
		}
		snake.WeightedHeuristics = append(snake.WeightedHeuristics, weightedHeuristic)
	}
	return snake
}

func (heuristicSnake *HeuristicSnake) Mutate(maxMutation int) {
	if maxMutation <= 0 {
		return
	}

	for i, _ := range heuristicSnake.WeightedHeuristics {
		originalWeight := heuristicSnake.WeightedHeuristics[i].Weight
		mutatedWeight := originalWeight + rand.Intn(maxMutation*2) - maxMutation // mutate between -x and +x
		if mutatedWeight > 100 {
			mutatedWeight = 100
		}
		if mutatedWeight < 0 {
			mutatedWeight = 0
		}
		heuristicSnake.WeightedHeuristics[i].Weight = mutatedWeight
	}
}

func (snake *HeuristicSnake) Move(gameState *GameState) string {
	start := time.Now()

	// do heuristics
	var heuristicWaitGroup sync.WaitGroup
	for i := range snake.WeightedHeuristics {
		heuristicWaitGroup.Add(1)
		go func(wh *WeightedHeuristic) {
			wh.Calculate(gameState)
			heuristicWaitGroup.Done()
		}(snake.WeightedHeuristics[i])
	}
	heuristicWaitGroup.Wait()
	if time.Since(start) > 500*time.Millisecond {
		fmt.Printf("Calculated heuristics %v\n", time.Since(start))
	}
	// calc weights of moves
	weights := map[string]int{
		UP:    0,
		DOWN:  0,
		LEFT:  0,
		RIGHT: 0,
		NOOP:  0,
	}
	for _, wh := range snake.WeightedHeuristics {
		for _, wd := range wh.WeightedDirections {
			weights[wd.Direction] += wd.Weight * wh.Weight
		}
	}
	weights[NOOP] = 0

	ch := make(chan WeightedDirection)
	go sortWeightsMap(weights, ch)

	for weightedDirection := range ch {
		//head := gameState.MySnake().Coords[0]
		//directionOfMovement := directionVector(weightedDirection.Direction)
		//possibleNewHead := head.Add(directionOfMovement)
		return weightedDirection.Direction
	}

	return NOOP
}

func sortWeightsMap(weights map[string]int, output chan WeightedDirection) {
	n := map[int][]string{}
	var a []int
	for k, v := range weights {
		n[v] = append(n[v], k)
	}
	for k := range n {
		a = append(a, k)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(a)))
	for _, direction := range a {
		for _, weight := range n[direction] {
			output <- WeightedDirection{
				Direction: weight,
				Weight:    direction,
			}
		}
	}
	close(output)
}
