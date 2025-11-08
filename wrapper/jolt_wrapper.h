/*
 * Jolt Physics C Wrapper Header
 *
 * Provides a C interface to Jolt Physics C++ library for use with Go's CGO.
 * Uses opaque pointers (void*) to hide C++ types from the C interface.
 */

#ifndef JOLT_WRAPPER_H
#define JOLT_WRAPPER_H

#ifdef __cplusplus
extern "C" {
#endif

// Opaque pointer types representing C++ objects
typedef void* JoltPhysicsSystem;
typedef void* JoltBodyInterface;
typedef void* JoltBodyID;

// Initialize Jolt Physics (call once at startup)
// Returns 1 on success, 0 on failure
int JoltInit();

// Shutdown Jolt Physics (call once at exit)
void JoltShutdown();

// Create a new physics world
JoltPhysicsSystem JoltCreatePhysicsSystem();

// Destroy a physics world
void JoltDestroyPhysicsSystem(JoltPhysicsSystem system);

// Step the physics simulation by deltaTime seconds
void JoltPhysicsSystemUpdate(JoltPhysicsSystem system, float deltaTime);

// Get the body interface for creating/manipulating bodies
JoltBodyInterface JoltPhysicsSystemGetBodyInterface(JoltPhysicsSystem system);

// Create a sphere body
// isDynamic: 1 = dynamic (affected by forces), 0 = static (immovable)
JoltBodyID JoltCreateSphere(JoltBodyInterface bodyInterface,
                            float radius,
                            float x, float y, float z,
                            int isDynamic);

// Get the position of a body
void JoltGetBodyPosition(JoltBodyInterface bodyInterface,
                        JoltBodyID bodyID,
                        float* x, float* y, float* z);

#ifdef __cplusplus
}
#endif

#endif // JOLT_WRAPPER_H
