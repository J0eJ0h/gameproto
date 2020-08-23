package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
)

// GOL is the game global state for Game Of Life
type GOL struct {
	grid      []int
	width     int
	height    int
	frame     int
	tileSize  int
	v         byte
	refresh   float64
	showDebug bool
	vp        viewport
	mx, my    int
	isPaused  bool
	showAge   bool
	ms        MouseState
	ageStep   int
}

type viewport struct {
	x, y, w, h int
}

// NewGOL returns a new GOL with the given width and height
func NewGOL(width, height, tileSize int) *GOL {
	g := &GOL{width: width, height: height, tileSize: tileSize}
	g.v = 128
	g.refresh = 1
	g.vp = viewport{0, 0, width, height}
	g.ageStep = 10
	return g
}

func (g *GOL) flat(x, y int) int {
	return x + y*g.width
}

func (g *GOL) expand(k int) (int, int) {
	return k % g.width, k / g.width
}
func (g *GOL) expandF(k int) (float64, float64) {
	return float64(k % g.width * g.tileSize), float64(k / g.width * g.tileSize)
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

func (g *GOL) makeTile(v byte, tileSize int) *ebiten.Image {
	//Err is explicitly always null here
	tile, _ := ebiten.NewImage(tileSize, tileSize, ebiten.FilterDefault)
	pixels := make([]byte, 4*tileSize*tileSize)
	border := 1 + (tileSize / 10)
	for y := 0; y < tileSize; y++ {
		for x := 0; x < tileSize; x++ {
			k := x + y*tileSize
			a := v
			b := 255 - v
			c := 0

			if g.isPaused {
				c = 255
			}

			pixels[4*k] = byte(a)
			pixels[4*k+1] = byte(b)
			pixels[4*k+2] = byte(c)
			if x < border || x >= tileSize-border || y < border || y >= tileSize-border {
				pixels[4*k+3] = 50
			} else {
				pixels[4*k+3] = 255
			}
		}

	}
	tile.ReplacePixels(pixels)
	return tile
}

// Draw is the draw function for GOL games
func (g *GOL) Draw(screen *ebiten.Image) {
	tile := g.makeTile(g.v, g.tileSize)
	for k := 0; k < g.width*g.height; k++ {
		if g.grid[k] == 0 {
			continue
		}
		xl, yl := g.expandF(k)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(xl, yl)
		if g.showAge {
			tile = g.makeTile(byte(g.grid[k]), g.tileSize)
		}
		screen.DrawImage(tile, op)
	}
	if g.showDebug {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Refresh: %v/sec", 1/g.refresh), 0, 0)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("MP x: %v y: %v", g.mx, g.my), 0, 10)
	}
}

func (g *GOL) doKeyboardUpdate() {
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.v = 255
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyG) {
		g.v = 0
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEqual) {
		g.refresh = g.refresh / 2
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyMinus) {
		g.refresh = g.refresh * 2
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF3) {
		g.showDebug = !g.showDebug
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.isPaused = !g.isPaused
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyV) {
		g.showAge = !g.showAge
	}
}

func (g *GOL) mapLocToTile(x, y int) int {
	tx, ty := x/g.tileSize, y/g.tileSize
	return g.flat(tx, ty)
}

func (g *GOL) doMouseUpdate() {
	if g.ms.LeftDown() {

		x, y := ebiten.CursorPosition()
		k := g.mapLocToTile(x, y)

		if g.grid[k] == 0 {
			g.grid[k] = 1
		} else {
			g.grid[k] = 0
		}

		g.mx, g.my = x/g.tileSize, y/g.tileSize
	}
}

// Update is the game state update function for GOL
func (g *GOL) Update(*ebiten.Image) error {
	// housekeeping/subsystems
	g.ms.Update()

	// process inputs
	g.doKeyboardUpdate()
	g.doMouseUpdate()

	// Do work
	tps := ebiten.MaxTPS()
	if g.isPaused {
		return nil
	}
	if g.refresh == 0 || g.frame%int(math.Ceil(g.refresh*float64(tps))) == 0 {
		g.DoGrid(g.lifeGrid)
		g.DoGrid(g.ageGrid)
	}

	// Final updates
	g.frame++
	return nil
}

// Layout sets the window : screen layout for GOL
func (g *GOL) Layout(int, int) (int, int) {
	return g.width * g.tileSize, g.height * g.tileSize
}

func main() {
	rand.Seed(time.Now().UnixNano())
	gol := NewGOL(20, 15, 10)
	gol.DoGrid(randGrid)

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("GOL")

	ebiten.RunGame(gol)
}
