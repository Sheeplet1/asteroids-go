package entities

import (
	"asteroids/internal/constants"
	"asteroids/internal/utils"
	"math"
	"math/rand/v2"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	SPAWN_MARGIN    = 100
	ASTEROID_HITBOX = 25
)

type Asteroid struct {
	Pos    rl.Vector2
	Vel    rl.Vector2
	Dir    rl.Vector2 // This will point towards the ship's location when it first spawned.
	Points []rl.Vector2

	// TODO: Add health
}

// Generates a random n-sided polygon shape generating randomized points around
// a circle. Each point is calculated based on polar coordinates before being
// converted into cartesian coordinates and returned.
func generateAsteroidShape(
	numSides int,
	minRadius float64,
	maxRadius float64,
) []rl.Vector2 {
	points := []rl.Vector2{}

	// Divide the circle into equal segments based on the number of sides.
	angleStep := math.Pi * 2 / float64(numSides)

	for i := 0; i < numSides; i++ {
		// Calculate the expected angle for the ith point.
		targetAngle := angleStep * float64(i)

		// Add some random angle variation for the irregularity.
		angle := targetAngle + (rand.Float64()-0.5)*angleStep*0.25

		// Pick a random radius between the given minimum and maximum radius.
		radius := minRadius + rand.Float64()*(maxRadius-minRadius)

		// Convert into cartesian coordinates.
		x := math.Cos(angle) * radius
		y := math.Sin(angle) * radius
		points = append(points, rl.Vector2{X: float32(x), Y: float32(y)})
	}

	return points
}

func newAsteroid(pos rl.Vector2, dir rl.Vector2) Asteroid {
	// Generates points for the polygon shape of the asteroid.
	const DEFAULT_NUM_SIDES = 11
	points := generateAsteroidShape(DEFAULT_NUM_SIDES, 0.5, 1)

	return Asteroid{
		Pos:    pos,
		Vel:    rl.Vector2{X: 2, Y: 2},
		Dir:    dir,
		Points: points,
	}
}

// Generates spawn point coordinates for the asteroids. Spawn point coordinates
// are contained to coordinates that are outside of the window dimensions. This
// is so that the asteroids can spawn outside and float into view.
func generateAsteroidSpawn() rl.Vector2 {
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
func DrawAsteroid(asteroid Asteroid) {
	utils.DrawLines(asteroid.Pos, constants.SCALE, constants.THICKNESS, 0.0, asteroid.Points)
}

// Spawns an asteroid and returns the Asteroid struct to be appended into the
// game state. Takes in the ship's position as the asteroid drifts towards
// the ship when it spawns.
func SpawnAsteroid(shipPos rl.Vector2) Asteroid {
	spawnPoint := generateAsteroidSpawn()
	asteroid := newAsteroid(
		spawnPoint,
		rl.Vector2Normalize(rl.Vector2Subtract(shipPos, spawnPoint)),
	)
	return asteroid
}
