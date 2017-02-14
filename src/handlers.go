package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func start(w http.ResponseWriter, r *http.Request) {
	var requestData GameStartRequest
	json.NewDecoder(r.Body).Decode(&requestData)
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	log.Printf("Game starting - %v\n", requestData.GameId)
	responseData := GameStartResponse{
		Color:   "#F7931D",
		Name:    "SWU Bounty Snake",
		HeadUrl: stringPtr(fmt.Sprintf("%v://%v/swu-logo.png", scheme, r.Host)),
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
	log.Printf("Move request - direction:%v\n", responseData.Move)
	b, err := json.Marshal(responseData)
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
	w.Write(b)
}
