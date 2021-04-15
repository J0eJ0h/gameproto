package main

import (
	"fmt"
	"math/rand"
)

func (g *GOL) flat(x, y int) int {
	return x + y*g.width
}

func (g *GOL) expand(k int) (int, int) {
	return k % g.width, k / g.width
}
func (g *GOL) expandF(k int) (float64, float64) {
	return float64((k % g.width) * g.tileSize), float64((k / g.width) * g.tileSize)
}

func (g *GOL) checkGrid(x, y int) (int, error) {
	k := g.flat(x, y)
	if k < 0 || g.width*g.height <= k {
		return -1, fmt.Errorf("(%v,%v) is out of grid bounds", x, y)
	}
	return k, nil
}

// UpdateGrid updates a given element of the grid with a new value
func (g *GOL) UpdateGrid(x, y, v int) error {
	k, err := g.checkGrid(x, y)
	if err != nil {
		return err
	}
	g.grid[k] = v
	return nil
}

// UpdateGridFlat updates a given element of the grid with a new value
func (g *GOL) updateGridFlat(k, v int) {
	g.grid[k] = v
}

// ReadGridFlat gets the current value of a location in the grid
func (g *GOL) readGridFlat(k int) int {
	return g.grid[k]
}

// ReadGrid gets the current value of a location in the grid
func (g *GOL) ReadGrid(x, y int) (int, error) {
	k, err := g.checkGrid(x, y)
	if err != nil {
		return 0, err
	}
	return g.grid[k], nil
}

// DoGrid applies a visitor pattern to the grid,
// replacing each element with
func (g *GOL) DoGrid(f func(int, int) int) {
	newGrid := make([]int, g.width*g.height)
	for y := 0; y < g.height; y++ {
		for x := 0; x < g.width; x++ {
			newGrid[g.flat(x, y)] = f(x, y)
		}
	}
	g.grid = newGrid
}

func randGrid(x, y int) int {
	if rand.Intn(2) != 0 {
		return 0
	}
	return rand.Intn(256)
}

func (g *GOL) getNeighborCount(x, y int) int {
	pop := 0
	for xi := x - 1; xi <= x+1; xi++ {
		for yi := y - 1; yi <= y+1; yi++ {
			if xi == x && yi == y {
				continue
			}
			if xi < 0 || xi >= g.width || yi < 0 || yi >= g.height {
				continue
			}
			if g.grid[g.flat(xi, yi)] > 0 {
				pop++
			}
		}
	}
	return pop
}

func (g *GOL) ageGrid(x, y int) int {
	v := g.grid[g.flat(x, y)]
	if v > 0 && v < 256-g.ageStep {
		return v + g.ageStep
	}
	return v
}

func (g *GOL) lifeGrid(x, y int) int {
	pop := g.getNeighborCount(x, y)
	v := g.grid[g.flat(x, y)]
	if pop == 3 || (pop == 2 && v > 0) {
		if v > 0 && v < 256 {
			return v
		}
		return 1
	}
	return 0
}
