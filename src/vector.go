package main

import "math"

func (v *Vector) Magnitude() float64 {
	return math.Sqrt(float64((v.X * v.X) + (v.Y * v.Y)))
}
