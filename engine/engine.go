package engine

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/qbradq/eye-engine/data"
	"github.com/qbradq/eye-engine/engine/types"
	"github.com/qbradq/eye-engine/formats"
)

// Engine is the engine structure. Engines are self-contained.
type Engine struct {
	ClearColor types.ColorIndex // Clear color
	Frame      *types.Buffer    // Current frame
	Palette    *types.Palette   // Engine palette
	FPS        float32          // Current FPS average
	MSPF       float32          // Current Milliseconds per Frame average

	w           int           // Frame width
	h           int           // Frame height
	cw          int           // Frame width in cells
	ch          int           // Frame height in cells
	lastFPSCalc time.Time     // Last time the FPS was calculated
	frameTime   time.Duration // Total running time for rendering
	frameCount  int           // Running frame count since last FPS update
	fpsDelay    time.Duration // Delay between FPS updates
	font        *types.Buffer // Font image
	testModel   *types.Model  // Test model, temporary
	Camera      *types.Camera // Camera for the main 3D viewport
}

// NewEngine returns a new Engine ready to use.
func NewEngine(w, h int) *Engine {
	var err error
	ret := &Engine{
		ClearColor:  15,
		fpsDelay:    time.Millisecond * 500,
		lastFPSCalc: time.Now(),
		Camera: &types.Camera{
			Entity: types.Entity{
				Position: types.Backward.Mul(2),
				Forward:  types.Forward,
				Up:       types.Up,
			},
			FOV:      90,
			NearClip: 0.001,
			FarClip:  1,
		},
	}
	ret.Palette, err = types.LoadPalette("default")
	if err != nil {
		slog.Error("error loading default palette", "error", err)
		return nil
	}
	r, err := data.FS.Open(filepath.Join("gfx", "font.png"))
	if err != nil {
		slog.Error("error opening default font", "error", err)
		return nil
	}
	defer r.Close()
	ret.font, err = types.DecodeBuffer(r, ret.Palette)
	if err != nil {
		slog.Error("error loading font into a buffer", "error", err)
		return nil
	}
	r, err = data.FS.Open(filepath.Join("models", "test.obj"))
	if err != nil {
		slog.Error("error opening test model", "error", err)
		return nil
	}
	defer r.Close()
	ret.testModel = formats.LoadModelFromOBJ(r)
	if ret.testModel == nil {
		slog.Error("error loading test model", "error", err)
		return nil
	}
	ret.Resize(w, h)
	return ret
}

// Resize resizes the output buffer.
func (e *Engine) Resize(w, h int) {
	e.w = w
	e.h = h
	e.cw = w / 8
	e.ch = h / 8
	e.Frame = types.NewBuffer(w, h, e.Palette)
	e.Frame.Fill(e.ClearColor)
}

// NextFrame renders the next frame.
func (e *Engine) NextFrame(delta float32) {
	startTime := time.Now()
	// Prep the frame
	e.Frame.Fill(e.ClearColor)
	e.Camera.Target = e.Frame
	e.Camera.Update()
	// Rendering process
	e.Camera.SetModelMatrix(mgl32.Ident4())
	e.Camera.DrawModel(
		e.testModel,
		types.DefaultLightBlue,
		types.DrawModeLines,
	)
	// e.Frame.DrawConvex(
	// 	[]types.PointI2D{
	// 		types.PointI2D{100, 100},
	// 		types.PointI2D{150, 150},
	// 		types.PointI2D{200, 200},
	// 		types.PointI2D{300, 250},
	// 		types.PointI2D{250, 100},
	// 	},
	// 	types.DefaultOrange,
	// )
	// Advance FPS measurement
	e.frameCount++
	endTime := time.Now()
	e.frameTime += endTime.Sub(startTime)
	d := endTime.Sub(e.lastFPSCalc)
	if d > e.fpsDelay {
		e.lastFPSCalc = endTime
		e.FPS = float32(e.frameCount) / float32(d.Seconds())
		e.MSPF = float32(e.frameTime.Seconds()/d.Seconds()) * 1000
		e.frameCount = 0
		e.frameTime = 0
	}
	// FPS display
	fpsStr := fmt.Sprintf("fps:%3.0f ms/f:%3.0f", e.FPS, e.MSPF)
	e.PrintString(e.Frame, fpsStr, e.cw-len(fpsStr), e.ch-1, types.DefaultLime)
}

// PrintChar prints a single character r from the engine font onto dest at cell
// location x, y. Cells are 8 pixels wide and tall.
func (e *Engine) PrintChar(dest *types.Buffer, r rune, x, y int, c types.ColorIndex) {
	sx := int(r) % 16
	sy := int(r) / 16
	dest.DrawMask(e.font, x*8, y*8, sx*8, sy*8, 8, 8, c)
}

// PrintString prints the string s with the engine font onto dest starting at
// cell location x, y. Cells are 8 pixels wide and tall. Does not respect
// special characters or buffer bounds.
func (e *Engine) PrintString(dest *types.Buffer, s string, x, y int, c types.ColorIndex) {
	for _, r := range s {
		e.PrintChar(dest, r, x, y, c)
		x++
	}
}
