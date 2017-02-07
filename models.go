package main

type GameStartRequest struct {
	GameId string `json:"game_id"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

type GameStartResponse struct {
	Color   string `json:"color"`
	HeadUrl string `json:"head_url"`
	Name    string `json:"name"`
	Taunt   string `json:"taunt"`
}

type MoveRequest struct {
	Board  [][]BoardCell `json:"board"`
	Food   []Point       `json:"food"`
	GameId string        `json:"game_id"`
	Height int           `json:"height"`
	Width  int           `json:"width"`
	Turn   int           `json:"turn"`
	Snakes []Snake       `json:"snakes"`
	You    string        `json:"you"`
}

type MoveResponse struct {
}

type BoardCell struct {
	State string `json:"state"`
	Snake string `json:"snake"`
}

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Snake struct {
}
