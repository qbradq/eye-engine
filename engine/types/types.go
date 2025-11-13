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
)

// PointI2D represents a point in 2D space with integers.
type PointI2D [2]int
