package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"math/rand"
	"sort"
)

func Train(numSnakes, numGamesPerGeneration int) {
	game := NewGame("", numSnakes, 1)
	trainingSnakes := game.currentGameState.HeuristicSnakes
	MutateSnakes(trainingSnakes)
	gamesWon := map[string]int{}
	for _, snake := range trainingSnakes {
		gamesWon[snake.Id] = 0
	}

	for i := 0; i < numGamesPerGeneration; i++ {
		game := NewGame(fmt.Sprintf("Game-%v", i), len(trainingSnakes), 1)
		game.currentGameState.HeuristicSnakes = trainingSnakes
		game.Run()
		//game.Print()

		for _, winner := range game.currentGameState.winners {
			gamesWon[winner.Id] += 1
		}
	}

	bestWeights := BestWeights(gamesWon, trainingSnakes)
	StoreWeights(bestWeights)

	keys := []string{}
	for key, _ := range bestWeights {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	s := "NEW BEST: "
	for _, key := range keys {
		s += fmt.Sprintf("%v:%v, ", key, bestWeights[key])
	}
	println(s)
}

func BestWeights(gamesWon map[string]int, snakes []*HeuristicSnake) map[string]int {

	bestSnakeWins := 0
	var bestSnake *HeuristicSnake
	for _, snake := range snakes {
		if gamesWon[snake.Id] >= bestSnakeWins {
			bestSnakeWins = gamesWon[snake.Id]
			bestSnake = snake
		}
	}

	weights := make(map[string]int)
	for _, weightedHeuristic := range bestSnake.WeightedHeuristics {
		weights[weightedHeuristic.Name] = weightedHeuristic.Weight
	}

	return weights
}

func getWeight(name string) int {
	c := redisConnectionPool.Get()
	defer c.Close()

	weight, err := redis.Int(c.Do("GET", name))
	if err != nil || weight == 0 {
		weight = rand.Intn(50) // figure out a good starting Weight for a new heuristic
	}
	return weight
}

func StoreWeights(weights map[string]int) {
	c := redisConnectionPool.Get()
	defer c.Close()

	for name, weight := range weights {
		c.Do("SET", name, weight)
	}
}

func MutateSnakes(snakes []*HeuristicSnake) {

	mutationAmount := []int{0, 2, 3, 20}

	for i := 0; i < len(snakes); i++ {
		snakes[i].Mutate(mutationAmount[i])
		weights := map[string]int{}
		for _, weight := range snakes[i].WeightedHeuristics {
			w := weight.Weight
			if w > 100 {
				w = 100
			}
			if w < 0 {
				w = 0
			}
			weights[weight.Name] = w
		}
	}
}
