package main

type Game struct {
	snakes       []HeuristicSnake
	moveCount    int
	currentBoard MoveRequest
	state        string
	winners      []HeuristicSnake
}

func NewGame(weights [][]int) Game {
	snakes := []HeuristicSnake{}
	for _, weights := range weights {
		snakes = append(snakes, NewHeuristicSnake(weights))
	}

	startingBoard := MoveRequest{} // TODO: populate

	return Game{
		snakes:       snakes,
		moveCount:    0,
		currentBoard: startingBoard,
		winners:      []HeuristicSnake{},
	}
}

func (game *Game) Run() []HeuristicSnake {
	game.state = "running"
	for {
		game.Tick()
		if game.state != "running" {
			println("game ended: %v on move %v", game.state, game.moveCount)
			for _, snake := range game.snakes {
				weights := []int{}
				for _, w := range snake.weightedHeuristics {
					weights = append(weights, w.weight)
				}
				println("winner: %v", weights)
			}
			break
		}
	}
	return game.winners
}

func (game *Game) Tick() {
	board := game.currentBoard
	for _, snake := range game.snakes {
		direction := snake.Move(&board)
		game.MakeMove(&snake, direction)
	}
	game.moveCount += 1
}

func (game *Game) MakeMove(snake *HeuristicSnake, direction string) {
	// TODO: update game.currentBoard based on direction
	// Possible side effects:
	//	- sensible move, board's updated
	//	- death
	//	- win state=win
	//	- draw state=draw
}
