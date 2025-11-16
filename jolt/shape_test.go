package jolt

import (
	"math"
	"testing"
)

func TestShapeCastRay(t *testing.T) {
	// Create a sphere shape at origin with radius 1.0
	sphere := CreateSphere(1.0)
	defer sphere.Destroy()

	settings := DefaultRayCastSettings()

	t.Run("Ray hits sphere from above", func(t *testing.T) {
		ray := RRayCast{
			Origin:    Vec3{X: 0, Y: 10, Z: 0},
			Direction: Vec3{X: 0, Y: -20, Z: 0}, // Ray goes down 20 units
		}
		result := RayCastResult{}

		if !sphere.CastRay(ray, settings, &result) {
			t.Fatal("Ray should have hit the sphere")
		}

		hitDistance := result.Fraction * 20.0
		expectedDistance := float32(9.0) // 10 (start) - 1 (radius) = 9
		if math.Abs(float64(hitDistance-expectedDistance)) > 0.1 {
			t.Errorf("Hit distance = %.2f, expected ~%.2f", hitDistance, expectedDistance)
		}
	})

	t.Run("Ray misses sphere", func(t *testing.T) {
		ray := RRayCast{
			Origin:    Vec3{X: 5, Y: 0, Z: 0},
			Direction: Vec3{X: 0, Y: 10, Z: 0}, // Ray goes up, sphere is at origin
		}
		result := RayCastResult{}

		if sphere.CastRay(ray, settings, &result) {
			t.Error("Ray should have missed the sphere")
		}
	})

	t.Run("Ray hits sphere from side", func(t *testing.T) {
		ray := RRayCast{
			Origin:    Vec3{X: -5, Y: 0, Z: 0},
			Direction: Vec3{X: 10, Y: 0, Z: 0}, // Ray goes right 10 units
		}
		result := RayCastResult{}

		if !sphere.CastRay(ray, settings, &result) {
			t.Fatal("Ray should have hit the sphere")
		}

		hitDistance := result.Fraction * 10.0
		expectedDistance := float32(4.0) // 5 (start) - 1 (radius) = 4
		if math.Abs(float64(hitDistance-expectedDistance)) > 0.1 {
			t.Errorf("Hit distance = %.2f, expected ~%.2f", hitDistance, expectedDistance)
		}
	})

	t.Run("Ray starts inside sphere", func(t *testing.T) {
		ray := RRayCast{
			Origin:    Vec3{X: 0, Y: 0, Z: 0}, // Inside the sphere
			Direction: Vec3{X: 10, Y: 0, Z: 0},
		}
		result := RayCastResult{}

		// With TreatConvexAsSolid = true, ray starting inside should hit at fraction 0
		if !sphere.CastRay(ray, settings, &result) {
			t.Fatal("Ray starting inside should hit the sphere")
		}

		if result.Fraction > 0.01 {
			t.Errorf("Ray starting inside should hit at fraction ~0, got %.4f", result.Fraction)
		}
	})
}

func TestShapeCastRayBox(t *testing.T) {
	// Create a box shape centered at origin with half-extents (1, 2, 3)
	box := CreateBox(Vec3{X: 1, Y: 2, Z: 3})
	defer box.Destroy()

	settings := DefaultRayCastSettings()

	t.Run("Ray hits box from above", func(t *testing.T) {
		ray := RRayCast{
			Origin:    Vec3{X: 0, Y: 10, Z: 0},
			Direction: Vec3{X: 0, Y: -20, Z: 0},
		}
		result := RayCastResult{}

		if !box.CastRay(ray, settings, &result) {
			t.Fatal("Ray should have hit the box")
		}

		hitDistance := result.Fraction * 20.0
		expectedDistance := float32(8.0) // 10 (start) - 2 (half-extent Y) = 8
		if math.Abs(float64(hitDistance-expectedDistance)) > 0.1 {
			t.Errorf("Hit distance = %.2f, expected ~%.2f", hitDistance, expectedDistance)
		}
	})
}

func TestShapeCastRayCapsule(t *testing.T) {
	// Create a capsule shape centered at origin
	capsule := CreateCapsule(2.0, 0.5) // halfHeight=2.0, radius=0.5
	defer capsule.Destroy()

	settings := DefaultRayCastSettings()

	t.Run("Ray hits capsule", func(t *testing.T) {
		ray := RRayCast{
			Origin:    Vec3{X: -5, Y: 0, Z: 0},
			Direction: Vec3{X: 10, Y: 0, Z: 0},
		}
		result := RayCastResult{}

		if !capsule.CastRay(ray, settings, &result) {
			t.Fatal("Ray should have hit the capsule")
		}

		hitDistance := result.Fraction * 10.0
		// Capsule extends 0.5 units in X direction (radius), so hit at 5 - 0.5 = 4.5
		if hitDistance < 4.0 || hitDistance > 5.0 {
			t.Errorf("Hit distance = %.2f, expected between 4.0 and 5.0", hitDistance)
		}
	})
}
