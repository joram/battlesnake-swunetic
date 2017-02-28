package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"sort"
	"time"
)

var weightCache = make(map[string]int)

func LogBestWeights(bestWeights map[string]int, numGames int, duration time.Duration) {
	keys := []string{}
	for key := range bestWeights {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	s := "Weights: "
	for _, key := range keys {
		s += fmt.Sprintf("sample#heuristic.%v=%v ", key, int(bestWeights[key]))
	}
	s += fmt.Sprintf("sample#games.Count=%v sample#games.ElapsedTime=%v", numGames, duration)
	println(s)
}

func SnakeQualities(games []*Game) map[string]float64 {
	winningBonusWeight := 500

	snakeQuality := make(map[string]float64)
	for _, game := range games {
		for _, snake := range game.currentGameState.winners {
			snakeQuality[snake.GetId()] += float64(snake.GetDiedOnTurn() + winningBonusWeight)
		}
		for _, snake := range game.currentGameState.losers {
			snakeQuality[snake.GetId()] += float64(snake.GetDiedOnTurn())
		}
	}

	totalQuality := float64(0)
	for _, quality := range snakeQuality {
		totalQuality += quality
	}

	for snakeId, quality := range snakeQuality {
		snakeQuality[snakeId] = quality / totalQuality
	}

	return snakeQuality
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
