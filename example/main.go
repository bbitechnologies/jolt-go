package main

import (
	"fmt"
	"math"
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

	// Create a large static platform (20x1x20 box at Y=0)
	bi := ps.GetBodyInterface()
	floor := bi.CreateBox(
		jolt.Vec3{X: 10, Y: 0.5, Z: 10}, // half-extents (creates 20x1x20 box)
		jolt.Vec3{X: 0, Y: 0, Z: 0},     // position
		false,                           // static
	)
	defer floor.Destroy()

	// Update physics once to initialize broad phase
	deltaTime := float32(1.0 / 60.0)
	ps.Update(deltaTime)

	// Create player character above the platform at Y=5 (will fall to ground)
	player := ps.CreateCharacterVirtual(jolt.Vec3{X: 0, Y: 5, Z: 0})
	defer player.Destroy()

	// Simulate for 5 seconds at ~60 FPS
	fmt.Println("Player controller demo - character walks in a circle")
	fmt.Println("==========================================================")

	gravity := jolt.Vec3{X: 0, Y: -9.81, Z: 0}
	elapsedTime := float32(0)
	verticalVelocity := float32(0) // Track vertical velocity for gravity

	for i := 0; i < 300; i++ {
		// Make the player walk in a circle (radius 5, period 5 seconds)
		angle := elapsedTime * 2.0 * math.Pi / 5.0
		speed := float32(3.0) // meters per second
		velocityX := float32(math.Cos(float64(angle))) * speed
		velocityZ := float32(math.Sin(float64(angle))) * speed

		// Apply gravity to vertical velocity (kinematic integration)
		verticalVelocity += gravity.Y * deltaTime

		// Set combined velocity (horizontal movement + gravity-affected vertical)
		player.SetLinearVelocity(jolt.Vec3{
			X: velocityX,
			Y: verticalVelocity,
			Z: velocityZ,
		})

		// Update physics
		ps.Update(deltaTime)
		player.ExtendedUpdate(deltaTime, gravity)

		// Reset vertical velocity if character is supported (on ground)
		if player.IsSupported() && verticalVelocity < 0 {
			verticalVelocity = 0
		}

		// Print position and ground state every 0.5 seconds
		if i%30 == 0 {
			pos := player.GetPosition()
			groundState := player.GetGroundState()
			fmt.Printf("[%.1fs] Position: X=% 6.2f Y=% 6.2f Z=% 6.2f | State: %s\n",
				elapsedTime, pos.X, pos.Y, pos.Z, groundState)
		}

		elapsedTime += deltaTime
		time.Sleep(time.Millisecond * 16)
	}

	fmt.Println("==========================================================")
	fmt.Println("Demo complete!")
}
