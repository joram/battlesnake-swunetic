package main

import (
	"math/rand"
	"sync"
)

func NewAStar(gameState *GameState, start *Point) AStar {
	aStar := AStar{
		gameState:     gameState,
		start:         start,
		turnsTo:       make(map[Point]int),
		pathToCache:   make(map[Point][]*Point),
		canVisitCount: 0,
	}
	aStar.process()
	return aStar
}

func (a *AStar) shouldVisit(p *Point) bool {
	_, set := a.turnsTo[*p]
	if set {
		return false
	}
	if a.gameState.IsEmpty(p) {
		return true
	}
	return false
}

func (a *AStar) process() {
	initial := AStarPoint{a.start, 0}
	var toVisit = []AStarPoint{initial}
	a.turnsTo[*a.start] = 0
	for len(toVisit) > 0 {
		p := toVisit[0]
		toVisit = toVisit[1:]
		a.canVisitCount += 1
		//println("visiting ", p.point.String())
		for _, neightbour := range p.point.Neighbours() {
			if a.shouldVisit(neightbour) {
				a.turnsTo[*neightbour] = p.turnsTo + 1
				next := AStarPoint{neightbour, p.turnsTo + 1}
				toVisit = append(toVisit, next)
			}
		}
	}
}

func (a *AStar) previousStep(to *Point) *Point {
	nextOptions := []*Point{}
	timeToCurr := a.turnsTo[*to]
	for _, neighbour := range to.Neighbours() {
		timeTo, set := a.turnsTo[*neighbour]
		if set && timeTo < timeToCurr {
			nextOptions = append(nextOptions, neighbour)
		}
	}

	// no path
	if len(nextOptions) == 0 {
		return nil
	}

	next := nextOptions[rand.Intn(len(nextOptions))]
	return next
}

var pathToCacheLock sync.Mutex

func (a *AStar) pathTo(to *Point) []*Point {
	pathToCacheLock.Lock()
	path := a.pathToCache[*to]
	pathToCacheLock.Unlock()

	if len(path) > 0 {
		return path
	}

	curr := to
	for !curr.Equals(*a.start) {
		path = append(path, curr)
		curr = a.previousStep(curr)
		if curr == nil {
			break
		}
	}

	// reverse
	if curr != nil {
		for i := len(path)/2 - 1; i >= 0; i-- {
			opp := len(path) - 1 - i
			path[i], path[opp] = path[opp], path[i]
		}
	}

	pathToCacheLock.Lock()
	a.pathToCache[*to] = path
	pathToCacheLock.Unlock()

	return path
}
