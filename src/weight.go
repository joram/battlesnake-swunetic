package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"sort"
	"time"
)

var weightCache = make(map[string]int)

func LogBestWeights(bestWeights map[string]int, numGames int, duration time.Duration, quality float64, averageTurns int) {
	keys := []string{}
	for key := range bestWeights {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	s := ""
	for _, key := range keys {
		s += fmt.Sprintf("sample#heuristic.%v=%v ", key, int(bestWeights[key]))
	}
	s += fmt.Sprintf(
		"sample#games.Count=%v sample#games.ElapsedTime=%v sample#games.Quality=%v, sample#games.AverageTurns=%v",
		numGames,
		duration,
		quality*float64(100),
		averageTurns,
	)
	println(s)
}

func SnakeQualities(games []*Game) map[string]float64 {
	snakeWins := make(map[string]float64)
	for _, snake := range games[0].currentGameState.Snakes {
		snakeWins[snake.Id] = 0
	}
	for _, game := range games {
		for _, snake := range game.currentGameState.winners {
			snakeWins[snake.GetId()] += 1
		}
	}

	return snakeWins
}

func PrimeWeightsCache() {
	for name, _ := range heuristics {
		weight := getWeight(name)
		fmt.Printf("%v\t%v\n", name, weight)
	}
}

func getWeight(name string) int {
	val, set := weightCache[name]
	if set {
		return val
	}

	c := redisConnectionPool.Get()
	defer c.Close()

	weight, err := redis.Int(c.Do("GET", name))
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
	return weight
}

func StoreWeights(weights map[string]int) {
	c := redisConnectionPool.Get()
	defer c.Close()

	average := 0
	for _, w := range weights {
		average += w
	}
	average = int(average / len(weights))
	offset := 50 - average

	for name, weight := range weights {
		c.Do("SET", name, weight+offset)
		weightCache[name] = weight
	}
}
