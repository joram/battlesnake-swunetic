package main

import (
	"flag"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/sendwithus/lib-go"
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

	bestWeights := BestWeights(gamesWon, trainingSnakes)
	StoreWeights(bestWeights)
	fmt.Printf("NEW BEST: %v", bestWeights)
}

func BestWeights(gamesWon map[string]int, snakes []HeuristicSnake) map[string]int {
	weights := make(map[string]int)
	for _, weightedHeuristic := range snakes[0].WeightedHeuristics {
		// TODO: DO this properly, not just picking the first
		weights[weightedHeuristic.Name] = weightedHeuristic.Weight
	}
	return weights
}

func StoreWeights(weights map[string]int) {
	c, err := redis.Dial("tcp", swu.GetEnvVariable("REDIS_URL", true))
	if err != nil {
		panic(err)
	}
	defer c.Close()

	for name, weight := range weights {
		c.Do("SET", name, weight)
	}
}

func MutateSnakes(snakes []HeuristicSnake) []HeuristicSnake {

	mutationAmount := []int{0, 5, 10, 15, 20, 20, 20, 20, 20, 20, 20, 20}

	for i := 0; i < len(snakes); i++ {
		snakes[i].Mutate(mutationAmount[i])
	}
	return snakes
}
