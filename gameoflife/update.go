package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

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
func (g *GOL) Update() error {
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
