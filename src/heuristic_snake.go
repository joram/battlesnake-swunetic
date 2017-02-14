package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"math/rand"
	"sort"
	"sync"
	"time"
)

func (weightedHeuristic *WeightedHeuristic) Calculate(gameState *GameState) {
	weightedHeuristic.move = weightedHeuristic.moveHeuristic(gameState)
}

func NewHeuristicSnake(id string) HeuristicSnake {
	snake := HeuristicSnake{
		Id:                 id,
		WeightedHeuristics: []WeightedHeuristic{},
	}

	heuristics := map[string]MoveHeuristic{
		"straight": GoStraightHeuristic,
	}

	for name, heuristic := range heuristics {
		snake.WeightedHeuristics = append(snake.WeightedHeuristics, WeightedHeuristic{
			weight:        getWeight(name),
			moveHeuristic: heuristic,
			Name:          name,
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
		UP:    0,
		DOWN:  0,
		LEFT:  0,
		RIGHT: 0,
		NOOP:  0,
	}
	for _, weightedHeuristic := range snake.WeightedHeuristics {
		weights[weightedHeuristic.move] += weightedHeuristic.weight
	}
	weights[NOOP] = 0
	fmt.Printf("weights: UP:%v, DOWN:%v, LEFT:%v, RIGHT:%V", weights[UP], weights[DOWN], weights[LEFT], weights[RIGHT])

	weightedDirections := WeightedDirections{
		WeightedDirection{Direction: UP, Weight: weights[UP]},
		WeightedDirection{Direction: DOWN, Weight: weights[DOWN]},
		WeightedDirection{Direction: LEFT, Weight: weights[LEFT]},
		WeightedDirection{Direction: RIGHT, Weight: weights[RIGHT]},
	}

	sort.Sort(weightedDirections)
	for _, weightedDirection := range weightedDirections {
		head := gameState.MySnake().Coords[0]
		directionOfMovement := directionVector(weightedDirection.Direction)
		possibleNewHead := head.Add(directionOfMovement)
		if !gameState.IsSolid(possibleNewHead) {
			return weightedDirection.Direction
		}
	}

	return NOOP
}
