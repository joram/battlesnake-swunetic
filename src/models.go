package main

type GameStartRequest struct {
	GameId string `json:"game_id"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

type GameStartResponse struct {
	Color   string  `json:"color"`
	HeadUrl *string `json:"head_url,omitempty"`
	Name    string  `json:"name"`
	Taunt   *string `json:"taunt,omitempty"`
}

type MoveRequest struct {
	Food   []Point `json:"food"`
	GameId string  `json:"game_id"`
	Height int     `json:"height"`
	Width  int     `json:"width"`
	Turn   int     `json:"turn"`
	Snakes []Snake `json:"snakes"`
	You    string  `json:"you"`
}

type MoveResponse struct {
	Move  string  `json:"Move"`
	Taunt *string `json:"taunt,omitempty"`
}

type Game struct {
	currentGameState *GameState
	foodFrequency    int
}

type GameState struct {
	HeuristicSnakes []HeuristicSnake
	Snakes          []Snake `json:"snakes"`
	Height          int     `json:"height"`
	Width           int     `json:"width"`
	Turn            int     `json:"turn"`
	Food            []Point `json:"food"`
	state           string
	winners         []HeuristicSnake
	losers          []HeuristicSnake
	You             string
}

type MoveHeuristic func(gameState *GameState) string

type WeightedHeuristic struct {
	Weight   int
	Move     string
	Name     string
	MoveFunc MoveHeuristic
}

type HeuristicSnake struct {
	Id                 string
	WeightedHeuristics []*WeightedHeuristic
	DiedOnTurn         int
}

type BoardCell struct {
	State string  `json:"state"`
	Snake *string `json:"snake,omitempty"`
}

type Vector Point
type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Snake struct {
	Coords       []Point `json:"coords"`
	HealthPoints int     `json:"health_points"`
	Id           string  `json:"id"`
	Name         string  `json:"name"`
	Taunt        string  `json:"taunt"`
}

type WeightedDirection struct {
	Direction string
	Weight    int
}

type WeightedDirections []WeightedDirection
