/*
Package jolt provides Go bindings for the Jolt Physics engine.

Architecture: Go code → CGO → C wrapper → Jolt C++ library
Uses opaque pointers to hide C++ types from Go.

Pre-built binaries are provided for darwin/arm64 and linux/amd64.
Platform-specific LDFLAGS are in jolt_<os>_<arch>.go files.
*/
package jolt

/*
#cgo CXXFLAGS: -std=c++17 -IJoltPhysics -DNDEBUG -DJPH_DISABLE_CUSTOM_ALLOCATOR -DJPH_PROFILE_ENABLED -DJPH_DEBUG_RENDERER -DJPH_OBJECT_STREAM
#include "wrapper/jolt_wrapper.h"
*/
import "C"
import (
	"fmt"
)

// PhysicsSystem represents a physics simulation world
type PhysicsSystem struct {
	handle C.JoltPhysicsSystem
}

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

// Vec3 represents a 3D vector
type Vec3 struct {
	X, Y, Z float32
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

// Init initializes Jolt Physics (call once at startup)
func Init() error {
	result := C.JoltInit()
	if result == 0 {
		return fmt.Errorf("failed to initialize Jolt")
	}
	return nil
}

// Shutdown cleans up Jolt resources (call once at exit)
func Shutdown() {
	C.JoltShutdown()
}

// NewPhysicsSystem creates a new physics world
func NewPhysicsSystem() *PhysicsSystem {
	handle := C.JoltCreatePhysicsSystem()
	return &PhysicsSystem{handle: handle}
}

// Destroy frees the physics system
func (ps *PhysicsSystem) Destroy() {
	C.JoltDestroyPhysicsSystem(ps.handle)
}

// Update advances the simulation by deltaTime seconds
func (ps *PhysicsSystem) Update(deltaTime float32) {
	C.JoltPhysicsSystemUpdate(ps.handle, C.float(deltaTime))
}

// GetBodyInterface returns the interface for creating/manipulating bodies
func (ps *PhysicsSystem) GetBodyInterface() *BodyInterface {
	handle := C.JoltPhysicsSystemGetBodyInterface(ps.handle)
	return &BodyInterface{handle: handle}
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

// CreateCharacterVirtual creates a virtual character at the specified initial position
func (ps *PhysicsSystem) CreateCharacterVirtual(position Vec3) *CharacterVirtual {
	handle := C.JoltCreateCharacterVirtual(
		ps.handle,
		C.float(position.X),
		C.float(position.Y),
		C.float(position.Z),
	)
	return &CharacterVirtual{handle: handle, ps: ps}
}

// CharacterVirtual represents a virtual character in the physics world
type CharacterVirtual struct {
	handle C.JoltCharacterVirtual
	ps     *PhysicsSystem
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
