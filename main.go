package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	SCREEN_WIDTH  = 1280
	SCREEN_HEIGHT = 960
)

func main() {
	fmt.Println("Hello, World!")

	rl.InitWindow(SCREEN_WIDTH, SCREEN_HEIGHT, "Asteroids 1979")
	defer rl.CloseWindow()

	rl.SetTargetFPS(120)

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
