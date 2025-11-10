package main

import (
	"fmt"
	"math"
	"time"

	"github.com/bbitechnologies/jolt-go"
)

const (
	MoveSpeed = float32(3.0)   // meters per second
	JumpSpeed = float32(7.0)   // meters per second
	GravityY  = float32(-9.81) // meters per second squared
)

// InputState represents player input at a given moment
type InputState struct {
	Forward  bool // W key
	Backward bool // S key
	Left     bool // A key
	Right    bool // D key
	Jump     bool // Space key
}

func (is InputState) String() string {
	s := ""
	if is.Forward {
		s += "W"
	}
	if is.Backward {
		s += "S"
	}
	if is.Left {
		s += "A"
	}
	if is.Right {
		s += "D"
	}
	if is.Jump {
		s += " [JUMP]"
	}
	if s == "" {
		s = "none"
	}
	return s
}

// InputRequest combines input state with timing information
type InputRequest struct {
	Input     InputState
	DeltaTime float32
}

// PlayerController manages character movement based on input
type PlayerController struct {
	character *jolt.CharacterVirtual
	velocity  jolt.Vec3
}

// NewPlayerController creates a new player controller
func NewPlayerController(character *jolt.CharacterVirtual) *PlayerController {
	return &PlayerController{
		character: character,
		velocity:  jolt.Vec3{X: 0, Y: 0, Z: 0},
	}
}

// ApplyInput processes input and updates character physics
func (pc *PlayerController) ApplyInput(req InputRequest) {
	// Calculate horizontal movement direction from input
	moveDir := jolt.Vec3{X: 0, Y: 0, Z: 0}

	if req.Input.Forward {
		moveDir.Z += 1
	}
	if req.Input.Backward {
		moveDir.Z -= 1
	}
	if req.Input.Right {
		moveDir.X += 1
	}
	if req.Input.Left {
		moveDir.X -= 1
	}

	// Normalize horizontal movement (prevent faster diagonal movement)
	magnitude := float32(math.Sqrt(float64(moveDir.X*moveDir.X + moveDir.Z*moveDir.Z)))
	if magnitude > 0 {
		moveDir.X /= magnitude
		moveDir.Z /= magnitude
	}

	// Get current velocity
	if pc.character.IsSupported() && pc.velocity.Y <= 0 {
		// On ground, read only horizontal velocity
		pc.velocity = pc.character.GetGroundVelocity()
	} else {
		// In air, read full velocity
		pc.velocity = pc.character.GetLinearVelocity()
	}

	// Apply move speed to horizontal velocity
	pc.velocity.X = moveDir.X * MoveSpeed
	pc.velocity.Z = moveDir.Z * MoveSpeed

	// Handle jumping (only when grounded)
	if req.Input.Jump && pc.character.IsSupported() && pc.velocity.Y <= 0 {
		pc.velocity.Y = JumpSpeed
	}

	// Apply gravity
	pc.velocity.Y += GravityY * req.DeltaTime

	// Set velocity and update character (ExtendedUpdate handles collision resolution)
	pc.character.SetLinearVelocity(pc.velocity)
	gravity := jolt.Vec3{X: 0, Y: GravityY, Z: 0}
	pc.character.ExtendedUpdate(req.DeltaTime, gravity)
}

// GetPosition returns the current character position
func (pc *PlayerController) GetPosition() jolt.Vec3 {
	return pc.character.GetPosition()
}

// GetGroundState returns the current ground state
func (pc *PlayerController) GetGroundState() jolt.GroundState {
	return pc.character.GetGroundState()
}

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
	character := ps.CreateCharacterVirtual(jolt.Vec3{X: 0, Y: 5, Z: 0})
	defer character.Destroy()

	// Create player controller
	controller := NewPlayerController(character)

	// Simulate player input for demo purposes
	fmt.Println("Player Controller Demo")
	fmt.Println("==========================================================")
	fmt.Println("Phase 1: Free fall and land")
	fmt.Println("Phase 2: Walk forward (simulated W key)")
	fmt.Println("Phase 3: Strafe right (simulated D key)")
	fmt.Println("Phase 4: Jump (simulated Space key)")
	fmt.Println("Phase 5: Diagonal movement (simulated S+D keys)")
	fmt.Println("==========================================================")

	elapsedTime := float32(0)

	for i := 0; i < 600; i++ {
		// Simulate different input patterns based on time
		input := InputState{}

		if elapsedTime >= 1.0 && elapsedTime < 3.0 {
			// Walk forward
			input.Forward = true
		} else if elapsedTime >= 3.0 && elapsedTime < 5.0 {
			// Strafe right
			input.Right = true
		} else if elapsedTime >= 5.0 && elapsedTime < 5.1 {
			// Jump
			input.Jump = true
		} else if elapsedTime >= 6.0 && elapsedTime < 9.0 {
			// Diagonal movement (backward + left)
			input.Backward = true
			input.Left = true
		}

		// Create input request
		req := InputRequest{
			Input:     input,
			DeltaTime: deltaTime,
		}

		// Apply input and update character
		controller.ApplyInput(req)

		// Print position and ground state every 0.5 seconds
		if i%30 == 0 {
			pos := controller.GetPosition()
			vel := character.GetLinearVelocity()
			groundState := controller.GetGroundState()
			fmt.Printf("[%.1fs] Position: X=% 6.2f Y=% 6.2f Z=% 6.2f | Velocity: X=% 5.2f Y=% 5.2f Z=% 5.2f | State: %s | Input: %s\n",
				elapsedTime, pos.X, pos.Y, pos.Z, vel.X, vel.Y, vel.Z, groundState, input.String())
		}

		elapsedTime += deltaTime
		time.Sleep(time.Millisecond * 16)
	}

	fmt.Println("==========================================================")
	fmt.Println("Demo complete!")
}
