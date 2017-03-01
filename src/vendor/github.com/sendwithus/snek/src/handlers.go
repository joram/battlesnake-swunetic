package snek

import (
	"encoding/json"
	"fmt"
	"github.com/sendwithus/lib-go"
	"io/ioutil"
	"log"
	"net/http"
)

func start(w http.ResponseWriter, r *http.Request) {
	var requestData GameStartRequest
	json.NewDecoder(r.Body).Decode(&requestData)

	log.Printf("Game starting - %v\n", requestData.GameId)
	responseData := GameStartResponse{
		Color:   "#35AA47",
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
	val, err := ioutil.ReadAll(r.Body)
	json.Unmarshal(val, &requestData)
	responseData := MoveResponse{
		Move: requestData.GenerateMove(),
	}
	log.Printf("Move request - direction:%v\n", responseData.Move)
	if err != nil {
		fmt.Printf("ERR: %#v\n", err)
	}
	log.Printf("%v\n", string(val))
	b, err := json.Marshal(responseData)
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
	w.Write(b)
}
