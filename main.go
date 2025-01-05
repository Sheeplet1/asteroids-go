package main

import (
	"asteroids/internal/constants"
	"asteroids/internal/entities"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	SCREEN_HEIGHT = constants.SCREEN_HEIGHT
	SCREEN_WIDTH  = constants.SCREEN_WIDTH

	// Default drawing parameters
	THICKNESS = 2.0
	SCALE     = 38.0

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
	if !state.ship.IsDead() {
		entities.DrawShip(
			state.ship.Pos,
			SCALE,
			THICKNESS,
			state.ship.Rot,
		)

		if rl.IsKeyDown(rl.KeyW) {
			entities.DrawShipWithThrusters(state.ship.Pos, SCALE, THICKNESS, state.ship.Rot)
		}
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
	entities.UpdateShip(&state.ship)

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

	// TODO: Check for collisions with the ship and any asteroid.
	for _, asteroid := range state.asteroids {
		// If the asteroid hits the ship, then the ship dies and we introduce
		// a 5 second death timer.
		if rl.CheckCollisionCircles(
			state.ship.Pos,
			entities.SHIP_HITBOX_RADIUS,
			asteroid.Pos,
			entities.ASTEROID_HITBOX,
		) && !state.ship.IsDead() {
			state.ship.DeathTimer += 5
			state.ship.Lives -= 1
		}
	}

	if state.ship.DeathTimer > 0 {
		state.ship.DeathTimer -= rl.GetFrameTime()
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
