// Package which holds all of the utility functions used throughout the codebase.
package utils

import (
	"asteroids/internal/constants"
	"math/rand/v2"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Main utility function for drawing lines which is used for drawing the entities
// of the game.
func DrawLines(
	origin rl.Vector2,
	scale float32,
	thickness float32,
	rotation float32,
	points []rl.Vector2,
) {
	// Lambda function to scale points
	scalePoints := func(origin rl.Vector2, scale float32, rot float32, p rl.Vector2) rl.Vector2 {
		return rl.Vector2Add(rl.Vector2Scale(rl.Vector2Rotate(p, rot), scale), origin)
	}

	for i := range points {
		rl.DrawLineEx(
			scalePoints(origin, scale, rotation, points[i]),
			scalePoints(origin, scale, rotation, points[(i+1)%len(points)]),
			thickness,
			rl.RayWhite,
		)
	}
}

// Returns a float32 in range of [minimum, maximum].
func RandInRange(minimum float32, maximum float32) float32 {
	return minimum + rand.Float32()*(maximum-minimum)
}

func DrawGameOverScreen() {
	rl.DrawText(
		"GAME OVER",
		constants.SCREEN_WIDTH/2-rl.MeasureText("GAME OVER", 40)/2,
		constants.SCREEN_HEIGHT/2-20,
		40,
		rl.Red,
	)
	rl.DrawText(
		"Press ENTER to Restart",
		constants.SCREEN_WIDTH/2-rl.MeasureText("Press ENTER to Restart", 20)/2,
		constants.SCREEN_HEIGHT/2+30,
		20,
		rl.RayWhite,
	)
}
