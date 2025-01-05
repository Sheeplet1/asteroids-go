package entities

import (
	"asteroids/internal/constants"
	"asteroids/internal/utils"
	"math/rand/v2"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	SPAWN_MARGIN    = 100
	ASTEROID_HITBOX = 25
)

type Asteroid struct {
	Pos rl.Vector2
	Vel rl.Vector2
	Dir rl.Vector2 // This will point towards the ship's location when it first spawned.
	// TODO: Store the coordinates that make the polygon of the asteroid
	// in order to redraw it.

	// TODO: Add health
}

func NewAsteroid(pos rl.Vector2, dir rl.Vector2) Asteroid {
	return Asteroid{
		Pos: pos,
		Vel: rl.Vector2{X: 2, Y: 2},
		Dir: dir,
	}
}

// Generates spawn point coordinates for the asteroids. Spawn point coordinates
// are contained to coordinates that are outside of the window dimensions. This
// is so that the asteroids can spawn outside and float into view.
func GenerateAsteroidSpawn() rl.Vector2 {
	// Generate a random zone for the asteroid to spawn in.
	zone := rand.IntN(4)

	const (
		// Defining zones for the asteroids to spawn in.
		TOP   = 0
		BOT   = 1
		LEFT  = 2
		RIGHT = 3
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
	rl.DrawCircleLinesV(pos, 25, rl.RayWhite)
}
