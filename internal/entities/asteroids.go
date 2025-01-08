package entities

import (
	"asteroids/internal/constants"
	"asteroids/internal/utils"
	"math"
	"math/rand/v2"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	SPAWN_MARGIN = constants.SPAWN_MARGIN

	// Hitbox sizes for the asteroids.
	SMALL_HITBOX = 10
	MED_HITBOX   = 25
	LARGE_HITBOX = 40

	// Scores that the asteroids will give when destroyed.
	SMALL_SCORE = 100 // Smaller asteroids are harder to hit, so they give more points.
	MED_SCORE   = 50
	LARGE_SCORE = 20
)

type AsteroidSize int

const (
	Small AsteroidSize = iota
	Medium
	Large
)

type Asteroid struct {
	Pos    rl.Vector2   // Position of the asteroid.
	Vel    rl.Vector2   // Velocity of the asteroid.
	Dir    rl.Vector2   // This will point towards the ship's location when it first spawned.
	Points []rl.Vector2 // Points which generate the asteroid's shape.
	Size   AsteroidSize // Size of the asteroid.
	Hitbox int          // Radius of the hitbox of the asteroid. Hitbox is in the shape of a circle.
	Health int          // Health of the asteroid. Asteroids will break into smaller asteroids when health reaches 0.
	Score  uint64       // Score that the player will receive when the asteroid is destroyed.
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

func newAsteroid(pos rl.Vector2, dir rl.Vector2, size AsteroidSize) Asteroid {
	var minRadius, maxRadius float64
	var velocity rl.Vector2
	var hitbox, health int
	var score uint64

	// Defining variables based on the size of the asteroid. Larger asteroids
	// will have more health but less speed etc.
	switch size {
	case Small:
		minRadius, maxRadius = 0, 0.5
		velocity = rl.Vector2{X: 3, Y: 3}
		hitbox = SMALL_HITBOX
		health = 1
		score = SMALL_SCORE
	case Medium:
		minRadius, maxRadius = 0.5, 1
		velocity = rl.Vector2{X: 2, Y: 2}
		hitbox = MED_HITBOX
		health = 2
		score = MED_SCORE
	case Large:
		minRadius, maxRadius = 1, 1.5
		velocity = rl.Vector2{X: 1, Y: 1}
		hitbox = LARGE_HITBOX
		health = 3
		score = LARGE_SCORE
	}

	// Generates points for the polygon shape of the asteroid.
	const DEFAULT_NUM_SIDES = 11
	points := generateAsteroidShape(DEFAULT_NUM_SIDES, minRadius, maxRadius)

	return Asteroid{
		Pos:    pos,
		Vel:    velocity,
		Dir:    dir,
		Points: points,
		Hitbox: hitbox,
		Health: health,
		Score:  score,
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
	rl.DrawCircleLinesV(asteroid.Pos, float32(asteroid.Hitbox), rl.Yellow)
	utils.DrawLines(asteroid.Pos, constants.SCALE, constants.THICKNESS, 0.0, asteroid.Points)
}

// Spawns an asteroid and returns the Asteroid struct to be appended into the
// game state. Takes in the ship's position as the asteroid drifts towards
// the ship when it spawns.
func SpawnAsteroid(shipPos rl.Vector2) Asteroid {
	spawnPoint := generateAsteroidSpawn()

	// Randomly generate a size for the asteroid.
	// 40% chance for large, 40% chance for medium, 20% chance for small.
	size := Large
	if rand.Float64() < 0.4 {
		size = Medium
	} else if rand.Float64() < 0.2 {
		size = Small
	}

	asteroid := newAsteroid(
		spawnPoint,
		rl.Vector2Normalize(rl.Vector2Subtract(shipPos, spawnPoint)),
		size,
	)

	return asteroid
}
