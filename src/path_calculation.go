package main

func NewPathCalculation(from *Point, to Points, gs *GameState) PathCalculation {
	return PathCalculation{
		gameState: gs,
		start:     from,
		goals:     to,
		visited:   []PointPair{{from: from, to: from}},
		toVisit:   Points{from},
	}
}

func (pc *PathCalculation) isAGoal(p *Point) bool {
	for _, goal := range pc.goals {
		if goal.Equals(*p) {
			return true
		}
	}
	return false
}

func (pc *PathCalculation) removeGoal(p *Point) {
	for i, goal := range pc.goals {
		if goal.Equals(*p) {
			first := pc.goals[:i]
			last := Points{}
			if len(pc.goals) > i {
				last = pc.goals[i+1:]
			}
			pc.goals = append(first, last...) // remove goal
		}
	}
}

func (pc *PathCalculation) haveVisited(p *Point) bool {
	for _, visited := range pc.visited {
		if visited.to.Equals(*p) {
			return true
		}
	}
	return false
}

func (pc *PathCalculation) Run() {
	for len(pc.toVisit) > 0 {
		p := pc.toVisit[0]
		pc.toVisit = pc.toVisit[1:]
		if pc.isAGoal(p) {
			pc.removeGoal(p)
			pc.achievedGoals = append(pc.achievedGoals, p)
			continue
		}

		for _, neighbour := range p.Neighbours() {
			if !pc.haveVisited(neighbour) && pc.gameState.IsEmpty(*neighbour) {
				pc.visited = append(pc.visited, PointPair{from: p, to: neighbour})
				pc.toVisit = append(pc.toVisit, neighbour)
			}
		}
	}
}

func (pc *PathCalculation) PathFrom(currentPoint *Point) Points {
	path := Points{currentPoint}
	for !currentPoint.Equals(*pc.start) {
		for _, pointPair := range pc.visited {
			if pointPair.to.Equals(*currentPoint) {
				path = append(path, pointPair.from)
				currentPoint = pointPair.from
				break
			}
		}
	}
	return path
}

func (pc *PathCalculation) Paths() []Points {
	paths := []Points{}
	for _, from := range pc.achievedGoals {
		paths = append(paths, pc.PathFrom(from))
	}
	return paths
}

func (pc *PathCalculation) PreviousPoint(p *Point) *Point {
	for _, pPair := range pc.visited {
		if pPair.to.Equals(*p) {
			return pPair.from
		}
	}
	return nil
}
