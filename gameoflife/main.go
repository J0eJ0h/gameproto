package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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
	sw, sh    float64
}

type viewport struct {
	x, y, w, h float64
}

// NewGOL returns a new GOL with the given width and height
func NewGOL(width, height, tileSize int) *GOL {
	g := &GOL{width: width, height: height, tileSize: tileSize}
	g.v = 128
	g.refresh = 1
	g.vp = viewport{1, 1, float64(width) - 2, float64(height) - 2}
	g.sw, g.sh = 16, 12
	g.ageStep = 10
	return g
}

func (g *GOL) makeTile(v byte, tileSize int) *ebiten.Image {
	//Err is explicitly always null here
	tile := ebiten.NewImage(tileSize, tileSize)
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

func (g *GOL) renderImage() *ebiten.Image {
	img := ebiten.NewImage(g.width*g.tileSize, g.height*g.tileSize)
	tile := g.makeTile(g.v, g.tileSize)
	for k := 0; k < g.width*g.height; k++ {
		if g.grid[k] == 0 {
			continue
		}

		// expand
		xl, yl := g.expandF(k)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(xl, yl)
		if g.showAge {
			tile = g.makeTile(byte(g.grid[k]), g.tileSize)
		}
		img.DrawImage(tile, op)
	}
	return img
}

// Draw is the draw function for GOL games
func (g *GOL) Draw(screen *ebiten.Image) {

	img := g.renderImage()
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(img, op)

	if g.showDebug {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Refresh: %v/sec", 1/g.refresh), 0, 0)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("MP x: %v y: %v", g.mx, g.my), 0, 10)
	}
}

// Layout sets the window : screen layout for GOL
func (g *GOL) Layout(int, int) (int, int) {
	return int(g.sw) * g.tileSize, int(g.sh) * g.tileSize
}

func main() {
	rand.Seed(time.Now().UnixNano())
	gol := NewGOL(20, 15, 25)
	gol.DoGrid(randGrid)

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("GOL")

	ebiten.RunGame(gol)
}
