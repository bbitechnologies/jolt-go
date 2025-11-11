package jolt

// #include "wrapper/query.h"
import "C"

// CollisionHit contains information about a single collision detected during a shape query
type CollisionHit struct {
	BodyID          *BodyID // The body that was hit
	ContactPoint    Vec3    // The contact point in world space
	PenetrationDepth float32 // How deep the shapes overlap (negative if separated)
}

// CollideShape checks if a shape at the given position collides with any bodies in the physics system.
// This performs a static overlap test - the shape itself is not added to the physics system.
//
// Parameters:
//   - shape: The collision shape to test
//   - position: Position in world space to place the shape
//
// Returns true if any collision is detected, false otherwise.
//
// Example usage:
//   sphere := jolt.CreateSphere(1.0)
//   defer sphere.Destroy()
//
//   if ps.CollideShape(sphere, jolt.Vec3{X: 0, Y: 5, Z: 0}) {
//       fmt.Println("Collision detected!")
//   }
func (ps *PhysicsSystem) CollideShape(shape *Shape, position Vec3) bool {
	result := C.JoltCollideShape(
		ps.handle,
		shape.handle,
		C.float(position.X),
		C.float(position.Y),
		C.float(position.Z),
	)
	return result != 0
}

// CollideShapeGetHits performs a shape collision query and returns detailed information about all hits.
// This is useful when you need to know which specific bodies were hit and where.
//
// Parameters:
//   - shape: The collision shape to test
//   - position: Position in world space to place the shape
//   - maxHits: Maximum number of hits to return (limits memory allocation)
//
// Returns a slice of CollisionHit containing information about each collision.
//
// Example usage:
//   sphere := jolt.CreateSphere(1.0)
//   defer sphere.Destroy()
//
//   hits := ps.CollideShapeGetHits(sphere, jolt.Vec3{X: 0, Y: 5, Z: 0}, 10)
//   for _, hit := range hits {
//       fmt.Printf("Hit body at %.2f, %.2f, %.2f (depth: %.2f)\n",
//           hit.ContactPoint.X, hit.ContactPoint.Y, hit.ContactPoint.Z,
//           hit.PenetrationDepth)
//   }
func (ps *PhysicsSystem) CollideShapeGetHits(shape *Shape, position Vec3, maxHits int) []CollisionHit {
	if maxHits <= 0 {
		return []CollisionHit{}
	}

	// Allocate C array for results
	cHits := make([]C.JoltCollisionHit, maxHits)

	numHits := C.JoltCollideShapeGetHits(
		ps.handle,
		shape.handle,
		C.float(position.X),
		C.float(position.Y),
		C.float(position.Z),
		&cHits[0],
		C.int(maxHits),
	)

	// Convert C results to Go
	hits := make([]CollisionHit, int(numHits))
	for i := 0; i < int(numHits); i++ {
		cHit := cHits[i]
		hits[i] = CollisionHit{
			BodyID: &BodyID{handle: cHit.bodyID},
			ContactPoint: Vec3{
				X: float32(cHit.contactPointX),
				Y: float32(cHit.contactPointY),
				Z: float32(cHit.contactPointZ),
			},
			PenetrationDepth: float32(cHit.penetrationDepth),
		}
	}

	return hits
}
