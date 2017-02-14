package main

import (
	"github.com/garyburd/redigo/redis"
	"math/rand"
	"sync"
)


func (weightedHeuristic *WeightedHeuristic) Calculate(gameState *GameState) {
	weightedHeuristic.move = weightedHeuristic.moveHeuristic(gameState)
}

func NewHeuristicSnake(id string) HeuristicSnake {
	snake := HeuristicSnake{
		Id: id,
		WeightedHeuristics: []WeightedHeuristic{},
	}

	heuristics := map[string]MoveHeuristic{
	//	 this is where we list all heuristics we've written
	}
	for name, heuristic := range heuristics {
		snake.WeightedHeuristics = append(snake.WeightedHeuristics, WeightedHeuristic{
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

func (snake *HeuristicSnake) Move(gameState *GameState) string {
	var heuristicWaitGroup sync.WaitGroup
	heuristicWaitGroup.Add(len(snake.WeightedHeuristics))

	// do heuristics
	for _, weightedHeuristic := range snake.WeightedHeuristics {
		go func(h *WeightedHeuristic) {
			h.Calculate(gameState)
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
	for _, weightedHeuristic := range snake.WeightedHeuristics {
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
