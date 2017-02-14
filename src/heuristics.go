package main

// NOTES:
// maybe split into multiple files if this gets too big
// these funcs will take in board state when the other branch gets in


func GoStraightHeuristic(request *MoveRequest) string {

	mySnake := Snake{}
	for _, snake := range request.Snakes {
		if snake.Id == request.You {
			mySnake = snake
		}
	}

	head := mySnake.Coords[0]
	neck := mySnake.Coords[1]
	directionOfMovement := Point{
		X: head.X - neck.X,
		Y: head.Y - neck.Y,
	}
	if directionOfMovement.X == 1 && directionOfMovement.Y == 0 {
		return UP
	}
	if directionOfMovement.X == -1 && directionOfMovement.Y == 0 {
		return DOWN
	}
	if directionOfMovement.X == 0 && directionOfMovement.Y == 1 {
		return RIGHT
	}
	if directionOfMovement.X == 0 && directionOfMovement.Y == -1 {
		return LEFT
	}

	// TODO: random choice of direction if suggested direction is solid
	return nil
}