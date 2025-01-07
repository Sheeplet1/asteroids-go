package main

import (
	"asteroids/internal/constants"
	"asteroids/internal/entities"
	"asteroids/internal/utils"
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	SCREEN_HEIGHT = constants.SCREEN_HEIGHT
	SCREEN_WIDTH  = constants.SCREEN_WIDTH

	// Default drawing parameters
	THICKNESS = constants.THICKNESS
	SCALE     = constants.SCALE

	// Default game parameters
	ASTEROID_SPAWN_INTERVAL = constants.ASTEROID_SPAWN_INTERVAL
)

type GameState struct {
	ship          entities.Ship
	asteroids     []entities.Asteroid // Slice of asteroids present in the game.
	asteroidTimer float32             // The spawn timer for the asteroids.
	lives         uint8
	isGameOver    bool
}

func NewGameState() GameState {
	return GameState{
		ship:          entities.NewShip(),
		asteroids:     []entities.Asteroid{},
		asteroidTimer: 0,
		lives:         3,
		isGameOver:    false,
	}
}

func render(state *GameState) {
	livesStr := fmt.Sprintf("Lives: %o", state.lives)
	rl.DrawText(livesStr, SCREEN_WIDTH-16-rl.MeasureText(livesStr, 30), 20, 30, rl.RayWhite)

	// If the ship is moving forward, then we draw thrusters onto the ship
	// for the effect.
	entities.RenderShip(&state.ship)

	// ------------------------------------------------------------------------
	// Asteroid rendering
	// ------------------------------------------------------------------------
	if state.asteroidTimer >= ASTEROID_SPAWN_INTERVAL {
		// Creating a new asteroid to spawn in.
		asteroid := entities.SpawnAsteroid(state.ship.Pos)

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
	// ------------------------------------------------------------------------

	// If the game is over, render the game over screen.
	if state.isGameOver {
		utils.DrawGameOverScreen()
	}
}

// Iterates through the existing asteroids in the game and updates their positions
// based on their velocity and direction.
func updateAsteroidPositions(state *GameState) {
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

// The ship and asteroid have circular hitboxes which we check for any collisions
// between them.
func checkForShipAsteroidCollisions(state *GameState) {
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
			state.lives -= 1
		}
	}
}

func update(state *GameState) {
	if state.isGameOver {
		if rl.IsKeyPressed(rl.KeyEnter) {
			*state = NewGameState()
		}
		return
	}

	// Updates the ship based on movement or death.
	entities.UpdateShip(&state.ship)

	// Update any asteroid's positions.
	updateAsteroidPositions(state)

	// Check for any asteroid and ship collisions. If there is any collision,
	// the ship dies and overall life count reduces by 1.
	checkForShipAsteroidCollisions(state)

	// Updating the death timer of the ship.
	if state.ship.DeathTimer > 0 {
		state.ship.DeathTimer -= rl.GetFrameTime()
	}

	// If there is no more lives left, set the game state to be over.
	if state.lives <= 0 {
		state.isGameOver = true
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
