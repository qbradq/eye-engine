package ui

import (
	"log/slog"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/qbradq/eye-engine/engine"
	"github.com/qbradq/eye-engine/internal/util"
)

const ScreenWidth int = 320  // Width of the engine screen
const ScreenHeight int = 240 // Height of the engine screen
const ZoomMin int = 1        // Minimum zoom level
const ZoomMax int = 9        // Maximum zoom level

// Main manages the entire UI.
type Main struct {
	Title string // Title of the application

	z             int            // Zoom level
	e             *engine.Engine // The graphics engine
	lastFrameTime time.Time      // Time of last frame
	backBytes     []byte         // Frame backing data
	backFrame     *ebiten.Image  // Frame backing image
}

// NewMain returns a new Main object ready for use.
func NewMain(title string) *Main {
	ret := &Main{
		Title:     title,
		z:         4,
		backBytes: make([]byte, ScreenWidth*ScreenHeight*4),
		backFrame: ebiten.NewImage(ScreenWidth, ScreenHeight),
	}
	ret.e = engine.NewEngine(ScreenWidth, ScreenHeight)
	if ret.e == nil {
		return nil
	}
	return ret
}

// Main is the app entry point.
func (m *Main) Main() {
	ebiten.SetTPS(60)
	ebiten.SetWindowDecorated(true)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)
	ebiten.SetWindowTitle(m.Title)
	m.Zoom(m.z)
	if err := ebiten.RunGame(m); err != nil {
		slog.Error("error executing main()", "error", err)
	}
}

// Zoom sets the zoom level.
func (m *Main) Zoom(z int) {
	m.z = util.BoundInt(z, ZoomMin, ZoomMax)
	ebiten.SetWindowSize(ScreenWidth*m.z, ScreenHeight*m.z)
}

// Layout returns the native resolution.
func (m *Main) Layout(w, h int) (int, int) {
	return ScreenWidth, ScreenHeight
}

// Update is called for input and physics updates.
func (m *Main) Update() error {
	// Handle keyboard rotation
	keyTurnSpeed := float32(360.0 * 1 / 60.0)
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		m.e.Camera.Rotate(keyTurnSpeed)
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		m.e.Camera.Rotate(-keyTurnSpeed)
	}
	// Handle mouse look
	// ebiten.SetCursorMode(ebiten.CursorModeCaptured)
	// mouseLookSpeed := float32(360.0 / 60.0)
	// var pos types.PointI2D
	// pos[0], pos[1] = ebiten.CursorPosition()
	// center := types.PointI2D{
	// 	ScreenWidth / 2,
	// 	ScreenHeight / 2,
	// }
	// dci := types.PointI2D{
	// 	pos[0] - center[0],
	// 	pos[1] - center[1],
	// }
	// dcf := mgl32.Vec2{
	// 	float32(dci[0]) / float32(ScreenWidth),
	// 	float32(dci[1]) / float32(ScreenWidth),
	// }
	// m.e.Camera.Rotate(dcf[0] * mouseLookSpeed)
	// Handle movement
	moveSpeed := float32(1.0 / 60.0)
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		m.e.Camera.MoveForward(moveSpeed)
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		m.e.Camera.MoveBackward(moveSpeed)
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		m.e.Camera.StrafeRight(moveSpeed)
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		m.e.Camera.StrafeLeft(moveSpeed)
	}
	return nil
}

// Draw is called to draw each frame.
func (m *Main) Draw(screen *ebiten.Image) {
	var d float32
	if m.lastFrameTime.IsZero() {
		d = 1.0 / 60.0
	} else {
		d = float32(time.Now().UnixMilli()-m.lastFrameTime.UnixMilli()) / 1000
	}
	m.e.NextFrame(d)
	m.e.Frame.Bytes(m.backBytes)
	m.backFrame.WritePixels(m.backBytes)
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(m.backFrame, op)
}
