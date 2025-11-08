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

// Vec3 represents a 3D vector
type Vec3 struct {
	X, Y, Z float32
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
