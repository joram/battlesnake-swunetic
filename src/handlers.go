package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func start(w http.ResponseWriter, r *http.Request) {
	var requestData GameStartRequest
	json.NewDecoder(r.Body).Decode(&requestData)

	log.Printf("Game starting - %v\n", requestData.GameId)
	responseData := GameStartResponse{
		Color:   "#F7931D",
		Name:    "SWU Bounty Snake",
		HeadUrl: stringPtr("https://s3.amazonaws.com/john-box-o-mysteries/swu-logo.png"),
	}
	b, err := json.Marshal(responseData)
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
	w.Write(b)
}

func move(w http.ResponseWriter, r *http.Request) {
	var requestData MoveRequest
	json.NewDecoder(r.Body).Decode(&requestData)
	snake := NewHeuristicSnake(requestData.GameId)
	gameState := NewGameState(requestData)
	responseData := MoveResponse{
		Move: snake.Move(&gameState),
	}
	b, err := json.Marshal(responseData)
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
	w.Write(b)
}
