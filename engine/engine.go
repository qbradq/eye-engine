package engine

import (
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"sync"
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
	pointPool   sync.Pool     // Pool of PointI2D slices
}

// NewEngine returns a new Engine ready to use.
func NewEngine(w, h int) (*Engine, error) {
	var err error
	ret := &Engine{
		ClearColor:  15,
		fpsDelay:    time.Millisecond * 500,
		lastFPSCalc: time.Now(),
		pointPool: sync.Pool{
			New: func() any {
				b := make([]PointI2D, 0, 512)
				return &b
			},
		},
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
	e.DrawLine(e.Frame, PointI2D{0, 0}, PointI2D{200, 239}, DefaultCyan)
	e.DrawTriangle(
		e.Frame,
		PointI2D{15, 60},
		PointI2D{45, 120},
		PointI2D{60, 90},
		DefaultEvergreen,
	)
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

// ReleasePointPoolSlice returns s to the available memory pool.
func (e *Engine) ReleasePointPoolSlice(s []PointI2D) {
	sPtr := &s
	*sPtr = (*sPtr)[:0]
	e.pointPool.Put(sPtr)
}

// Line returns a slice of ints representing the points along a line from p0 to
// p1 in 2D space. When done with the returned slice, release it with
// ReleasePointPoolSlice.
func (e *Engine) Line(p0, p1 PointI2D) []PointI2D {
	// Allocate return slice
	retPtr, ok := e.pointPool.Get().(*[]PointI2D)
	if !ok {
		panic(errors.New("e.pointPool.Get() did not return *[]PointI2D"))
	}
	ret := *retPtr
	dx := p1[0] - p0[0]
	dy := p1[1] - p0[1]
	stepX := 1
	if dx < 0 {
		stepX = -1
	}
	stepY := 1
	if dy < 0 {
		stepY = -1
	}
	absDx := dx
	if absDx < 0 {
		absDx = -absDx
	}
	absDy := dy
	if absDy < 0 {
		absDy = -absDy
	}
	if absDx >= absDy {
		// X-major axis
		ret = append(ret, p0)
		p := 2*absDy - absDx
		for range absDx {
			if p >= 0 {
				p0[1] += stepY
				p += 2 * (absDy - absDx)
			} else {
				p += 2 * absDy
			}
			p0[0] += stepX
			ret = append(ret, p0)
		}
	} else {
		// Y-major axis
		ret = append(ret, p0)
		p := 2*absDx - absDy
		for range absDy {
			if p >= 0 {
				p0[0] += stepX
				p += 2 * (absDx - absDy)
			} else {
				p += 2 * absDx
			}
			p0[1] += stepY
			ret = append(ret, p0)
		}
	}
	return ret
}

// DrawLine draws a line on buf from p0 to p1 using color c.
func (e *Engine) DrawLine(buf *Buffer, p0, p1 PointI2D, c ColorIndex) {
	points := e.Line(p0, p1)
	for _, p := range points {
		buf.Pixels[p[1]*buf.Width+p[0]] = c
	}
	e.ReleasePointPoolSlice(points)
}

// TriangleScanLines returns a slice of an even number of points. For each pair
// of points, they specify the left and right limits of a scan line contained
// within the triangle.
func (e *Engine) TriangleScanLines(p0, p1, p2 PointI2D) []PointI2D {
	// Allocate return slice
	retPtr, ok := e.pointPool.Get().(*[]PointI2D)
	if !ok {
		panic(errors.New("e.pointPool.Get() did not return *[]PointI2D"))
	}
	ret := *retPtr
	// Sort points top to bottom
	p := []PointI2D{p0, p1, p2}
	sort.SliceStable(p, func(i, j int) bool {
		return p[i][1] < p[j][1]
	})
	p0 = p[0]
	p1 = p[1]
	p2 = p[2]
	// Setup the interpolation loops
	s0 := float32(p2[0]-p1[0]) / float32(p2[1]-p1[1])
	s1 := float32(p1[0]-p0[0]) / float32(p1[1]-p0[1])
	s2 := float32(p2[0]-p0[0]) / float32(p2[1]-p0[1])
	xLeft := float32(p0[0])
	xRight := float32(p0[0])
	// Top interpolation loop
	lSlope := s1
	rSlope := s2
	if s1 > s2 {
		lSlope = s2
		rSlope = s1
	}
	for y := p0[1]; y < p1[1]; y++ {
		ret = append(ret, PointI2D{int(xLeft), y}, PointI2D{int(xRight), y})
		xLeft += lSlope
		xRight += rSlope
	}
	// Bottom interpolation loop
	if lSlope == s2 {
		rSlope = s0
	} else {
		lSlope = s0
	}
	for y := p1[1]; y <= p2[1]; y++ {
		ret = append(ret, PointI2D{int(xLeft), y}, PointI2D{int(xRight), y})
		xLeft += lSlope
		xRight += rSlope
	}
	return ret
}

// DrawTriangle draws triangle p0, p1, p2 on buf using color c.
func (e *Engine) DrawTriangle(buf *Buffer, p0, p1, p2 PointI2D, c ColorIndex) {
	points := e.TriangleScanLines(p0, p1, p2)
	for i := 0; i < len(points); i += 2 {
		pl := points[i+0]
		pr := points[i+1]
		for x := pl[0]; x <= pr[0]; x++ {
			buf.Pixels[pl[1]*buf.Width+x] = c
		}
	}
	e.ReleasePointPoolSlice(points)
}
