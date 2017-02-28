package main

import (
	"encoding/json"
	"github.com/sendwithus/lib-go"
	"log"
	"net/http"
)

func start(w http.ResponseWriter, r *http.Request) {
	var requestData GameStartRequest
	json.NewDecoder(r.Body).Decode(&requestData)

	log.Printf("Game starting - %v\n", requestData.GameId)
	responseData := GameStartResponse{
		Color:   "#00FF00",
		Name:    "dsnek",
		HeadUrl: swu.String("https://s3.amazonaws.com/john-box-o-mysteries/swu-logo.png"),
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
	responseData := MoveResponse{
		Move: requestData.GenerateMove(),
	}
	log.Printf("Move request - direction:%v\n", responseData.Move)
	b, err := json.Marshal(responseData)
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
	w.Write(b)
}
