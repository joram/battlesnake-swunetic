package main

func (moveRequest *MoveRequest) IsSolid(point Point) bool{
	if point.X < 0 || point.X >= moveRequest.Width {
		return true
	}
	if point.Y < 0 || point.Y >= moveRequest.Height {
		return true
	}

	// TODO: take in to account tale shrinks (when no food was eaten)
	for _, snake := range moveRequest.Snakes {
		for _, coord := range snake.Coords {
			if coord.X == point.X && coord.Y == point.Y {
				return true
			}
		}
	}
	return false
}

func (moveRequest *MoveRequest) KillSnake(snakeId string) {
	newSnakes := []Snake{}
	for _, snake := range moveRequest.Snakes {
		if snake.Id != snakeId {
			newSnakes = append(newSnakes, snake)
		}
	}
	moveRequest.Snakes = newSnakes
}


func (moveRequest *MoveRequest) GetSnake(snakeId string) *Snake {
	for i, snake := range moveRequest.Snakes {
		if snake.Id == snakeId {
			return &moveRequest.Snakes[i]
		}
	}
	return nil
}

