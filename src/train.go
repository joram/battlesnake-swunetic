package main

import (
	"fmt"
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

		winPercent := heuristicQuality / float64(len(games)) * 100
		LogBestWeights(bestWeights, numGamesPerGeneration, time.Since(start), winPercent, averageTurns)
		fmt.Printf("\n\t%.2f%% wins ", winPercent)
	} else {
		print(".")
	}
	return heuristicQuality
}

func RunGames(snakeAIs []SnakeAI, snakeNames []string, numGamesPerGeneration int) []*Game {
	doneGamesChan := make(chan *Game)
	gamesChan := make(chan *Game)

	games := []*Game{}
	go func() {
		for game := range doneGamesChan {
			games = append(games, game)
			if len(games) >= numGamesPerGeneration {
				close(gamesChan)
				close(doneGamesChan)
			}
		}
	}()

	// run 10 games in parallel
	workersCount := 10
	for i := 0; i < workersCount; i++ {
		go func() {
			for game := range gamesChan {
				game.Run()
				doneGamesChan <- game
			}
		}()
	}

	// add games
	for i := 0; i < numGamesPerGeneration; i++ {
		game := NewGame(fmt.Sprintf("Game-%v", i), snakeNames, 1)
		game.currentGameState.SnakeAIs = snakeAIs
		gamesChan <- game
	}

	return games
}
