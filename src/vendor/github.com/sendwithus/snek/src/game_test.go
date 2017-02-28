package snek

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMoveRequest_IsValidMove(t *testing.T) {
	m := MoveRequest{
		Width:  20,
		Height: 20,
		You:    "1",
		Snakes: []Snake{
			{
				Coords: [][]int{
					{1, 1},
					{1, 2},
				},
				Id: "1",
			},
		},
	}

	assert.True(t, m.IsValidMove(UP, false))
	assert.True(t, m.IsValidMove(LEFT, false))
	assert.True(t, m.IsValidMove(RIGHT, false))
	assert.False(t, m.IsValidMove(DOWN, false))

	m.Snakes[0].Coords[0][0] = 0
	m.Snakes[0].Coords[0][1] = 0
	assert.False(t, m.IsValidMove(UP, false))
	assert.False(t, m.IsValidMove(LEFT, false))

	m.Snakes[0].Coords[0][0] = m.Width - 1
	m.Snakes[0].Coords[0][1] = m.Height - 1
	assert.False(t, m.IsValidMove(DOWN, false))
	assert.False(t, m.IsValidMove(RIGHT, false))
}

func TestMoveRequest_IsLocationEmpty(t *testing.T) {
	m := MoveRequest{
		Width:  20,
		Height: 20,
	}
	snake := Snake{
		Coords: [][]int{
			{1, 1},
			{1, 2},
		},
	}
	m.Snakes = append(m.Snakes, snake)
	assert.False(t, m.IsLocationEmpty(Point{1, 1}))
	assert.True(t, m.IsLocationEmpty(Point{2, 1}))
}

func TestMoveRequest_FindMoveToNearestFood(t *testing.T) {
	m := MoveRequest{
		Width:  20,
		Height: 20,
		Food: [][]int{
			{4, 1},
		},
		You: "1",
		Snakes: []Snake{
			{
				Coords: [][]int{
					{1, 1},
					{1, 2},
				},
				Id: "1",
			},
		},
	}

	move := m.FindMoveToNearestFood()
	assert.Equal(t, RIGHT, move)

	m.Food[0][0] = 0
	move = m.FindMoveToNearestFood()
	assert.Equal(t, LEFT, move)

	m.Food[0][0] = 1
	m.Food[0][1] = 0
	move = m.FindMoveToNearestFood()
	assert.Equal(t, UP, move)

	m.Food[0][1] = 5
	move = m.FindMoveToNearestFood()
	assert.Equal(t, NOOP, move)
}

func TestMoveRequest_AddNodes(t *testing.T) {
	m := MoveRequest{
		Width:  20,
		Height: 20,
		Food: [][]int{
			{4, 1},
		},
		You: "1",
		Snakes: []Snake{
			{
				Coords: [][]int{
					{2, 1},
					{1, 1},
					{1, 2},
					{1, 3},
					{2, 3},
					{3, 3},
					{3, 2},
				},
				Id: "1",
			},
		},
	}
	assert.True(t, m.SearchForClosedArea(Point{2, 2}))
	assert.False(t, m.SearchForClosedArea(Point{2, 0}))
}

func TestMoveRequest_CheckForPossibleKills(t *testing.T) {
	m := MoveRequest{
		Width:  20,
		Height: 20,
		Food: [][]int{
			{4, 1},
		},
		You: "1",
		Snakes: []Snake{
			{
				Coords: [][]int{
					{1, 1},
					{1, 2},
				},
				Id: "1",
			},
			{
				Id: "2",
				Coords: [][]int{
					{2, 0},
				},
			},
		},
	}
	assert.Equal(t, UP, m.CheckForPossibleKills())

	m.Snakes[1].Coords = append(m.Snakes[1].Coords, []int{2, 1})
	assert.Equal(t, NOOP, m.CheckForPossibleKills())
}
