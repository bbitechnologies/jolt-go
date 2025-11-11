/*
 * Jolt Physics C Wrapper - Character Virtual
 *
 * Handles virtual character controller for player/NPC movement.
 */

#ifndef JOLT_WRAPPER_CHARACTER_H
#define JOLT_WRAPPER_CHARACTER_H

#ifdef __cplusplus
extern "C" {
#endif

// Opaque pointer types
typedef void* JoltCharacterVirtual;
typedef void* JoltPhysicsSystem;

// Ground state enum (matches Jolt's EGroundState)
typedef enum {
    JoltGroundStateOnGround = 0,       // Character is on the ground and can move freely
    JoltGroundStateOnSteepGround = 1,  // On a slope too steep to climb
    JoltGroundStateNotSupported = 2,   // Touching object but not supported (should fall)
    JoltGroundStateInAir = 3           // In the air, not touching anything
} JoltGroundState;

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

#endif // JOLT_WRAPPER_CHARACTER_H
