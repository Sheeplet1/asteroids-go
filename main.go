package main

import (
	"asteroids/internal/constants"
	"asteroids/internal/ship"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	SCREEN_HEIGHT = constants.SCREEN_HEIGHT
	SCREEN_WIDTH  = constants.SCREEN_WIDTH
)

type GameState struct {
	ship ship.Ship
}

func render(state *GameState) {
	ship.Draw(rl.Vector2{X: 0, Y: 0}, 38.0, 2.0, 0.0)
}

func update(state *GameState) {
}

func main() {
	rl.InitWindow(SCREEN_WIDTH, SCREEN_HEIGHT, "Asteroids 1979")
	defer rl.CloseWindow()

	rl.SetTargetFPS(120)

	gameState := GameState{
		ship.New(),
	}

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.Black)

		render(&gameState)

		rl.EndDrawing()
	}
}
