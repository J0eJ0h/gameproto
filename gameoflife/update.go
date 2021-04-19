package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (g *GOL) doKeyboardUpdate() {
	// game behavior controls
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

	// camera controls
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.camera.Position[0] -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.camera.Position[0] += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.camera.Position[1] -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.camera.Position[1] += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		g.camera.Rotation -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		g.camera.Rotation += 1
	}
}

func (g *GOL) WorldToTile(x, y float64) (int, int) {
	return int(math.Floor(x / float64(g.tileSize))), int(math.Floor(y / float64(g.tileSize)))
}

func (g *GOL) doMouseUpdate() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {

		x, y := g.camera.ScreenToWorld(ebiten.CursorPosition())
		tx, ty := g.WorldToTile(x, y)

		v, _ := g.grid.ReadGrid(tx, ty)
		if v == 0 {
			v = 1
		} else {
			v = 0
		}
		g.grid.UpdateGrid(tx, ty, v)
		g.mx, g.my, g.mv = tx, ty, v
	}
	_, mwy := g.ms.ConsumeWheel()
	g.camera.ZoomFactor += mwy

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
		g.grid.DoGrid(g.lifeGrid)
		g.grid.DoGrid(g.ageGrid)
	}

	// Final updates
	g.frame++
	return nil
}
