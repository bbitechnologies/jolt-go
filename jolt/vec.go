package jolt

// Vec3 represents a 3D vector
type Vec3 struct {
	X, Y, Z float32
}

// Dot returns the dot product of this vector with another vector
func (v Vec3) Dot(other Vec3) float32 {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z
}
