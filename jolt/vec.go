package jolt

// Vec3 represents a 3D vector
type Vec3 struct {
	X, Y, Z float32
}

// Add returns the sum of this vector and another vector
func (v Vec3) Add(other Vec3) Vec3 {
	return Vec3{X: v.X + other.X, Y: v.Y + other.Y, Z: v.Z + other.Z}
}

// Sub returns the difference of this vector and another vector
func (v Vec3) Sub(other Vec3) Vec3 {
	return Vec3{X: v.X - other.X, Y: v.Y - other.Y, Z: v.Z - other.Z}
}

// Mul returns this vector multiplied by a scalar
func (v Vec3) Mul(scalar float32) Vec3 {
	return Vec3{X: v.X * scalar, Y: v.Y * scalar, Z: v.Z * scalar}
}

// Dot returns the dot product of this vector with another vector
func (v Vec3) Dot(other Vec3) float32 {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z
}
