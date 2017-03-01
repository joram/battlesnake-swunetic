package main

import (
	"fmt"
	"sync"
	"time"
	_ "time"
)

func TrainAgainstSnek(numGamesPerGeneration int, mutation int, bestQualitySoFar float64) float64 {
	start := time.Now()
	heuristicSnakeId := "MutatedSnake"
	snake := NewHeuristicSnake(heuristicSnakeId)
	bestWeights := snake.GetWeights()
	snake.Mutate(mutation)
	snek := NewSnekSnake()
	snakeAIs := []SnakeAI{snake, snek}
	snakeNames := []string{heuristicSnakeId, "snek"}

	averageTurns := -1
	games := RunGames(snakeAIs, snakeNames, numGamesPerGeneration)
	qualities := SnakeQualities(games)
	heuristicQuality := qualities[heuristicSnakeId]
	if heuristicQuality > bestQualitySoFar {
		bestWeights = snake.GetWeights()
		StoreWeights(bestWeights)

		averageTurns = 0
		for _, game := range games {
			averageTurns += game.currentGameState.Turn
		}
		averageTurns = averageTurns / len(games)
		LogBestWeights(bestWeights, numGamesPerGeneration, time.Since(start), heuristicQuality, averageTurns)
		fmt.Printf("\n\t%.2f%% wins ", heuristicQuality/float64(len(games))*100)
	} else {
		print(".")
	}
	return heuristicQuality
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
