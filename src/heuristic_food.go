package main

func NearestFoodHeuristic(gameState *GameState) WeightedDirections {
	foods := []*Point{}
	for _, food := range gameState.Food {
		foods = append(foods, &food)
	}

	allOtherSnakesWillStarve := true
	for _, snake := range gameState.OtherSnakes() {
		canGetToFood := false
		for _, food := range gameState.Food {
			if gameState.aStar[snake.Id].turnsTo[food] < snake.HealthPoints {
				canGetToFood = true
				break
			}
		}
		if canGetToFood {
			allOtherSnakesWillStarve = false
			break
		}
	}

	if allOtherSnakesWillStarve {
		healthiestSnake := true
		for _, snake := range gameState.OtherSnakes() {
			if snake.HealthPoints >= gameState.MySnake().HealthPoints {
				healthiestSnake = false
				break
			}
		}
		if healthiestSnake {
			return WeightedDirections{}
		}
	}

	return MoveTo(gameState, foods)
}
