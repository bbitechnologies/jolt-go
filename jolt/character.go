package jolt

// #include "wrapper/character.h"
import "C"

// CharacterVirtual represents a virtual character in the physics world
type CharacterVirtual struct {
	handle C.JoltCharacterVirtual
	ps     *PhysicsSystem
}

// GroundState indicates the ground contact state of a CharacterVirtual
type GroundState int

const (
	// GroundStateOnGround - Character is on the ground and can move freely
	GroundStateOnGround GroundState = 0
	// GroundStateOnSteepGround - Character is on a slope too steep to climb
	GroundStateOnSteepGround GroundState = 1
	// GroundStateNotSupported - Character is touching an object but not supported (should fall)
	GroundStateNotSupported GroundState = 2
	// GroundStateInAir - Character is in the air and not touching anything
	GroundStateInAir GroundState = 3
)

func (gs GroundState) String() string {
	switch gs {
	case GroundStateOnGround:
		return "OnGround"
	case GroundStateOnSteepGround:
		return "OnSteepGround"
	case GroundStateNotSupported:
		return "NotSupported"
	case GroundStateInAir:
		return "InAir"
	default:
		return "Unknown"
	}
}

// CreateCharacterVirtual creates a virtual character with the specified shape at the initial position
func (ps *PhysicsSystem) CreateCharacterVirtual(shape *Shape, position Vec3) *CharacterVirtual {
	handle := C.JoltCreateCharacterVirtual(
		ps.handle,
		shape.handle,
		C.float(position.X),
		C.float(position.Y),
		C.float(position.Z),
	)
	return &CharacterVirtual{handle: handle, ps: ps}
}

// ExtendedUpdate advances the character simulation with combined movement logic
// Combines Update, StickToFloor, and WalkStairs into a unified operation
// deltaTime: duration of simulation step in seconds
// gravity: acceleration vector (e.g., Vec3{0, -9.81, 0} for Earth gravity)
func (cv *CharacterVirtual) ExtendedUpdate(deltaTime float32, gravity Vec3) {
	C.JoltCharacterVirtualExtendedUpdate(
		cv.handle,
		cv.ps.handle,
		C.float(deltaTime),
		C.float(gravity.X),
		C.float(gravity.Y),
		C.float(gravity.Z),
	)
}

// SetLinearVelocity sets the character's linear velocity
func (cv *CharacterVirtual) SetLinearVelocity(velocity Vec3) {
	C.JoltCharacterVirtualSetLinearVelocity(
		cv.handle,
		C.float(velocity.X),
		C.float(velocity.Y),
		C.float(velocity.Z),
	)
}

// GetLinearVelocity returns the current linear velocity of the character
func (cv *CharacterVirtual) GetLinearVelocity() Vec3 {
	var x, y, z C.float
	C.JoltCharacterVirtualGetLinearVelocity(cv.handle, &x, &y, &z)
	return Vec3{
		X: float32(x),
		Y: float32(y),
		Z: float32(z),
	}
}

// GetGroundVelocity returns the velocity clamped to the ground plane
func (cv *CharacterVirtual) GetGroundVelocity() Vec3 {
	var x, y, z C.float
	C.JoltCharacterVirtualGetGroundVelocity(cv.handle, &x, &y, &z)
	return Vec3{
		X: float32(x),
		Y: float32(y),
		Z: float32(z),
	}
}

// SetPosition sets the character's position in the world
func (cv *CharacterVirtual) SetPosition(position Vec3) {
	C.JoltCharacterVirtualSetPosition(
		cv.handle,
		C.float(position.X),
		C.float(position.Y),
		C.float(position.Z),
	)
}

// GetPosition returns the current position of the character
func (cv *CharacterVirtual) GetPosition() Vec3 {
	var x, y, z C.float
	C.JoltCharacterVirtualGetPosition(cv.handle, &x, &y, &z)
	return Vec3{
		X: float32(x),
		Y: float32(y),
		Z: float32(z),
	}
}

// Destroy frees the character resources
func (cv *CharacterVirtual) Destroy() {
	C.JoltDestroyCharacterVirtual(cv.handle)
}

// GetGroundState returns the current ground contact state
func (cv *CharacterVirtual) GetGroundState() GroundState {
	state := C.JoltCharacterVirtualGetGroundState(cv.handle)
	return GroundState(state)
}

// IsSupported returns true if the character is on ground or steep ground (not falling)
func (cv *CharacterVirtual) IsSupported() bool {
	result := C.JoltCharacterVirtualIsSupported(cv.handle)
	return result != 0
}

// SetShape changes the collision shape of the character
// shape: new collision shape for the character
// maxPenetrationDepth: maximum allowed penetration depth (typically 0.1)
func (cv *CharacterVirtual) SetShape(shape *Shape, maxPenetrationDepth float32) {
	C.JoltCharacterVirtualSetShape(
		cv.handle,
		shape.handle,
		C.float(maxPenetrationDepth),
		cv.ps.handle,
	)
}

// GetShape retrieves the current collision shape of the character
func (cv *CharacterVirtual) GetShape() *Shape {
	handle := C.JoltCharacterVirtualGetShape(cv.handle)
	return &Shape{handle: handle}
}

// PhysicsSystem returns the physics system that this character belongs to
func (cv *CharacterVirtual) PhysicsSystem() *PhysicsSystem {
	return cv.ps
}

// GetGroundNormal returns the normal vector of the ground surface
func (cv *CharacterVirtual) GetGroundNormal() Vec3 {
	var x, y, z C.float
	C.JoltCharacterVirtualGetGroundNormal(cv.handle, &x, &y, &z)
	return Vec3{
		X: float32(x),
		Y: float32(y),
		Z: float32(z),
	}
}
