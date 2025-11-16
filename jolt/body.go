package jolt

// #include "wrapper/body.h"
import "C"

// MotionType determines how a body responds to forces
type MotionType int

const (
	MotionTypeStatic    MotionType = C.JoltMotionTypeStatic    // Immovable, zero velocity
	MotionTypeKinematic MotionType = C.JoltMotionTypeKinematic // Movable by user, doesn't respond to forces
	MotionTypeDynamic   MotionType = C.JoltMotionTypeDynamic   // Affected by forces
)

// BodyInterface provides methods to create and manipulate physics bodies
type BodyInterface struct {
	handle C.JoltBodyInterface
}

// GetBodyInterface returns the interface for creating/manipulating bodies
func (ps *PhysicsSystem) GetBodyInterface() *BodyInterface {
	handle := C.JoltPhysicsSystemGetBodyInterface(ps.handle)
	return &BodyInterface{handle: handle}
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

// SetPosition updates the position of a body
func (bi *BodyInterface) SetPosition(bodyID *BodyID, position Vec3) {
	C.JoltSetBodyPosition(
		bi.handle,
		bodyID.handle,
		C.float(position.X),
		C.float(position.Y),
		C.float(position.Z),
	)
}

// CreateBodyWithMotionType creates a body with specific motion type and sensor flag.
// This is useful for creating kinematic sensor bodies that are detected by raycasts
// but don't participate in collision resolution.
//
// Parameters:
//   - shape: The collision shape
//   - position: Initial position
//   - motionType: MotionTypeStatic, MotionTypeKinematic, or MotionTypeDynamic
//   - isSensor: If true, body is detected by queries but doesn't generate contact forces
//
// Example - Create kinematic sensor for player detection:
//
//	capsule := jolt.CreateCapsule(0.5, 1.8)
//	sensorBody := bi.CreateBodyWithMotionType(
//	    capsule,
//	    jolt.Vec3{X: 0, Y: 1, Z: 0},
//	    jolt.MotionTypeKinematic,
//	    true,  // isSensor
//	)
//	bi.ActivateBody(sensorBody)
func (bi *BodyInterface) CreateBodyWithMotionType(shape *Shape, position Vec3, motionType MotionType, isSensor bool) *BodyID {
	sensor := C.int(0)
	if isSensor {
		sensor = C.int(1)
	}

	handle := C.JoltCreateBodyWithMotionType(
		bi.handle,
		shape.handle,
		C.float(position.X),
		C.float(position.Y),
		C.float(position.Z),
		C.JoltMotionType(motionType),
		sensor,
	)

	return &BodyID{handle: handle}
}

// ActivateBody makes a body participate in the simulation
func (bi *BodyInterface) ActivateBody(bodyID *BodyID) {
	C.JoltActivateBody(bi.handle, bodyID.handle)
}

// DeactivateBody removes a body from active simulation
func (bi *BodyInterface) DeactivateBody(bodyID *BodyID) {
	C.JoltDeactivateBody(bi.handle, bodyID.handle)
}
