package jolt

// #include "wrapper/jolt_wrapper.h"
import "C"

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

// CreateCapsule creates a capsule body (cylinder with hemispherical caps)
// halfHeight: half-height of the cylindrical part (not including the caps)
// radius: radius of the capsule
// isDynamic: true = affected by forces, false = static/immovable
func (bi *BodyInterface) CreateCapsule(halfHeight, radius float32, position Vec3, isDynamic bool) *BodyID {
	dynamic := C.int(0)
	if isDynamic {
		dynamic = C.int(1)
	}

	handle := C.JoltCreateCapsule(
		bi.handle,
		C.float(halfHeight),
		C.float(radius),
		C.float(position.X),
		C.float(position.Y),
		C.float(position.Z),
		dynamic,
	)

	return &BodyID{handle: handle}
}

// CreateConvexHull creates a convex hull body from a set of points
// points: slice of Vec3 vertices that define the convex hull
// isDynamic: true = affected by forces, false = static/immovable
func (bi *BodyInterface) CreateConvexHull(points []Vec3, position Vec3, isDynamic bool) *BodyID {
	dynamic := C.int(0)
	if isDynamic {
		dynamic = C.int(1)
	}

	// Flatten Vec3 slice to float array
	floatPoints := make([]C.float, len(points)*3)
	for i, p := range points {
		floatPoints[i*3] = C.float(p.X)
		floatPoints[i*3+1] = C.float(p.Y)
		floatPoints[i*3+2] = C.float(p.Z)
	}

	handle := C.JoltCreateConvexHull(
		bi.handle,
		&floatPoints[0],
		C.int(len(points)),
		C.float(position.X),
		C.float(position.Y),
		C.float(position.Z),
		dynamic,
	)

	return &BodyID{handle: handle}
}

// CreateMesh creates a mesh body from vertices and triangle indices
// vertices: slice of Vec3 vertices
// indices: slice of triangle indices (must be multiple of 3, each triangle is 3 indices)
// isDynamic: true = affected by forces, false = static/immovable
// Note: Mesh shapes are typically used for static geometry (e.g., terrain, buildings)
func (bi *BodyInterface) CreateMesh(vertices []Vec3, indices []int32, position Vec3, isDynamic bool) *BodyID {
	dynamic := C.int(0)
	if isDynamic {
		dynamic = C.int(1)
	}

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
		bi.handle,
		&floatVertices[0],
		C.int(len(vertices)),
		&cIndices[0],
		C.int(len(indices)),
		C.float(position.X),
		C.float(position.Y),
		C.float(position.Z),
		dynamic,
	)

	return &BodyID{handle: handle}
}
