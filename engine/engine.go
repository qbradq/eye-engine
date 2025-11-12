package engine

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/qbradq/eye-engine/data"
)

// Engine is the engine structure. Engines are self-contained.
type Engine struct {
	ClearColor ColorIndex // Clear color
	Frame      *Buffer    // Current frame
	Palette    *Palette   // Engine palette
	FPS        float32    // Current FPS average
	MSPF       float32    // Current Milliseconds per Frame average

	w           int           // Frame width
	h           int           // Frame height
	cw          int           // Frame width in cells
	ch          int           // Frame height in cells
	lastFPSCalc time.Time     // Last time the FPS was calculated
	frameTime   time.Duration // Total running time for rendering
	frameCount  int           // Running frame count since last FPS update
	fpsDelay    time.Duration // Delay between FPS updates
	font        *Buffer       // Font image
}

// NewEngine returns a new Engine ready to use.
func NewEngine(w, h int) (*Engine, error) {
	var err error
	ret := &Engine{
		ClearColor:  15,
		fpsDelay:    time.Millisecond * 500,
		lastFPSCalc: time.Now(),
	}
	ret.Palette, err = LoadPalette("default")
	if err != nil {
		return nil, err
	}
	r, err := data.FS.Open(filepath.Join("gfx", "font.png"))
	if err != nil {
		return nil, err
	}
	ret.font, err = DecodeBuffer(r, ret.Palette)
	if err != nil {
		return nil, err
	}
	ret.Resize(w, h)
	return ret, nil
}

// Resize resizes the output buffer.
func (e *Engine) Resize(w, h int) {
	e.w = w
	e.h = h
	e.cw = w / 8
	e.ch = h / 8
	e.Frame = NewBuffer(w, h, e.Palette)
	e.Frame.Fill(e.ClearColor)
}

// NextFrame renders the next frame.
func (e *Engine) NextFrame(delta float32) {
	startTime := time.Now()
	// Prep the frame
	e.Frame.Fill(e.ClearColor)
	// Rendering process
	e.PrintChar(e.Frame, '@', 4, 1, DefaultBlue)
	e.PrintString(e.Frame, "Hello, Eye Engine!", 0, 0, DefaultLime)
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
	e.PrintString(e.Frame, fpsStr, e.cw-len(fpsStr), e.ch-1, DefaultLime)
}

// PrintChar prints a single character r from the engine font onto dest at cell
// location x, y. Cells are 8 pixels wide and tall.
func (e *Engine) PrintChar(dest *Buffer, r rune, x, y int, c ColorIndex) {
	sx := int(r) % 16
	sy := int(r) / 16
	dest.DrawMask(e.font, x*8, y*8, sx*8, sy*8, 8, 8, c)
}

// PrintString prints the string s with the engine font onto dest starting at
// cell location x, y. Cells are 8 pixels wide and tall. Does not respect
// special characters or buffer bounds.
func (e *Engine) PrintString(dest *Buffer, s string, x, y int, c ColorIndex) {
	for _, r := range s {
		e.PrintChar(dest, r, x, y, c)
		x++
	}
}
