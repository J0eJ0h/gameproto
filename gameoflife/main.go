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
	return &Camera{Viewport: f64.Vec2{screenWidth, screenHeight}}
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
	grid       Grid
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

	g := &GOL{tileSize: tileSize}
	g.grid = FlatGrid(width, height)
	g.v = 128
	g.refresh = 1
	g.sw, g.sh = 16, 12
	g.camera = NewCamera(g.sw*float64(g.tileSize), g.sh*float64(g.tileSize))
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
	minx, miny, maxx, maxy := g.grid.Dim()
	img := ebiten.NewImage((maxx-minx+1)*g.tileSize, (maxy-miny+1)*g.tileSize)
	tile := g.makeTile(g.v, g.tileSize)

	// Render the whole thing for now
	g.grid.DoGrid(func(x, y int) int {
		v, _ := g.grid.ReadGrid(x, y)
		if v == 0 {
			return v
		}

		// expand
		xl, yl := float64(x*g.tileSize), float64(y*g.tileSize)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(xl, yl)
		if g.showAge {
			tile = g.makeTile(byte(v), g.tileSize)
		}
		img.DrawImage(tile, op)
		return v
	})

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
	return int(g.camera.Viewport[0]), int(g.camera.Viewport[1])
}

func main() {
	rand.Seed(time.Now().UnixNano())
	gol := NewGOL(20, 15, 25)
	gol.grid.DoGrid(randGrid)

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("GOL")

	ebiten.RunGame(gol)
}
