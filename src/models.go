package main

import (
	"sync"
	"time"
)

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
	duration         time.Duration
	name             string
}

type GameState struct {
	SnakeAIs   []SnakeAI
	all        []SnakeAI
	winners    []SnakeAI
	losers     []SnakeAI
	Snakes     []*Snake
	Height     int
	Width      int
	Turn       int
	Food       []Point
	state      string
	You        string
	aStar      map[string]*AStar
	DiedOnTurn map[string]int
}

type MoveHeuristic func(gameState *GameState) WeightedDirections

type WeightedHeuristic struct {
	Weight             int
	WeightedDirections WeightedDirections
	Name               string
	MoveFunc           MoveHeuristic `json:"-"`
}

type SnakeAI interface {
	Move(gameState *GameState) string
	SetDiedOnTurn(int)
	GetDiedOnTurn() int
	GetWeights() map[string]int
	GetId() string
	Mutate(int)
}

type HeuristicSnake struct {
	Id                 string
	WeightedHeuristics []*WeightedHeuristic
	DiedOnTurn         int
}

type SnekSnake struct {
	DiedOnTurn int
}

type Points []*Point
type Vector Point
type Point struct {
	X int
	Y int
}
type PointPair struct {
	from *Point
	to   *Point
}

type PathCalculation struct {
	start         *Point
	goals         Points
	achievedGoals Points
	visited       []PointPair
	toVisit       Points
	gameState     *GameState
}

type AStar struct {
	gameState       *GameState
	start           *Point
	turnsTo         map[Point]int
	visited         map[Point]bool
	pathToCache     map[Point][]*Point
	pathToCacheLock sync.Mutex
	canVisitCount   int
}

type AStarPoint struct {
	point   *Point
	turnsTo int
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
