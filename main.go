package main

import (
	"asteroids/internal/constants"
	"asteroids/internal/entities"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	SCREEN_HEIGHT = constants.SCREEN_HEIGHT
	SCREEN_WIDTH  = constants.SCREEN_WIDTH

	// Default drawing parameters
	THICKNESS = 2.0
	SCALE     = 38.0

	// Ship movement constants
	ROTATION_SPEED = 0.1
	ACCEL          = 0.15
	DECEL          = 0.01
	MIN_VEL        = 2.0
	MAX_VEL        = 6.0
	DRAG           = 0.01

	// Default game parameters
	ASTEROID_SPAWN_INTERVAL = 2.5 // Asteroids spawn every 2.5 seconds
)

type GameState struct {
	ship          entities.Ship
	asteroids     []entities.Asteroid // Slice of asteroids present in the game.
	asteroidTimer float32             // The spawn timer for the asteroids.
}

func NewGameState() GameState {
	return GameState{
		ship:          entities.NewShip(),
		asteroids:     []entities.Asteroid{},
		asteroidTimer: 0,
	}
}

func render(state *GameState) {
	// If the ship is moving forward, then we draw thrusters onto the ship
	// for the effect.
	if rl.IsKeyDown(rl.KeyW) {
		entities.DrawShipWithThrusters(state.ship.Pos, SCALE, THICKNESS, state.ship.Rot)
	} else {
		entities.DrawShip(
			state.ship.Pos,
			SCALE,
			THICKNESS,
			state.ship.Rot,
		)
	}

	if state.asteroidTimer >= ASTEROID_SPAWN_INTERVAL {
		// Creating a new asteroid to spawn in.
		spawnPoint := entities.GenerateAsteroidSpawn()
		asteroid := entities.NewAsteroid(
			spawnPoint,
			rl.Vector2Normalize(rl.Vector2Subtract(state.ship.Pos, spawnPoint)),
		)

		// Add new asteroid into the game state.
		state.asteroids = append(state.asteroids, asteroid)
		state.asteroidTimer = 0
	}

	// Render any asteroids that are already in the game.
	for _, asteroid := range state.asteroids {
		entities.DrawAsteroid(asteroid.Pos)
	}

	// Increment the spawn timer for the asteroids.
	state.asteroidTimer += rl.GetFrameTime()
}

func update(state *GameState) {
	// Side movements only handle the direction that the ship is facing.
	if rl.IsKeyDown(rl.KeyA) {
		state.ship.Rot -= ROTATION_SPEED
	}
	if rl.IsKeyDown(rl.KeyD) {
		state.ship.Rot += ROTATION_SPEED
	}

	// Handle forward and backward movements for the ship.
	if rl.IsKeyDown(rl.KeyW) {
		state.ship.Vel = rl.Vector2ClampValue(
			rl.Vector2Scale(state.ship.Vel, 1.0+ACCEL),
			MIN_VEL,
			MAX_VEL,
		)
	}
	if rl.IsKeyDown(rl.KeyS) {
		state.ship.Vel = rl.Vector2ClampValue(
			rl.Vector2Scale(state.ship.Vel, 1.0-DECEL),
			0,
			MAX_VEL,
		)
	}

	// Calculate the ship's velocity after accounting for drag. Creates that
	// floating through space feel.
	state.ship.Vel = rl.Vector2Scale(state.ship.Vel, 1.0-DRAG)

	// Updating the ship's position after accounting for all velocity changes.
	shipDirection := rl.Vector2{
		X: float32(-math.Sin(float64(state.ship.Rot))),
		Y: float32(math.Cos(float64(state.ship.Rot))),
	}
	state.ship.Pos = rl.Vector2Add(
		state.ship.Pos,
		rl.Vector2Multiply(state.ship.Vel, shipDirection),
	)

	// Handle out of bounds movements of the ship. The ship going out of bounds
	// will simply teleport the ship to the opposite side of where it was going.
	state.ship.Pos.X = float32(math.Mod(float64(state.ship.Pos.X), float64(SCREEN_WIDTH)))
	state.ship.Pos.Y = float32(math.Mod(float64(state.ship.Pos.Y), float64(SCREEN_HEIGHT)))

	// Update any asteroid's positions.
	if len(state.asteroids) > 0 {
		for i := range state.asteroids {
			state.asteroids[i].Pos = rl.Vector2Add(
				state.asteroids[i].Pos,
				rl.Vector2Multiply(state.asteroids[i].Vel, state.asteroids[i].Dir),
			)
		}

		for i := len(state.asteroids) - 1; i >= 0; i-- {
			// PERF: Could have some performance optimisations here but not
			// needed due to the size of the game.

			// If the asteroid moves out of the window dimensions, we remove it
			// from the game.
			if state.asteroids[i].Pos.X > SCREEN_WIDTH+entities.SPAWN_MARGIN ||
				state.asteroids[i].Pos.X < -entities.SPAWN_MARGIN {
				state.asteroids = append(state.asteroids[:i], state.asteroids[i+1:]...)
			}
			if state.asteroids[i].Pos.Y > SCREEN_HEIGHT+entities.SPAWN_MARGIN ||
				state.asteroids[i].Pos.Y < -entities.SPAWN_MARGIN {
				state.asteroids = append(state.asteroids[:i], state.asteroids[i+1:]...)
			}
		}
	}
}

func main() {
	rl.InitWindow(SCREEN_WIDTH, SCREEN_HEIGHT, "Asteroids 1979")
	defer rl.CloseWindow()

	rl.SetTargetFPS(120)

	gameState := NewGameState()

	for !rl.WindowShouldClose() {
		update(&gameState)

		rl.BeginDrawing()

		rl.ClearBackground(rl.Black)

		render(&gameState)

		rl.EndDrawing()
	}
}
