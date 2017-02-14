package main

func (slice WeightedDirections) Len() int {
	return len(slice)
}

func (slice WeightedDirections) Less(i, j int) bool {
	return slice[i].Weight < slice[j].Weight
}

func (slice WeightedDirections) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
