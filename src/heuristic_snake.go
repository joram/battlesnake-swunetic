package main

import (
	"github.com/garyburd/redigo/redis"
	"math/rand"
	"sync"
)

type MoveHeuristic func(request *MoveRequest) string

type WeightedHeuristic struct {
	weight        int
	move          string
	moveHeuristic MoveHeuristic
}

func (weightedHeuristic *WeightedHeuristic) Calculate(request *MoveRequest) {
	weightedHeuristic.move = weightedHeuristic.moveHeuristic(request)
}

type HeuristicSnake struct {
	id string
	weightedHeuristics []WeightedHeuristic
}

func NewHeuristicSnake() HeuristicSnake {
	snake := HeuristicSnake{
		weightedHeuristics: []WeightedHeuristic{},
	}

	heuristics := map[string]MoveHeuristic{
	//	 this is where we list all heuristics we've written
	}
	for name, heuristic := range heuristics {
		snake.weightedHeuristics = append(snake.weightedHeuristics, WeightedHeuristic{
			weight:        getWeight(name),
			moveHeuristic: heuristic,
		})
	}
	return snake
}

func getWeight(name string) int {
	c, err := redis.Dial("tcp", "sendwithus.local.web-app.redis:6379") // TODO: update to redis on heroku
	if err != nil {
		panic(err)
	}
	defer c.Close()

	weight, err := redis.Int(c.Do("GET", name))
	if err != nil || weight == 0 {
		weight = rand.Intn(50) // figure out a good starting weight for a new heuristic
	}
	return weight
}

func (snake *HeuristicSnake) Move(request *MoveRequest) string {
	var heuristicWaitGroup sync.WaitGroup
	heuristicWaitGroup.Add(len(snake.weightedHeuristics))

	// do heuristics
	for _, weightedHeuristic := range snake.weightedHeuristics {
		go func(h *WeightedHeuristic) {
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

	return bestDirection
}
