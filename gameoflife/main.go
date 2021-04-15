package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/math/f64"
)

type Camera struct {
	Viewport   f64.Vec2
	Position   f64.Vec2
	ZoomFactor float64
	Rotation   int
}

func NewCamera(screenWidth, screenHeight float64) *Camera {
	return &Camera{Viewport: f64.Vec2{screenWidth, screenWidth}}
}

func (c *Camera) viewportCenter() f64.Vec2 {
	return f64.Vec2{c.Viewport[0] * 0.5, c.Viewport[1] * 0.5}
}

func (c *Camera) worldMatrix() (m ebiten.GeoM) {

	// Move to position
	m.Translate(-c.Position[0], -c.Position[1])

	// Center viewport
	m.Translate(-c.viewportCenter()[0], -c.viewportCenter()[1])

	// Scale camera
	s := math.Pow(1.01, c.ZoomFactor)
	m.Scale(s, s)

	// Do rotations
	m.Rotate(float64(c.Rotation) * 2 * math.Pi / 360)

	// Back to screen space
	m.Translate(c.viewportCenter()[0], c.viewportCenter()[1])
	return
}

func (c *Camera) Render(world, screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{GeoM: c.worldMatrix()}
	screen.DrawImage(world, op)
}

func (c *Camera) ScreenToWorld(posX, posY int) (float64, float64) {
	inverseMatrix := c.worldMatrix()
	if inverseMatrix.IsInvertible() {
		inverseMatrix.Invert()
		return inverseMatrix.Apply(float64(posX), float64(posY))
	} else {
		// When scaling it can happend that matrix is not invertable
		return math.NaN(), math.NaN()
	}
}

// GOL is the game global state for Game Of Life
type GOL struct {
	grid       []int
	width      int
	height     int
	frame      int
	tileSize   int
	v          byte
	refresh    float64
	showDebug  bool
	mx, my, mv int
	mwx, mwy   float64
	isPaused   bool
	showAge    bool
	ms         MouseState
	ageStep    int
	sw, sh     float64
	camera     *Camera
}

// NewGOL returns a new GOL with the given width and height
func NewGOL(width, height, tileSize int) *GOL {
	g := &GOL{width: width, height: height, tileSize: tileSize}
	g.v = 128
	g.refresh = 1
	g.camera = NewCamera(16, 12)
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
	g.camera.Render(img, screen)

	if g.showDebug {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Refresh: %v/sec", 1/g.refresh), 0, 0)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("MP x: %v y: %v v: %v", g.mx, g.my, g.mv), 0, 10)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Scale: %v", g.camera.ZoomFactor), 0, 20)
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
