package main

import (
	"fmt"
	"sync"
	"time"
)

func TrainAgainstSnek(numGamesPerGeneration int, bestQualitySoFar float64) {
	start := time.Now()
	heuristicSnakeId := "MutatedSnake"
	snake := NewHeuristicSnake(heuristicSnakeId)
	snake.Mutate(3)
	snek := NewSnekSnake()
	snakeAIs := []SnakeAI{snake, snek}
	snakeNames := []string{heuristicSnakeId, "snek"}

	games := RunGames(snakeAIs, snakeNames, numGamesPerGeneration)
	qualities := SnakeQualities(games)
	heuristicQuality := qualities[heuristicSnakeId]
	bestWeights := map[string]int{}
	if heuristicQuality > bestQualitySoFar {
		bestWeights = snake.GetWeights()
		StoreWeights(bestWeights)
	} else {
		bestWeights = map[string]int{}
	}
	LogBestWeights(bestWeights, numGamesPerGeneration, time.Since(start))
}

func RunGames(snakeAIs []SnakeAI, snakeNames []string, numGamesPerGeneration int) []*Game {
	games := []*Game{}

	wg := sync.WaitGroup{}
	wg.Add(numGamesPerGeneration)
	for i := 0; i < numGamesPerGeneration; i++ {
		game := NewGame(fmt.Sprintf("Game-%v", i), snakeNames, 1)
		game.currentGameState.SnakeAIs = snakeAIs
		games = append(games, game)
		go func(game *Game, wg *sync.WaitGroup) {
			game.Run()
			wg.Done()
		}(game, &wg)
	}
	wg.Wait()
	return games
}
