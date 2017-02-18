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
	Food   [][]int            `json:"food"`
	GameId string             `json:"game_id"`
	Height int                `json:"height"`
	Width  int                `json:"width"`
	Turn   int                `json:"turn"`
	Snakes []MoveRequestSnake `json:"snakes"`
	You    string             `json:"you"`
}

type MoveRequestSnake struct {
	Coords       [][]int `json:"coords"`
	HealthPoints int     `json:"health_points"`
	Id           string  `json:"id"`
	Name         string  `json:"name"`
	Taunt        string  `json:"taunt"`
}

type MoveResponse struct {
	Move  string  `json:"move"`
	Taunt *string `json:"taunt,omitempty"`
}

type Game struct {
	currentGameState *GameState
	foodFrequency    int
}

type GameState struct {
	HeuristicSnakes []HeuristicSnake
	Snakes          []Snake
	Height          int
	Width           int
	Turn            int
	Food            []Point
	state           string
	winners         []HeuristicSnake
	losers          []HeuristicSnake
	You             string
}

type MoveHeuristic func(gameState *GameState) WeightedDirections

type WeightedHeuristic struct {
	Weight             int
	WeightedDirections WeightedDirections
	Name               string
	MoveFunc           MoveHeuristic
}

type HeuristicSnake struct {
	Id                 string
	WeightedHeuristics []*WeightedHeuristic
	DiedOnTurn         int
}

type Vector Point
type Point struct {
	X int
	Y int
}

type Snake struct {
	Coords       []Point
	HealthPoints int
	Id           string
	Name         string
	Taunt        string
}

type WeightedDirection struct {
	Direction string
	Weight    int
}

type WeightedDirections []WeightedDirection
