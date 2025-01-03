package main

import (
	"asteroids/internal/constants"
	"asteroids/internal/ship"
	"fmt"

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
}

func update(state *GameState) {
}

func main() {
	fmt.Println("Hello, World!")

	rl.InitWindow(SCREEN_WIDTH, SCREEN_HEIGHT, "Asteroids 1979")
	defer rl.CloseWindow()

	rl.SetTargetFPS(120)

	gameState := GameState{
		ship.New(),
	}

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.Black)
		rl.DrawText(
			"Congrats! You created your first window!",
			SCREEN_WIDTH/2,
			SCREEN_HEIGHT/2,
			18,
			rl.RayWhite,
		)

		rl.EndDrawing()
	}
}
