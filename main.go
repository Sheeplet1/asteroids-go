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
	bullets       []entities.Bullet
	bulletTimer   float32
	lives         uint8
	isGameOver    bool
	Score         uint64
}

func NewGameState() GameState {
	return GameState{
		ship:          entities.NewShip(),
		asteroids:     []entities.Asteroid{},
		asteroidTimer: 0,
		bullets:       []entities.Bullet{},
		bulletTimer:   0,
		lives:         3,
		isGameOver:    false,
		Score:         0,
	}
}

func render(state *GameState) {
	// Renders the lives counter in the top right of the window.
	livesStr := fmt.Sprintf("Lives: %o", state.lives)
	rl.DrawText(livesStr, SCREEN_WIDTH-16-rl.MeasureText(livesStr, 30), 20, 30, rl.RayWhite)

	// Renders the death timer of the ship in the middle of the screen.
	if state.ship.IsDead() {
		deathStr := fmt.Sprintf("Respawning in %.0f", state.ship.DeathTimer)
		rl.DrawText(
			deathStr,
			SCREEN_WIDTH/2-rl.MeasureText(deathStr, 30)/2,
			SCREEN_HEIGHT/2,
			30,
			rl.RayWhite,
		)
	}

	// Renders the score in the top left of the window.
	scoreStr := fmt.Sprintf("Score: %d", state.Score)
	rl.DrawText(scoreStr, 16, 20, 30, rl.RayWhite)

	// If the ship is moving forward, then we draw thrusters onto the ship
	// for the effect.
	entities.RenderShip(&state.ship)

	// ------------------------------------------------------------------------
	// Bullet rendering
	// ------------------------------------------------------------------------
	if rl.IsKeyDown(rl.KeySpace) && state.bulletTimer <= 0 {
		bullet := entities.NewBullet(state.ship.Pos, state.ship.Rot)
		fmt.Println(bullet)
		entities.DrawBullet(bullet)
		state.bullets = append(state.bullets, bullet)
		state.bulletTimer = 1
	}

	for _, bullet := range state.bullets {
		entities.DrawBullet(bullet)
	}
	// ------------------------------------------------------------------------

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
		entities.DrawAsteroid(asteroid)
	}
	// ------------------------------------------------------------------------

	// Increment/decrement timers.
	state.asteroidTimer += rl.GetFrameTime()
	state.bulletTimer -= rl.GetFrameTime()

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
				state.asteroids[i].Pos.X < -entities.SPAWN_MARGIN ||
				state.asteroids[i].Pos.Y > SCREEN_HEIGHT+entities.SPAWN_MARGIN ||
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
			float32(asteroid.Hitbox),
		) && !state.ship.IsDead() {
			state.ship.DeathTimer += 5
			state.lives -= 1
		}
	}
}

// Update all bullets positions.
func updateBulletPositions(state *GameState) {
	if len(state.bullets) > 0 {
		// Update bullet start and ending position.
		for i := range state.bullets {
			state.bullets[i].Start = rl.Vector2Add(
				state.bullets[i].Start,
				rl.Vector2Multiply(state.bullets[i].Vel, state.bullets[i].Dir),
			)

			state.bullets[i].End = rl.Vector2Add(
				state.bullets[i].Start,
				rl.Vector2Scale(state.bullets[i].Dir, constants.BULLET_LENGTH),
			)
		}

		// Remove any bullets that go out of bounds.
		for i := len(state.bullets) - 1; i >= 0; i-- {
			if state.bullets[i].Start.X > SCREEN_WIDTH+entities.SPAWN_MARGIN ||
				state.bullets[i].Start.X < -entities.SPAWN_MARGIN ||
				state.bullets[i].Start.Y > SCREEN_HEIGHT+entities.SPAWN_MARGIN ||
				state.bullets[i].Start.Y < -entities.SPAWN_MARGIN {
				state.bullets = append(state.bullets[:i], state.bullets[i+1:]...)
			}
		}
	}
}

func checkForBulletAsteroidCollisions(state *GameState) {
	if len(state.bullets) == 0 || len(state.asteroids) == 0 {
		return
	}
	// TODO: Need to efficiently check for bullet-asteroid collisions since
	// this could be a performance bottleneck as it runs every frame.
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

	// Update foreign entities positions.
	updateAsteroidPositions(state)
	updateBulletPositions(state)

	// Check for any entity collisions.
	checkForShipAsteroidCollisions(state)
	checkForBulletAsteroidCollisions(state)

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
