package main

import (
	"fmt"
	"log"
	"time"
)

func Train() {
	log.Println("Simulate a game to train swunetics!")
	numWorkers := 2
	numGames := 200
	numFood := 6
	mutation := 2
	bestQuality := TrainAgainstSnek(numGames, 0, numFood, numWorkers, 0)
	fmt.Printf("\nstarting quality: %v\n", bestQuality)
	for {
		quality := TrainAgainstSnek(numGames, mutation, numFood, numWorkers, bestQuality)
		if quality > bestQuality {
			bestQuality = quality
		}
	}
}

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
	}
	
	averageTurns = 0
	for _, game := range games {
		averageTurns += game.currentGameState.Turn
	}
	averageTurns = averageTurns / len(games)

	winPercent := WinPercent(games, heuristicSnakeId)
	LogBestWeights(bestWeights, numGamesPerGeneration, time.Since(start), winPercent, averageTurns)
	fmt.Printf("\n\t%.2f%% wins\n", winPercent)

	fmt.Printf("sample#quality=%v\n", heuristicQuality)
	return heuristicQuality
}

func WinPercent(games []Game, snakeId string) float64 {
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

func SnakeQualities(games []Game) map[string]float64 {

	snakeWins := make(map[string]float64)
	for _, snake := range games[0].currentGameState.SnakeAIs {
		snakeWins[snake.GetId()] = WinPercent(games, snake.GetId())
	}
	fmt.Printf("%v\n", snakeWins)
	return snakeWins
}

func RunGames(snakeAIs []SnakeAI, snakeNames []string, numGamesPerGeneration, amountOfFood, workerCount int) []Game {
	doneGamesChan := make(chan *Game)
	gamesChan := make(chan *Game)

	games := []Game{}
	go func() {
		for game := range doneGamesChan {
			games = append(games, *game)
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
		game := NewGame(fmt.Sprintf("Game-%v", i), snakeNames, amountOfFood, 15, 15)
		game.currentGameState.SnakeAIs = snakeAIs
		gamesChan <- game
	}

	for len(games) < numGamesPerGeneration {
		time.Sleep(time.Millisecond)
	}

	return games
}
