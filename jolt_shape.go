package jolt

// #include "wrapper/jolt_wrapper.h"
import "C"

// Shape represents collision geometry that can be used to create bodies
type Shape struct {
	handle C.JoltShape
}

// Destroy frees the shape (decrements ref count)
func (s *Shape) Destroy() {
	C.JoltDestroyShape(s.handle)
}

// CreateSphereShape creates a sphere collision shape
func CreateSphere(radius float32) *Shape {
	handle := C.JoltCreateSphere(C.float(radius))
	return &Shape{handle: handle}
}

// CreateBoxShape creates a box collision shape
// halfExtent: half-size of the box in each dimension (e.g., Vec3{5,0.5,5} creates 10x1x10 box)
func CreateBox(halfExtent Vec3) *Shape {
	handle := C.JoltCreateBox(
		C.float(halfExtent.X),
		C.float(halfExtent.Y),
		C.float(halfExtent.Z),
	)
	return &Shape{handle: handle}
}

// CreateCapsuleShape creates a capsule collision shape (cylinder with hemispherical caps)
// halfHeight: half-height of the cylindrical part (not including the caps)
// radius: radius of the capsule
func CreateCapsule(halfHeight, radius float32) *Shape {
	handle := C.JoltCreateCapsule(
		C.float(halfHeight),
		C.float(radius),
	)
	return &Shape{handle: handle}
}

// CreateConvexHullShape creates a convex hull collision shape from a set of points
// points: slice of Vec3 vertices that define the convex hull
func CreateConvexHull(points []Vec3) *Shape {
	// Flatten Vec3 slice to float array
	floatPoints := make([]C.float, len(points)*3)
	for i, p := range points {
		floatPoints[i*3] = C.float(p.X)
		floatPoints[i*3+1] = C.float(p.Y)
		floatPoints[i*3+2] = C.float(p.Z)
	}

	handle := C.JoltCreateConvexHull(
		&floatPoints[0],
		C.int(len(points)),
	)
	return &Shape{handle: handle}
}

// CreateMeshShape creates a mesh collision shape from vertices and triangle indices
// vertices: slice of Vec3 vertices
// indices: slice of triangle indices (must be multiple of 3, each triangle is 3 indices)
// Note: Mesh shapes are typically used for static geometry (e.g., terrain, buildings)
func CreateMesh(vertices []Vec3, indices []int32) *Shape {
	// Flatten Vec3 slice to float array
	floatVertices := make([]C.float, len(vertices)*3)
	for i, v := range vertices {
		floatVertices[i*3] = C.float(v.X)
		floatVertices[i*3+1] = C.float(v.Y)
		floatVertices[i*3+2] = C.float(v.Z)
	}

	// Convert int32 slice to C int array
	cIndices := make([]C.int, len(indices))
	for i, idx := range indices {
		cIndices[i] = C.int(idx)
	}

	handle := C.JoltCreateMesh(
		&floatVertices[0],
		C.int(len(vertices)),
		&cIndices[0],
		C.int(len(indices)),
	)
	return &Shape{handle: handle}
}
