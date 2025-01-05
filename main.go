package main

import (
	"asteroids/internal/constants"
	"asteroids/internal/ship"

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
)

type GameState struct {
	ship ship.Ship
}

func render(state *GameState) {
	ship.Draw(
		rl.Vector2{X: SCREEN_WIDTH / 2, Y: SCREEN_HEIGHT / 2},
		SCALE,
		THICKNESS,
		state.ship.Rot,
	)
}

func update(state *GameState) {
	// Side movements only handle the direction that the ship is facing.
	if rl.IsKeyDown(rl.KeyA) {
		state.ship.Rot -= ROTATION_SPEED
	}
	if rl.IsKeyDown(rl.KeyD) {
		state.ship.Rot += ROTATION_SPEED
	}
}

func main() {
	rl.InitWindow(SCREEN_WIDTH, SCREEN_HEIGHT, "Asteroids 1979")
	defer rl.CloseWindow()

	rl.SetTargetFPS(120)

	gameState := GameState{
		ship.New(),
	}

	for !rl.WindowShouldClose() {
		update(&gameState)

		rl.BeginDrawing()

		rl.ClearBackground(rl.Black)

		render(&gameState)

		rl.EndDrawing()
	}
}
