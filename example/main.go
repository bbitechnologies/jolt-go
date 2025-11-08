package main

import (
	"fmt"
	"time"

	"github.com/bbitechnologies/jolt-go"
)

func main() {
	// Initialize Jolt Physics
	if err := jolt.Init(); err != nil {
		panic(err)
	}
	defer jolt.Shutdown()

	// Create physics world
	ps := jolt.NewPhysicsSystem()
	defer ps.Destroy()

	// Create a dynamic sphere at Y=20
	bi := ps.GetBodyInterface()
	sphere := bi.CreateSphere(1.0, jolt.Vec3{X: 0, Y: 20, Z: 0}, true)

	// Simulate for 3 seconds at 60 FPS
	fmt.Println("Simulating falling sphere...")
	deltaTime := float32(1.0 / 60.0)

	for i := 0; i < 180; i++ {
		ps.Update(deltaTime)

		// Print position every 0.5 seconds
		if i%30 == 0 {
			pos := bi.GetPosition(sphere)
			fmt.Printf("Frame %3d: Y = %.2f\n", i, pos.Y)
		}

		time.Sleep(time.Millisecond * 16)
	}

	fmt.Println("Done!")
}
