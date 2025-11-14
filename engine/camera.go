package engine

import (
	"log/slog"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/qbradq/eye-engine/engine/types"
)

// Camera represents a point camera in 3D space.
type Camera struct {
	types.Entity
	FOV      float32       // Field of View in degrees
	NearClip float32       // Near clipping plane distance
	FarClip  float32       // Far clipping  plane distance
	Target   *types.Buffer // Buffer we are rendering onto

	m   mgl32.Mat4 // Model matrix
	v   mgl32.Mat4 // View matrix
	p   mgl32.Mat4 // Perspective matrix
	mvp mgl32.Mat4 // Model View Perspective transform matrix
}

// Update updates the internal camera matrixes and should only be called once
// per frame.
func (c *Camera) Update() {
	c.v = c.View()
	c.p = c.Projection()
}

// SetModelMatrix sets the model matrix and updates the internal transform
// matrix.
func (c *Camera) SetModelMatrix(m mgl32.Mat4) {
	c.m = m
	c.mvp = c.p.Mul4(c.v).Mul4(c.m)
}

// View returns the view matrix.
func (c *Camera) View() mgl32.Mat4 {
	t := c.Position.Add(c.Forward)
	return mgl32.LookAt(
		c.Position[0], c.Position[1], c.Position[2],
		// c.Forward[0], c.Forward[1], c.Forward[2],
		t[0], t[1], t[2],
		c.Up[0], c.Up[1], c.Up[2],
	)
}

// Projection returns the projection matrix for the given buffer.
func (c *Camera) Projection() mgl32.Mat4 {
	return mgl32.Perspective(
		mgl32.DegToRad(c.FOV),
		float32(c.Target.Width)/float32(c.Target.Height),
		c.NearClip,
		c.FarClip,
	)
}

// Transform returns the point p transformed by the camera matrix. The x and y
// elements are in screen space. The third is the depth.
func (c *Camera) Transform(p mgl32.Vec3) mgl32.Vec3 {
	v := c.mvp.Mul4x1(mgl32.Vec4{p[0], p[1], p[2], 1.0})
	u := mgl32.Vec3{v[0] / v[3], v[1] / v[3], v[2] / v[3]}
	di := mgl32.Vec3{u[0] / u[2], u[1] / u[2], u[2]}
	return mgl32.Vec3{
		((di[0] + 1) / 2) * float32(c.Target.Width),
		float32(c.Target.Height) - ((di[1]+1)/2)*float32(c.Target.Height),
		di[2],
	}
}

// TransformToScreen returns the screen point for the given 3D position.
func (c *Camera) TransformToScreen(p mgl32.Vec3) types.PointI2D {
	p = c.Transform(p)
	return types.PointI2D{
		int(p[0]),
		int(p[1]),
	}
}

// DrawPoint draws a slice of 3D points onto the target using the color.
func (c *Camera) DrawPoints(verts []types.Vertex, color types.ColorIndex) {
	for i := range verts {
		sv := c.Transform(verts[i].Position)
		c.Target.SetPixel(int(sv[0]), int(sv[1]), color)
	}
}

// DrawLineLoop draws a line loop in 3D space onto the target using the color.
func (c *Camera) DrawLineLoop(verts []types.Vertex, color types.ColorIndex) {
	points := types.PointI2DPool.Get()
	for i := range verts {
		sv := c.Transform(verts[i].Position)
		points = append(points, types.PointI2D{int(sv[0]), int(sv[1])})
	}
	c.Target.DrawLineLoop(points, color)
	types.PointI2DPool.Release(points)
}

// DrawTriangle draws a triangle in 3D space onto the target using the color.
func (c *Camera) DrawTriangle(verts []types.Vertex, color types.ColorIndex) {
	var buf [3]types.PointI2D
	buf[0] = c.TransformToScreen(verts[0].Position)
	buf[1] = c.TransformToScreen(verts[1].Position)
	buf[2] = c.TransformToScreen(verts[2].Position)
	c.Target.DrawTriangle(buf[:], color)
}

// DrawModel draws the model onto the target using the color and draw mode.
func (c *Camera) DrawModel(m *types.Model, color types.ColorIndex, mode types.DrawMode) {
	// Prep a point buffer if needed
	var buf []types.PointI2D
	if mode == types.DrawModeFlat {
		buf = types.PointI2DPool.Get()
		defer types.PointI2DPool.Release(buf)
	}
	// Draw all faces of the model
	for i := range m.Faces {
		// Skip back-facing polygons
		if m.Faces[i].BackFacing(c.Position) {
			continue
		}
		// Draw
		switch mode {
		case types.DrawModePoints:
			c.DrawPoints(m.Faces[i].Vertexes[:], color)
		case types.DrawModeLines:
			c.DrawLineLoop(m.Faces[i].Vertexes[:], color)
		case types.DrawModeFlat:
			c.DrawTriangle(m.Faces[i].Vertexes[:], color)
		default:
			slog.Error("invalid draw mode", "mode", mode)
		}
		// color = (color + 1) % 16
	}
}
