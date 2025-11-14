package types

import "github.com/go-gl/mathgl/mgl32"

// Up is the normal up vector
var Up = mgl32.Vec3{0, 1, 0}

// Down points down
var Down = mgl32.Vec3{0, -1, 0}

// Forward points down the negative Z axis
var Forward = mgl32.Vec3{0, 0, -1}

// Backward points down the positive Z axis
var Backward = mgl32.Vec3{0, 0, 1}

// Left points down the negative X axis
var Left = mgl32.Vec3{-1, 0, 0}

// Right points down the positive X axis
var Right = mgl32.Vec3{1, 0, 0}

// DrawMode informs the engine how to draw objects.
type DrawMode uint8

// Draw modes.
const (
	DrawModePoints DrawMode = iota
	DrawModeLines
	DrawModeFlat
)

// PointI2D represents a point in 2D space with integers (x, y)
type PointI2D [2]int

// RectI2D represents a 2D rectangular area with integers (x, y, w, h)
type RectI2D [4]int

// OutCode tells where a point is relative to a rect.
type OutCode uint8

const (
	OutCodeTop    OutCode = 0b0001
	OutCodeBottom OutCode = 0b0010
	OutCodeRight  OutCode = 0b0100
	OutCodeLeft   OutCode = 0b1000
)

// OutCode calculates the OutCode of the point relative to r.
func (r RectI2D) OutCode(p PointI2D) OutCode {
	var out OutCode
	if p[1] < r[1] {
		out |= OutCodeTop
	}
	if p[1] >= r[1]+r[3] {
		out |= OutCodeBottom
	}
	if p[0] >= r[0]+r[2] {
		out |= OutCodeRight
	}
	if p[0] < r[0] {
		out |= OutCodeLeft
	}
	return out
}
