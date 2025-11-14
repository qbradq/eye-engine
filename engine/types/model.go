package types

import "github.com/go-gl/mathgl/mgl32"

// Vertex holds all of the data for a single vertex.
type Vertex struct {
	Position mgl32.Vec3 // Position
	Normal   mgl32.Vec3 // Normal
	UV       mgl32.Vec2 // Texture UV
}

// Face holds all of the data of a triangular face.
type Face struct {
	Vertexes [3]Vertex  // The vertexes of the face in counter-clockwise order
	Normal   mgl32.Vec3 // Face normal
}

// BackFacing returns true if the face is back-facing relative to c.
func (f *Face) BackFacing(c mgl32.Vec3) bool {
	v := c.Sub(f.Vertexes[0].Position)
	return f.Normal.Dot(v) <= 0
}

// Model holds all of the data of a model.
type Model struct {
	Faces []Face // All of the faces of the model
}

// UpdateGeometry recalculates the internal geometry caches. Must be called
// after altering the face contents.
func (m *Model) UpdateGeometry() {
	// Calculate face normals
	for iFace := range m.Faces {
		f := &m.Faces[iFace]
		ab := f.Vertexes[1].Position.Sub(f.Vertexes[0].Position)
		ac := f.Vertexes[2].Position.Sub(f.Vertexes[0].Position)
		f.Normal = ab.Cross(ac)
	}
}
