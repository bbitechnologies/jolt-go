/*
 * Jolt Physics C Wrapper - Body Operations
 *
 * Handles rigid body creation and manipulation.
 */

#ifndef JOLT_WRAPPER_BODY_H
#define JOLT_WRAPPER_BODY_H

#ifdef __cplusplus
extern "C" {
#endif

// Opaque pointer types
typedef void* JoltBodyInterface;
typedef void* JoltBodyID;
typedef void* JoltShape;

// Get the position of a body
void JoltGetBodyPosition(const JoltBodyInterface bodyInterface,
                        const JoltBodyID bodyID,
                        float* x, float* y, float* z);

// Create a body from a shape
// isDynamic: 1 = dynamic (affected by forces), 0 = static (immovable)
JoltBodyID JoltCreateBody(JoltBodyInterface bodyInterface,
                          JoltShape shape,
                          float x, float y, float z,
                          int isDynamic);

// Destroy a body ID
void JoltDestroyBodyID(JoltBodyID bodyID);

#ifdef __cplusplus
}
#endif

#endif // JOLT_WRAPPER_BODY_H
