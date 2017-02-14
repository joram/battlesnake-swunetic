package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {

	simulate := flag.Bool("sim", false, "simulate instead of starting a web snake")
	flag.Parse()

	if !*simulate {
		fs := http.FileServer(http.Dir("assets"))
		http.Handle("/", fs)

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
		game := NewGame(4)
		game.Run()
	}
}
