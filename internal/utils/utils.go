// Package which holds all of the utility functions used throughout the codebase.
package utils

import rl "github.com/gen2brain/raylib-go/raylib"

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
