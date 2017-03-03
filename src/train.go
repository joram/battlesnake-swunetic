package main

import (
	"fmt"
	"time"
	_ "time"
)

func TrainAgainstSnek(numGamesPerGeneration, mutation, amountOfFood, workerCount int, bestQualitySoFar float64) float64 {
	start := time.Now()
	heuristicSnakeId := "MutatedSnake"
	snake := NewHeuristicSnake(heuristicSnakeId)
	bestWeights := snake.GetWeights()
	snake.Mutate(mutation)
	snek := NewSnekSnake()
	snakeAIs := []SnakeAI{snake, snek}
	snakeNames := []string{heuristicSnakeId, "snek"}

	averageTurns := -1
	games := RunGames(snakeAIs, snakeNames, numGamesPerGeneration, amountOfFood, workerCount)
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

		winPercent := WinPercent(games, heuristicSnakeId)
		LogBestWeights(bestWeights, numGamesPerGeneration, time.Since(start), winPercent, averageTurns)
		fmt.Printf("\n\t%.2f%% wins\n", winPercent)
	}
	fmt.Printf("sample#quality=%v\n", heuristicQuality)
	return heuristicQuality
}

func WinPercent(games []*Game, snakeId string) float64 {
	count := float64(0)
	for _, game := range games {
		for _, winner := range game.currentGameState.winners {
			if winner.GetId() == snakeId {
				count += 1
			}
		}
	}

	return count * float64(100) / float64(len(games))
}

func RunGames(snakeAIs []SnakeAI, snakeNames []string, numGamesPerGeneration, amountOfFood, workerCount int) []*Game {
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
	for i := 0; i < workerCount; i++ {
		go func() {
			for game := range gamesChan {
				game.Run()
				doneGamesChan <- game
			}
		}()
	}

	// add games
	for i := 0; i < numGamesPerGeneration; i++ {
		game := NewGame(fmt.Sprintf("Game-%v", i), snakeNames, amountOfFood)
		game.currentGameState.SnakeAIs = snakeAIs
		gamesChan <- game
	}

	for len(games) < numGamesPerGeneration {
		time.Sleep(time.Millisecond)
	}

	return games
}
