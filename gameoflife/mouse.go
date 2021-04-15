package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// MouseState tracks debounces for the mouse to allow MouseDown\Up
// TODO:  finish other buttons
// TODO:  event registration (?via channels)
type MouseState struct {
	leftNow, leftLast, rightLast, rightNow, midLast, midNow bool
	wheelDx, wheelDy                                        float64
}

// Update must be added to your update game loop to handle mouse debounces
func (m *MouseState) Update() {
	m.leftLast = m.leftNow
	m.leftNow = ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)

	m.rightLast = m.rightNow
	m.rightNow = ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight)

	m.midLast = m.midNow
	m.midNow = ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle)

	dx, dy := ebiten.Wheel()
	m.wheelDx += dx
	m.wheelDy += dy
}

func (m *MouseState) ConsumeWheel() (dx, dy float64) {
	dx, dy = m.wheelDx, m.wheelDy
	m.wheelDx, m.wheelDy = 0, 0
	return
}

// LeftDown eturns true if the left mouse button has been pressed since the last call to Update
func (m *MouseState) LeftDown() bool {
	return m.leftNow && !m.leftLast
}

// LeftUp returns true if the left mouse button has been released since the last call to Update
func (m *MouseState) LeftUp() bool {
	return !m.leftNow && m.leftLast
}
