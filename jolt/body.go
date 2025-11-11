package jolt

// #include "wrapper/body.h"
import "C"

// BodyInterface provides methods to create and manipulate physics bodies
type BodyInterface struct {
	handle C.JoltBodyInterface
}

// BodyID uniquely identifies a physics body
type BodyID struct {
	handle C.JoltBodyID
}

// Destroy frees the body ID
func (b *BodyID) Destroy() {
	C.JoltDestroyBodyID(b.handle)
}

// GetPosition returns the current position of a body
func (bi *BodyInterface) GetPosition(bodyID *BodyID) Vec3 {
	var x, y, z C.float
	C.JoltGetBodyPosition(bi.handle, bodyID.handle, &x, &y, &z)
	return Vec3{
		X: float32(x),
		Y: float32(y),
		Z: float32(z),
	}
}

// CreateBody creates a body from a shape
// shape: the collision shape
// position: initial position
// isDynamic: true = affected by forces, false = static/immovable
func (bi *BodyInterface) CreateBody(shape *Shape, position Vec3, isDynamic bool) *BodyID {
	dynamic := C.int(0)
	if isDynamic {
		dynamic = C.int(1)
	}

	handle := C.JoltCreateBody(
		bi.handle,
		shape.handle,
		C.float(position.X),
		C.float(position.Y),
		C.float(position.Z),
		dynamic,
	)

	return &BodyID{handle: handle}
}

// CreateStaticBody creates a static (immovable) body from a shape
func (bi *BodyInterface) CreateStaticBody(shape *Shape, position Vec3) *BodyID {
	return bi.CreateBody(shape, position, false)
}
