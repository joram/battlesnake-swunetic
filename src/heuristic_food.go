package main

func NearestFoodHeuristic(gameState *GameState) WeightedDirections {
	foods := []*Point{}
	for _, food := range gameState.Food {
		foods = append(foods, &food)
	}
	return MoveTo(gameState, foods)
}
