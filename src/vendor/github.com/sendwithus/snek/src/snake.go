package main

func (s Snake) Head() Point {
	return Point{s.Coords[0][0], s.Coords[0][1]}
}

func (s Snake) GetBody() []Point {
	parts := []Point{}
	for _, p := range s.Coords {
		parts = append(parts, Point{p[0], p[1]})
	}
	return parts
}
