package main

import (
	"encoding/json"
	"log"
	"math/rand"
)

func NewGameState(mr MoveRequest) GameState {

	snakes := []*Snake{}
	for _, snake := range mr.Snakes {
		body := []Point{}
		for _, coord := range snake.Coords {
			part := Point{X: coord[0], Y: coord[1]}
			body = append(body, part)
		}
		snakes = append(snakes, &Snake{
			Coords:       body,
			HealthPoints: snake.HealthPoints,
			Id:           snake.Id,
			Name:         snake.Name,
			Taunt:        snake.Taunt,
		})
	}

	foods := []Point{}
	for _, coord := range mr.Food {
		food := Point{X: coord[0], Y: coord[1]}
		foods = append(foods, food)

	}

	gameState := GameState{
		Snakes:     snakes,
		Width:      mr.Width,
		Height:     mr.Height,
		Turn:       mr.Turn,
		Food:       foods,
		winners:    []SnakeAI{},
		state:      "running",
		You:        mr.You,
		DiedOnTurn: map[string]int{},
	}

	gameState.aStar = map[string]AStar{}
	for _, snake := range gameState.Snakes {
		gameState.aStar[snake.Id] = NewAStar(&gameState, snake.Head())
	}
	return gameState
}

func CloneGameState(gs *GameState) GameState {
	snakes := []*Snake{}
	for _, s := range gs.Snakes {
		clonedSnake := Snake{
			s.Coords,
			s.HealthPoints,
			s.Id,
			s.Name,
			s.Taunt,
		}
		snakes = append(snakes, &clonedSnake)
	}
	return GameState{
		gs.SnakeAIs,
		gs.all,
		gs.winners,
		gs.losers,
		snakes,
		gs.Height,
		gs.Width,
		gs.Turn,
		gs.Food,
		gs.state,
		gs.You,
		map[string]AStar{},
		gs.DiedOnTurn,
	}

}

func (gameState *GameState) MySnake() *Snake {
	snake := gameState.GetSnake(gameState.You)
	return snake
}

func (gameState *GameState) OtherSnakes() []*Snake {
	snakes := []*Snake{}
	for _, snake := range gameState.Snakes {
		if snake.Id != gameState.You {
			snakes = append(snakes, snake)
		}
	}
	return snakes
}

func (gameState *GameState) CountSurroundingWalls(p *Point) int {
	solids := 0
	for _, neighbour := range p.NeighboursWithDiagonals() {
		if gameState.IsSolid(neighbour, "") {
			solids += 1
		}
	}
	return solids
}

func (gameState *GameState) String() string {
	b, err := json.Marshal(gameState)
	if err != nil {
		log.Fatalf("%v", err)
		panic("cant string")
	}
	return string(b)
}

func (gameState *GameState) NextGameState() *GameState {
	nextGameState := GameState{
		Turn:       gameState.Turn + 1,
		SnakeAIs:   gameState.SnakeAIs,
		Snakes:     gameState.Snakes,
		Width:      gameState.Width,
		Height:     gameState.Height,
		Food:       gameState.Food,
		winners:    gameState.winners,
		state:      gameState.state,
		You:        gameState.You,
		DiedOnTurn: gameState.DiedOnTurn,
	}

	// get all moves
	moveDirections := map[string]string{}
	for i, snakeAI := range gameState.SnakeAIs {
		snakeId := gameState.Snakes[i].Id
		gameState.You = snakeId
		moveDirections[snakeId] = snakeAI.Move(gameState)
	}

	// extend all snakes
	newHeads := map[string]Point{}
	for _, snake := range nextGameState.Snakes {
		direction := moveDirections[snake.Id]
		newHead := snake.Extend(direction)
		newHeads[snake.Id] = newHead
	}

	// eat or shrink
	foodEaten := []Point{}
	for snakeId, newHead := range newHeads {
		snake := nextGameState.GetSnake(snakeId)
		if nextGameState.FoodAt(&newHead) {
			foodEaten = append(foodEaten, newHead)
			snake.HealthPoints = 100
		} else {
			snake.HealthPoints -= 1
			snake.Coords = snake.Coords[0 : len(snake.Coords)-1]
		}
		if snake.HealthPoints <= 0 {
			nextGameState.KillSnake(snakeId)
		}
	}

	nextGameState.UpdateFood(foodEaten)

	// collision
	for snakeId, newHead := range newHeads {
		if gameState.IsSolid(&newHead, snakeId) {
			nextGameState.KillSnake(snakeId)
		}
	}

	nextGameState.CheckWinStates()

	nextGameState.aStar = map[string]AStar{}
	for _, snake := range nextGameState.Snakes {
		nextGameState.aStar[snake.Id] = NewAStar(&nextGameState, snake.Head())
	}

	return &nextGameState
}

func (gameState *GameState) UpdateFood(foodEaten []Point) {

	// remove food
	for _, eatenFood := range foodEaten {
		newFoodList := []Point{}
		for _, food := range gameState.Food {
			if !eatenFood.Equals(food) {
				newFoodList = append(newFoodList, food)
			}
		}
		gameState.Food = newFoodList
	}

	// spawn food
	for i := 0; i < len(foodEaten); i++ {
		gameState.SpawnFood()
	}
}

func (gameState *GameState) CheckWinStates() {
	numSnakes := len(gameState.Snakes)
	if numSnakes == 1 {
		gameState.state = "Won"
		snake := gameState.Snakes[0]
		gameState.MarkSnakeAsWinner(snake.Id)
	}
	if numSnakes == 0 {
		gameState.state = "Draw"
		gameState.winners = []SnakeAI{}
		for _, snake := range gameState.Snakes {
			gameState.KillSnake(snake.Id)
		}
	}
}

func (gameState *GameState) SpawnFood() {
	emptyPoints := []Point{}
	for x := 0; x < gameState.Width; x += 1 {
		for y := 0; y < gameState.Height; y += 1 {
			p := Point{X: x, Y: y}
			if !gameState.IsSolid(&p, "") && !gameState.FoodAt(&p) {
				emptyPoints = append(emptyPoints, p)
			}
		}
	}

	l := len(emptyPoints)
	if l == 0 {
		return
	}
	p := emptyPoints[rand.Intn(l)]
	gameState.Food = append(gameState.Food, p)
}

func (gameState *GameState) FoodAt(p *Point) bool {
	for _, food := range gameState.Food {
		if food.X == p.X && food.Y == p.Y {
			return true
		}
	}
	return false
}

func (gameState *GameState) IsEmpty(point *Point) bool {
	return !gameState.IsSolid(point, "")
}

func (gameState *GameState) IsSolid(point *Point, ignoreSnakeHead string) bool {
	if point.X < 0 || point.X >= gameState.Width {
		return true
	}
	if point.Y < 0 || point.Y >= gameState.Height {
		return true
	}

	// TODO: take in to account tale shrinks (when no food was eaten)
	for _, snake := range gameState.Snakes {
		for i, coord := range snake.Coords {
			if coord.X == point.X && coord.Y == point.Y {
				if !(i == 0 && snake.Id == ignoreSnakeHead) {
					return true
				}
			}
		}
	}
	return false
}

func (gameState *GameState) IsPossiblySolid(point *Point, ignoreSnakeHead string) bool {
	if gameState.IsSolid(point, ignoreSnakeHead) {
		return true
	}
	for _, snake := range gameState.OtherSnakes() {
		for _, neighbour := range snake.Head().Neighbours() {
			if neighbour.Equals(*point) {
				if snake.Length() > gameState.MySnake().Length() {
					return true
				}
			}
		}
	}

	return false
}

func (gameState *GameState) RemoveSnake(snakeId string) {
	// remove snake
	newSnakes := []*Snake{}
	for _, snake := range gameState.Snakes {
		if snake.Id != snakeId {
			newSnakes = append(newSnakes, snake)
		}
	}
	gameState.Snakes = newSnakes
}

func (gameState *GameState) KillSnake(snakeId string) {
	gameState.losers = append(gameState.losers, gameState.GetSnakeAI(snakeId))
	gameState.RemoveSnake(snakeId)
	gameState.DiedOnTurn[snakeId] = gameState.Turn
}

func (gameState *GameState) MarkSnakeAsWinner(snakeId string) {
	gameState.winners = append(gameState.winners, gameState.GetSnakeAI(snakeId))
	gameState.RemoveSnake(snakeId)
}

func (gameState *GameState) GetSnake(snakeId string) *Snake {
	for _, snake := range gameState.Snakes {
		if snake.Id == snakeId {
			return snake
		}
	}
	return nil
}

func (gameState *GameState) GetSnakeAI(snakeId string) SnakeAI {
	for _, snakeAI := range gameState.SnakeAIs {
		if snakeAI.GetId() == snakeId {
			return snakeAI
		}
	}
	return nil
}
