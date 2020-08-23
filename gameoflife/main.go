package main

import (
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten"
)

// GOL is the game global state for Game Of Life
type GOL struct {
	grid   []int
	width  int
	height int
	frame  int
}

// NewGOL returns a new GOL with the given width and height
func NewGOL(width, height int) *GOL {
	g := &GOL{width: width, height: height}
	return g
}

func (g *GOL) flat(x, y int) int {
	return x + y*g.width
}

func (g *GOL) expand(k int) (int, int) {
	return k % g.width, k / g.width
}

// UpdateGrid updates a given element of the grid with a new value
func (g *GOL) UpdateGrid(x, y, v int) {
	g.grid[g.flat(x, y)] = v
}

// UpdateGridFlat updates a given element of the grid with a new value
func (g *GOL) UpdateGridFlat(k, v int) {
	g.grid[k] = v
}

// ReadGridFlat gets the current value of a location in the grid
func (g *GOL) ReadGridFlat(k int) int {
	return g.grid[k]
}

// ReadGrid gets the current value of a location in the grid
func (g *GOL) ReadGrid(x, y int) int {
	return g.grid[g.flat(x, y)]
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

func (g *GOL) lifeGrid(x, y int) int {
	pop := g.getNeighborCount(x, y)
	if pop == 3 {
		return 255
	}
	if pop == 2 && g.grid[g.flat(x, y)] > 0 {
		return 1
	}
	return 0
}

// Draw is the draw function for GOL games
func (g *GOL) Draw(screen *ebiten.Image) {
	pixels := make([]byte, 4*g.width*g.height)
	for k := 0; k < g.width*g.height; k++ {
		a := g.grid[k]
		b := 255 - g.grid[k]
		c := 0
		if g.grid[k] == 0 {
			a = 0
			b = 0
			c = 255
		}

		pixels[4*k] = byte(a)
		pixels[4*k+1] = byte(b)
		pixels[4*k+2] = byte(c)
		pixels[4*k+3] = 255

	}
	screen.ReplacePixels(pixels)

}

// Update is the game state update function for GOL
func (g *GOL) Update(*ebiten.Image) error {
	g.frame++
	tps := ebiten.MaxTPS()
	if g.frame%(1*tps) == 0 {
		g.DoGrid(g.lifeGrid)
	}
	return nil
}

// Layout sets the window : screen layout for GOL
func (g *GOL) Layout(int, int) (int, int) {
	return g.width, g.height
}

func main() {
	rand.Seed(time.Now().UnixNano())
	gol := NewGOL(80, 60)
	gol.DoGrid(randGrid)

	ebiten.SetWindowSize(320, 240)
	ebiten.SetWindowTitle("GOL")

	ebiten.RunGame(gol)
}
