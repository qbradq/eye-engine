package types

import "github.com/go-gl/mathgl/mgl32"

// Vertex holds all of the data for a single vertex.
type Vertex struct {
	Position mgl32.Vec3 // Position
	Normal   mgl32.Vec3 // Normal
	UV       mgl32.Vec2 // Texture UV
}

// Face holds all of the data of a face.
type Face struct {
	Vertexes []Vertex // The vertexes of the face in counter-clockwise order
}

// Model holds all of the data of a model.
type Model struct {
	Faces []Face // All of the faces of the model
}
