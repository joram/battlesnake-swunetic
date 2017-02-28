package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"sort"
	"sync"
)

var weightCache = make(map[string]float64)

func Train(numSnakes, numGamesPerGeneration int) {
	games := RunGames(numSnakes, numGamesPerGeneration)
	bestWeights := BestWeights(games)
	StoreWeights(bestWeights)
	LogBestWeights(bestWeights)
}

func RunGames(numSnakes int, numGamesPerGeneration int) []*Game {
	games := []*Game{}

	snakes := NewGame("", numSnakes, 1).currentGameState.HeuristicSnakes
	for i, snake := range snakes {
		snake.Mutate(i * 2)
	}

	wg := sync.WaitGroup{}
	wg.Add(numGamesPerGeneration)
	for i := 0; i < numGamesPerGeneration; i++ {
		game := NewGame(fmt.Sprintf("Game-%v", i), numSnakes, 1)
		for j, snake := range game.currentGameState.HeuristicSnakes {
			snake.WeightedHeuristics = snakes[j].WeightedHeuristics
		}
		games = append(games, game)
		go func(game *Game, wg *sync.WaitGroup) {
			game.Run()
			wg.Done()
		}(game, &wg)
	}
	wg.Wait()
	return games
}

func LogBestWeights(bestWeights map[string]float64) {
	keys := []string{}
	for key := range bestWeights {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	s := "Weights: "
	for _, key := range keys {
		s += fmt.Sprintf("sample#heuristic.%v=%v ", key, int(bestWeights[key]))
	}
	println(s)
}

func BestWeights(games []*Game) map[string]float64 {
	winningBonusWeight := 500

	snakeQuality := make(map[string]float64)
	for _, game := range games {
		for _, snake := range game.currentGameState.winners {
			snakeQuality[snake.Id] += float64(snake.DiedOnTurn + winningBonusWeight)
		}
		for _, snake := range game.currentGameState.losers {
			snakeQuality[snake.Id] += float64(snake.DiedOnTurn)
		}
	}

	totalQuality := float64(0)
	for _, quality := range snakeQuality {
		totalQuality += quality
	}

	for snakeId, quality := range snakeQuality {
		snakeQuality[snakeId] = quality / totalQuality
	}

	weights := make(map[string]float64)
	for _, snake := range games[0].currentGameState.HeuristicSnakes {
		snakeWeights := make(map[string]int)
		for _, weightedHeuristic := range snake.WeightedHeuristics {
			snakeWeights[weightedHeuristic.Name] = weightedHeuristic.Weight
		}
		for name, weight := range snakeWeights {
			weights[name] += float64(weight) * snakeQuality[snake.Id]
		}
	}

	return weights
}

func getWeight(name string) float64 {
	val, set := weightCache[name]
	if set {
		return val
	}

	c := redisConnectionPool.Get()
	defer c.Close()

	weight, err := redis.Float64(c.Do("GET", name))
	if err != nil {
		println(err.Error())
		weight = 50
	}
	if weight < 0 {
		weight = 0
	}
	if weight > 100 {
		weight = 100
	}

	weightCache[name] = weight
	println(name, ": ", int(weight))
	return weight
}

func StoreWeights(weights map[string]float64) {
	c := redisConnectionPool.Get()
	defer c.Close()

	average := float64(0)
	for _, w := range weights {
		average += w
	}
	average = average / float64(len(weights))
	offset := float64(50) - average

	for name, weight := range weights {
		c.Do("SET", name, weight+offset)
		weightCache[name] = weight
	}
}
