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
func NewMain(title string) (*Main, error) {
	var err error
	ret := &Main{
		Title:     title,
		z:         4,
		backBytes: make([]byte, ScreenWidth*ScreenHeight*4),
		backFrame: ebiten.NewImage(ScreenWidth, ScreenHeight),
	}
	ret.e, err = engine.NewEngine(ScreenWidth, ScreenHeight)
	if err != nil {
		return nil, err
	}
	return ret, nil
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
