package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {

	simulate := flag.Bool("sim", false, "simulate instead of starting a web snake")
	flag.Parse()

	if !*simulate {
		http.HandleFunc("/start", start)
		http.HandleFunc("/move", move)
		port := os.Getenv("PORT")
		if port == "" {
			port = "9000"
		}

		log.Printf("Running server on port %s...\n", port)
		http.ListenAndServe(":"+port, nil)
	} else {
		log.Println("Simulate a game to train swunetics!")
		Train(4, 100)

	}
}

func Train(numSnakes, numGamesPerGeneration int) {
	game := NewGame(numSnakes, 1)
	trainingSnakes := game.currentGameState.HeuristicSnakes
	trainingSnakes = MutateSnakes(trainingSnakes)
	gamesWon := map[string]int{}
	for _, snake := range trainingSnakes {
		gamesWon[snake.Id] = 0
	}

	for i := 0; i < numGamesPerGeneration; i++ {
		game := NewGame(len(trainingSnakes), 1)
		game.currentGameState.HeuristicSnakes = trainingSnakes
		game.Run()
		for _, winner := range game.currentGameState.winners {
			gamesWon[winner.Id] += 1
		}

		s := ""
		for _, snake := range trainingSnakes {
			s += fmt.Sprint("%v:%v\t", snake.Id, gamesWon[snake.Id])
		}
	}
}

func MutateSnakes(snakes []HeuristicSnake) []HeuristicSnake {
	mutationAmount := []int{0, 5, 10, 15, 20, 20, 20, 20, 20, 20, 20, 20}

	for i := 0; i < len(snakes); i++ {
		snakes[i].Mutate(mutationAmount[i])
	}
	return snakes
}
