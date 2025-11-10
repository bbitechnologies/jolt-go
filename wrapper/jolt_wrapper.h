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
typedef void* JoltCharacterVirtual;

// Ground state enum (matches Jolt's EGroundState)
typedef enum {
    JoltGroundStateOnGround = 0,       // Character is on the ground and can move freely
    JoltGroundStateOnSteepGround = 1,  // On a slope too steep to climb
    JoltGroundStateNotSupported = 2,   // Touching object but not supported (should fall)
    JoltGroundStateInAir = 3           // In the air, not touching anything
} JoltGroundState;

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

// Get the position of a body
void JoltGetBodyPosition(const JoltBodyInterface bodyInterface,
                        const JoltBodyID bodyID,
                        float* x, float* y, float* z);

// Create a sphere body
// isDynamic: 1 = dynamic (affected by forces), 0 = static (immovable)
JoltBodyID JoltCreateSphere(JoltBodyInterface bodyInterface,
                            float radius,
                            float x, float y, float z,
                            int isDynamic);

// Create a box body
// halfExtentX/Y/Z: half-size of the box in each dimension
// isDynamic: 1 = dynamic (affected by forces), 0 = static (immovable)
JoltBodyID JoltCreateBox(JoltBodyInterface bodyInterface,
                         float halfExtentX, float halfExtentY, float halfExtentZ,
                         float x, float y, float z,
                         int isDynamic);

// Destroy a body ID
void JoltDestroyBodyID(JoltBodyID bodyID);

// Create a new virtual character at initial position (x, y, z)
JoltCharacterVirtual JoltCreateCharacterVirtual(JoltPhysicsSystem system,
                                              float x, float y, float z);

// Destroy a virtual character
void JoltDestroyCharacterVirtual(JoltCharacterVirtual character);

// Update virtual character with extended update (combines Update, StickToFloor, WalkStairs)
// gravityX/Y/Z: gravity vector applied when character stands on another object
void JoltCharacterVirtualExtendedUpdate(JoltCharacterVirtual character,
                                        JoltPhysicsSystem system,
                                        float deltaTime,
                                        float gravityX, float gravityY, float gravityZ);

// Set the linear velocity of a virtual character
void JoltCharacterVirtualSetLinearVelocity(JoltCharacterVirtual character,
                                           float x, float y, float z);

// Get the linear velocity of a virtual character
void JoltCharacterVirtualGetLinearVelocity(const JoltCharacterVirtual character,
                                           float* x, float* y, float* z);

// Get the ground velocity of a virtual character
void JoltCharacterVirtualGetGroundVelocity(const JoltCharacterVirtual character,
                                           float* x, float* y, float* z);

// Get the position of a virtual character
void JoltCharacterVirtualGetPosition(const JoltCharacterVirtual character,
                                     float* x, float* y, float* z);

// Get the ground state of a virtual character
JoltGroundState JoltCharacterVirtualGetGroundState(const JoltCharacterVirtual character);

// Check if character is supported (on ground or steep ground)
int JoltCharacterVirtualIsSupported(const JoltCharacterVirtual character);

#ifdef __cplusplus
}
#endif

#endif // JOLT_WRAPPER_H
