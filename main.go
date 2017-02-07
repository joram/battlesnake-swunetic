package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/start", start)
	http.HandleFunc("/move", move)
	http.ListenAndServe(":8000", nil)
}
