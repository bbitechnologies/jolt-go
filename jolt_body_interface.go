package jolt

// #include "wrapper/jolt_wrapper.h"
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

// CreateSphere creates a sphere body
// isDynamic: true = affected by forces, false = static/immovable
func (bi *BodyInterface) CreateSphere(radius float32, position Vec3, isDynamic bool) *BodyID {
	dynamic := C.int(0)
	if isDynamic {
		dynamic = C.int(1)
	}

	handle := C.JoltCreateSphere(
		bi.handle,
		C.float(radius),
		C.float(position.X),
		C.float(position.Y),
		C.float(position.Z),
		dynamic,
	)

	return &BodyID{handle: handle}
}

// CreateBox creates a box body
// halfExtent: half-size of the box in each dimension (e.g., Vec3{5,0.5,5} creates 10x1x10 box)
// isDynamic: true = affected by forces, false = static/immovable
func (bi *BodyInterface) CreateBox(halfExtent Vec3, position Vec3, isDynamic bool) *BodyID {
	dynamic := C.int(0)
	if isDynamic {
		dynamic = C.int(1)
	}

	handle := C.JoltCreateBox(
		bi.handle,
		C.float(halfExtent.X),
		C.float(halfExtent.Y),
		C.float(halfExtent.Z),
		C.float(position.X),
		C.float(position.Y),
		C.float(position.Z),
		dynamic,
	)

	return &BodyID{handle: handle}
}
