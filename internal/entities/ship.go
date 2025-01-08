package entities

import (
	"asteroids/internal/constants"
	"asteroids/internal/utils"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	SHIP_HITBOX_RADIUS = 15

	// Ship movement constants
	ROTATION_SPEED = 0.08
	ACCEL          = 0.15
	DECEL          = 0.01
	MIN_VEL        = 2.0
	MAX_VEL        = 6.0
	DRAG           = 0.01
)

type Ship struct {
	Pos        rl.Vector2 // Initial position of the ship
	Vel        rl.Vector2 // Initial velocity of the ship
	Rot        float32    // Rotation angle of the ship
	DeathTimer float32    // Death timer for the ship
}

// Returns true/false whether the ship is dead or not.
func (ship Ship) IsDead() bool {
	return ship.DeathTimer > 0
}

// Initialises a new Ship struct. Position defaults to the middle of the window,
// velocity defaults to 2 and rotation defaults to 0.0 which is facing upwards.
func NewShip() Ship {
	return Ship{
		Pos:        rl.Vector2{X: constants.SCREEN_WIDTH / 2, Y: constants.SCREEN_HEIGHT / 2},
		Vel:        rl.Vector2{X: 2, Y: 2},
		Rot:        0,
		DeathTimer: 0,
	}
}

// Draws the base ship without thrusters.
func drawShip(pos rl.Vector2, scale float32, thickness float32, rotation float32) {
	shipLines := []rl.Vector2{
		{X: -0.4, Y: -0.5},
		{X: 0.0, Y: 0.5},
		{X: 0.4, Y: -0.5},
		{X: 0.2, Y: -0.4},
		{X: -0.2, Y: -0.4},
	}

	utils.DrawLines(pos, scale, thickness, rotation, shipLines)
}

func drawShipWithThrusters(pos rl.Vector2, scale float32, thickness float32, rotation float32) {
	shipWithThrusters := []rl.Vector2{
		{X: -0.4, Y: -0.5},
		{X: 0.0, Y: 0.5},
		{X: 0.4, Y: -0.5},
		{X: 0.2, Y: -0.4},

		// Drawing the thrusters
		{X: 0.0, Y: -0.8},
		{X: -0.2, Y: -0.4},
		{X: 0.2, Y: -0.4},
		{X: -0.2, Y: -0.4},
		{X: -0.4, Y: -0.5},
	}

	utils.DrawLines(pos, scale, thickness, rotation, shipWithThrusters)
}

// Updates the ship depending on whether its dead and movement keys.
func UpdateShip(ship *Ship) {
	if ship.IsDead() {
		return
	}

	// Side movements only handle the direction that the ship is facing.
	if rl.IsKeyDown(rl.KeyA) {
		ship.Rot -= ROTATION_SPEED
	}
	if rl.IsKeyDown(rl.KeyD) {
		ship.Rot += ROTATION_SPEED
	}

	// Handle forward and backward movements for the ship.
	if rl.IsKeyDown(rl.KeyW) {
		ship.Vel = rl.Vector2ClampValue(
			rl.Vector2Scale(ship.Vel, 1.0+ACCEL),
			MIN_VEL,
			MAX_VEL,
		)
	}
	if rl.IsKeyDown(rl.KeyS) {
		ship.Vel = rl.Vector2ClampValue(
			rl.Vector2Scale(ship.Vel, 1.0-DECEL),
			0,
			MAX_VEL,
		)
	}

	// Calculate the ship's velocity after accounting for drag. Creates that
	// floating through space feel.
	ship.Vel = rl.Vector2Scale(ship.Vel, 1.0-DRAG)

	// Updating the ship's position after accounting for all velocity changes.
	shipDirection := rl.Vector2{
		X: float32(-math.Sin(float64(ship.Rot))),
		Y: float32(math.Cos(float64(ship.Rot))),
	}
	ship.Pos = rl.Vector2Add(
		ship.Pos,
		rl.Vector2Multiply(ship.Vel, shipDirection),
	)

	// Handle out of bounds movements of the ship. The ship going out of bounds
	// will simply teleport the ship to the opposite side of where it was going.
	ship.Pos.X = float32(math.Mod(float64(ship.Pos.X), float64(constants.SCREEN_WIDTH)))
	ship.Pos.Y = float32(math.Mod(float64(ship.Pos.Y), float64(constants.SCREEN_HEIGHT)))
}

// Renders the ship based on whether its dead. The ship will render with thrusters
// if the ship is moving forward.
func RenderShip(ship *Ship) {
	if !ship.IsDead() {
		drawShip(
			ship.Pos,
			constants.SCALE,
			constants.THICKNESS,
			ship.Rot,
		)

		if rl.IsKeyDown(rl.KeyW) {
			drawShipWithThrusters(ship.Pos, constants.SCALE, constants.THICKNESS, ship.Rot)
		}
	}
}
