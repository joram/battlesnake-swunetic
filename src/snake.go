package main

func (snake *Snake) Extend(direction string) Point {
	curHead := snake.Coords[0]
	newHead := curHead.Add(directionVector(direction))
	snake.Coords = append([]Point{newHead}, snake.Coords...)
	return newHead
}

func (snake *Snake) Shrink() {
	snake.Coords = snake.Coords[len(snake.Coords)-1:]
}

func (snake *Snake) Length() int {
	return len(snake.Coords)
}

func (snake *Snake) Head() *Point {
	if len(snake.Coords) == 0 {
		return nil
	}
	return &snake.Coords[0]
}
