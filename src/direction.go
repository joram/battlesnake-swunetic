package main

import "fmt"

const UP = "up"
const DOWN = "down"
const LEFT = "left"
const RIGHT = "right"
const NOOP = "no-op"

func directionVector(direction string) Point {
	if direction == UP {
		return Point{X: 0, Y: 1}
	}
	if direction == DOWN {
		return Point{X: 0, Y: -1}
	}
	if direction == LEFT {
		return Point{X: -1, Y: 0}
	}
	if direction == RIGHT {
		return Point{X: 1, Y: 0}
	}
	panic(fmt.Sprintf("invalid direction '%v'", direction))
}
