package grid

import (
	"fmt"
)

type Grid interface {
	// ReadGrid gets the current value of a location in the grid
	ReadGrid(x, y int) (int, error)

	// UpdateGrid updates a given element of the grid with a new value
	UpdateGrid(x, y, v int) error

	// DoGrid applies a visitor pattern to the grid,
	// replacing each element with the result
	// after calculating the updates
	DoGrid(f func(int, int) int)

	// Dim returns the minx, miny, maxx and maxy values for the grid
	// Depending on the underlying implementation, not all spaces may
	// be populated.
	// Also all values are inclusive, so both min* and max* are addressable
	// locations in the grid and should not error
	Dim() (int, int, int, int)
}

type flatGrid struct {
	grid   []int
	width  int
	height int
}

func FlatGrid(width, height int) Grid {
	return &flatGrid{grid: make([]int, width*height), width: width, height: height}
}

func (g *flatGrid) flat(x, y int) int {
	return x + y*g.width
}

func (g *flatGrid) checkGrid(x, y int) (int, error) {
	k := g.flat(x, y)
	if k < 0 || g.width*g.height <= k {
		return -1, fmt.Errorf("(%v,%v) is out of grid bounds", x, y)
	}
	return k, nil
}

func (g *flatGrid) ReadGrid(x, y int) (int, error) {
	k, err := g.checkGrid(x, y)
	if err != nil {
		return 0, err
	}
	return g.grid[k], nil
}

func (g *flatGrid) UpdateGrid(x, y, v int) error {
	k, err := g.checkGrid(x, y)
	if err != nil {
		return err
	}
	g.grid[k] = v
	return nil
}

func (g *flatGrid) Dim() (int, int, int, int) {
	return 0, 0, g.width - 1, g.height - 1
}

func (g *flatGrid) DoGrid(f func(int, int) int) {
	newGrid := make([]int, g.width*g.height)
	for y := 0; y < g.height; y++ {
		for x := 0; x < g.width; x++ {
			newGrid[g.flat(x, y)] = f(x, y)
		}
	}
	g.grid = newGrid
}
