package ship

import (
	"asteroids/internal/constants"
	"asteroids/internal/utils"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Ship struct {
	Pos rl.Vector2 // Initial position of the ship
	Vel rl.Vector2 // Initial velocity of the ship
	Rot float32    // Rotation angle of the ship
}

// Initialises a new Ship struct. Position defaults to the middle of the window,
// velocity defaults to 2 and rotation defaults to 0.0 which is facing upwards.
func New() Ship {
	return Ship{
		Pos: rl.Vector2{X: constants.SCREEN_WIDTH / 2, Y: constants.SCREEN_HEIGHT / 2},
		Vel: rl.Vector2{X: 2, Y: 2},
		Rot: 0.0,
	}
}

func Draw(pos rl.Vector2, scale float32, thickness float32, rotation float32) {
	shipLines := []rl.Vector2{
		{X: -0.4, Y: -0.5},
		{X: 0.0, Y: 0.5},
		{X: 0.4, Y: -0.5},
		{X: 0.2, Y: -0.4},
		{X: -0.2, Y: -0.4},
	}

	utils.DrawLines(pos, scale, thickness, rotation, shipLines)
}
