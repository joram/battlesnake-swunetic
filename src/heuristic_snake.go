package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/sendwithus/lib-go"
	"math/rand"
	"sort"
	"sync"
)

func (h *WeightedHeuristic) Calculate(gameState *GameState) {
	direction := h.MoveFunc(gameState)
	h.Move = direction
}

func NewHeuristicSnake(id string) HeuristicSnake {
	snake := HeuristicSnake{
		Id:                 id,
		WeightedHeuristics: []*WeightedHeuristic{},
	}

	heuristics := map[string]MoveHeuristic{
		"nearest-food": NearestFoodHeuristic,
		"straight":     GoStraightHeuristic,
		"random":       RandomHeuristic,
	}

	for name, heuristic := range heuristics {
		weightedHeuristic := &WeightedHeuristic{
			Weight:   getWeight(name),
			MoveFunc: heuristic,
			Move:     NOOP,
			Name:     name,
		}
		snake.WeightedHeuristics = append(snake.WeightedHeuristics, weightedHeuristic)
	}
	return snake
}

func getWeight(name string) int {
	c, err := redis.Dial("tcp", swu.GetEnvVariable("REDIS_URL", true))
	if err != nil {
		panic(err)
	}
	defer c.Close()

	weight, err := redis.Int(c.Do("GET", name))
	if err != nil || weight == 0 {
		weight = rand.Intn(50) // figure out a good starting Weight for a new heuristic
	}
	return weight
}

func (snake *HeuristicSnake) Move(gameState *GameState) string {

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

	// calc weights of moves
	weights := map[string]int{
		UP:    0,
		DOWN:  0,
		LEFT:  0,
		RIGHT: 0,
		NOOP:  0,
	}
	for _, wh := range snake.WeightedHeuristics {
		weights[wh.Move] += wh.Weight
	}
	weights[NOOP] = 0

	ch := make(chan WeightedDirection)
	go sortWeightsMap(weights, ch)

	for weightedDirection := range ch {
		head := gameState.MySnake().Coords[0]
		directionOfMovement := directionVector(weightedDirection.Direction)
		possibleNewHead := head.Add(directionOfMovement)
		if !gameState.IsSolid(possibleNewHead, snake.Id) {
			fmt.Println("Choosing:", weightedDirection.Direction)
			return weightedDirection.Direction
		}
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
