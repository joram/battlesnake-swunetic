package main

func (mr *MoveRequest) HorribleJsonMassagingHack() {
	for i, snake := range mr.Snakes {
		body := []Point{}
		for _, coord := range snake.RawCoords {
			part := Point{X: coord[0], Y: coord[1]}
			body = append(body, part)
		}
		mr.Snakes[i].Coords = body
	}

	for _, coord := range mr.RawFood {
		food := Point{X: coord[0], Y: coord[1]}
		mr.Food = append(mr.Food, food)

	}
}
