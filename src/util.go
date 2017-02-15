package main

func stringPtr(str string) *string {
	return &str
}

func getDistanceBetween(p1, p2 Point) *Vector {
	return &Vector{
		X: p2.X - p1.X,
		Y: p2.Y - p1.Y,
	}
}
