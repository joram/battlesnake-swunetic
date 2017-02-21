package main

import "fmt"

func (p *Point) Add(p2 Point) Point {
	return Point{
		X: p.X + p2.X,
		Y: p.Y + p2.Y,
	}
}

func (p *Point) Subtract(p2 Point) Point {
	return Point{
		X: p.X - p2.X,
		Y: p.Y - p2.Y,
	}
}

func (p *Point) Left() *Point {
	return &Point{X: p.X - 1, Y: p.Y}
}

func (p *Point) Right() *Point {
	return &Point{X: p.X + 1, Y: p.Y}
}

func (p *Point) Up() *Point {
	return &Point{X: p.X, Y: p.Y - 1}
}

func (p *Point) Down() *Point {
	return &Point{X: p.X, Y: p.Y + 1}
}

func (p *Point) Neighbours() []*Point {
	return []*Point{
		p.Up(),
		p.Down(),
		p.Left(),
		p.Right(),
	}
}

func (p *Point) Equals(p2 Point) bool {
	return p.X == p2.X && p.Y == p2.Y
}

func (p *Point) String() string {
	return fmt.Sprintf("(%v,%v)", p.X, p.Y)
}

func (p *Point) isCloser(p2 *Point, goal *Point) bool {
	dist1 := getDistanceBetween(*p, *goal).Magnitude()
	dist2 := getDistanceBetween(*p2, *goal).Magnitude()
	return dist1 < dist2
}
