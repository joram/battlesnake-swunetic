package main

import "fmt"

const UP = "u"
const DOWN = "d"
const LEFT = "l"
const RIGHT = "r"


func directionVector(direction string) Point {
	if direction == UP {
		return Point{X:0, Y:1}
	}
	if direction == DOWN {
		return Point{X:0, Y:-1}
	}
	if direction == LEFT {
		return Point{X:-1, Y:0}
	}
	if direction ==RIGHT {
		return Point{X:1, Y:0}
	}
	panic(fmt.Sprintf("invalid direction '%v'", direction))
}