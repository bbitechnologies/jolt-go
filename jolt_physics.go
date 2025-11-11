package jolt

// #include "wrapper/jolt_wrapper.h"
import "C"

// PhysicsSystem represents a physics simulation world
type PhysicsSystem struct {
	handle C.JoltPhysicsSystem
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
