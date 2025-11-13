package types

import "github.com/go-gl/mathgl/mgl32"

// Entity is the basic object type for objects within the engine's game world.
type Entity struct {
	Position mgl32.Vec3 // Absolute position
	Up       mgl32.Vec3 // Up vector for this entity
	Forward  mgl32.Vec3 // Unit vector showing the orientation of the entity
}

// MoveForward moves the entity in the direction of it's forward vector.
func (e *Entity) MoveForward(d float32) {
	e.Position = e.Position.Add(e.Forward.Mul(d))
}

// MoveBackward moves the entity in the direction of it's forward vector.
func (e *Entity) MoveBackward(d float32) {
	e.Position = e.Position.Add(e.Forward.Mul(-d))
}

// StrafeRight strafes the entity right relative to the forward vector.
func (e *Entity) StrafeRight(d float32) {
	r := e.Forward.Cross(e.Up).Normalize()
	e.Position = e.Position.Add(r.Mul(d))
}

// StrafeLeft strafes the entity left relative to the forward vector.
func (e *Entity) StrafeLeft(d float32) {
	r := e.Forward.Cross(e.Up).Normalize()
	e.Position = e.Position.Add(r.Mul(-d))
}

// Rotate rotates the entity along the Y axis (yaw) by d degrees.
func (e *Entity) Rotate(d float32) {
	r := mgl32.DegToRad(d)
	m := mgl32.HomogRotate3D(r, e.Up)
	f := m.Mul4x1(e.Forward.Vec4(0))
	e.Forward = f.Vec3().Normalize()
}
