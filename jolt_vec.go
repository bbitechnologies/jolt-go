package jolt

// #include "wrapper/jolt_wrapper.h"
import "C"

// Vec3 represents a 3D vector
type Vec3 struct {
	X, Y, Z float32
}
