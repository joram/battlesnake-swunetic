package main

func (snake *Snake) Extend(direction string) Point {
	curHead := snake.Coords[0]
	newHead := curHead.Add(directionVector(direction))
	newCoords := []Point{newHead}
	for _, coord := range snake.Coords {
		newCoords = append(newCoords, coord)
	}
	snake.Coords = newCoords
	return newHead
}

func (snake *Snake) Shrink() {
	snake.Coords = snake.Coords[len(snake.Coords)-1:]
}
