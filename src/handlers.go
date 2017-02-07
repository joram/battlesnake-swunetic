package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func start(w http.ResponseWriter, r *http.Request) {
	var requestData GameStartRequest
	json.NewDecoder(r.Body).Decode(&requestData)
	responseData := GameStartResponse{
		Color: "#00FF00",
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
		Move: "up",
	}
	b, err := json.Marshal(responseData)
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
	w.Write(b)
}
