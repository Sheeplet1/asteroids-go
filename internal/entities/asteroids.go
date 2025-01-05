package entities

import (
	"asteroids/internal/constants"
	"asteroids/internal/utils"
	"math/rand/v2"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Generates spawn point coordinates for the asteroids. Spawn point coordinates
// are contained to coordinates that are outside of the window dimensions. This
// is so that the asteroids can spawn outside and float into view.
func GenerateAsteroidSpawn() rl.Vector2 {
	zone := rand.IntN(4)

	const (
		// Defining zones for the asteroids to spawn in.
		TOP          = 0
		BOT          = 1
		LEFT         = 2
		RIGHT        = 3
		SPAWN_MARGIN = 100
	)

	switch zone {
	case TOP:
		x := utils.RandInRange(0, constants.SCREEN_WIDTH)
		y := utils.RandInRange(-SPAWN_MARGIN, 0)
		return rl.Vector2{X: x, Y: y}
	case BOT:
		x := utils.RandInRange(0, constants.SCREEN_WIDTH)
		y := utils.RandInRange(constants.SCREEN_HEIGHT, constants.SCREEN_HEIGHT+SPAWN_MARGIN)
		return rl.Vector2{X: x, Y: y}
	case LEFT:
		x := utils.RandInRange(-SPAWN_MARGIN, 0)
		y := utils.RandInRange(0, constants.SCREEN_HEIGHT)
		return rl.Vector2{X: x, Y: y}
	case RIGHT:
		x := utils.RandInRange(constants.SCREEN_WIDTH, constants.SCREEN_WIDTH+SPAWN_MARGIN)
		y := utils.RandInRange(0, constants.SCREEN_HEIGHT)
		return rl.Vector2{X: x, Y: y}
	default:
		panic("unreachable: unexpected zone value")
	}
}

// Draws the asteroid at any given `pos`.
func DrawAsteroid(pos rl.Vector2) {
	rl.DrawCircleV(pos, 9, rl.RayWhite)
}
