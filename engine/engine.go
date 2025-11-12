package engine

import "time"

// Engine is the engine structure. Engines are self-contained.
type Engine struct {
	ClearColor ColorIndex // Clear color
	Frame      *Buffer    // Current frame
	Palette    *Palette   // Engine palette
	FPS        float32    // Current FPS average
	LastMSPF   float32    // Current Milliseconds per Frame average

	w           int           // Frame width
	h           int           // Frame height
	lastFPSCalc time.Time     // Last time the FPS was calculated
	frameTime   time.Duration // Total running time for rendering
	frameCount  int           // Running frame count since last FPS update
	fpsDelay    time.Duration // Delay between FPS updates
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
	ret.Resize(w, h)
	return ret, nil
}

// Resize resizes the output buffer.
func (e *Engine) Resize(w, h int) {
	e.w = w
	e.h = h
	e.Frame = NewBuffer(w, h, e.Palette)
	e.Frame.Fill(e.ClearColor)
}

// NextFrame renders the next frame.
func (e *Engine) NextFrame(delta float32) {
	startTime := time.Now()
	// Prep the frame
	e.Frame.Fill(e.ClearColor)
	// Rendering process
	e.Frame.SetPixel(0, 0, 2)
	e.Frame.SetPixel(319, 239, 2)
	// Advance FPS measurement
	e.frameCount++
	endTime := time.Now()
	e.frameTime += endTime.Sub(startTime)
	d := endTime.Sub(e.lastFPSCalc)
	if d > e.fpsDelay {
		e.lastFPSCalc = endTime
		e.FPS = float32(e.frameCount) / float32(d.Seconds())
		e.LastMSPF = float32(e.frameTime.Seconds()/d.Seconds()) * 1000
		e.frameCount = 0
		e.frameTime = 0
	}
	// FPS display
}
