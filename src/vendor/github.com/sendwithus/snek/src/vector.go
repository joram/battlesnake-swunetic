package snek

import "math"

type Vectors []Vector
type Points []Point

func (l Points) Contains(search Point) bool {
	for _, p := range l {
		if p.Equals(search) {
			return true
		}

	}
	return false
}

func (v Vector) Magnitude() float64 {
	return math.Sqrt(float64((v.X * v.X) + (v.Y * v.Y)))
}

func (p Point) DistanceTo(other Point) Vector {
	return Vector{
		X: other.X - p.X,
		Y: other.Y - p.Y,
	}
}

func (p Point) Add(dir string) Point {
	newLocation := p
	if dir == UP {
		newLocation.Y--
	} else if dir == DOWN {
		newLocation.Y++
	} else if dir == LEFT {
		newLocation.X--
	} else if dir == RIGHT {
		newLocation.X++
	}
	return newLocation
}

func (p Point) Equals(other Point) bool {
	return p.X == other.X && p.Y == other.Y
}

func (slice Vectors) Len() int {
	return len(slice)
}

func (slice Vectors) Less(i, j int) bool {
	return slice[i].Magnitude() < slice[j].Magnitude()
}

func (slice Vectors) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
