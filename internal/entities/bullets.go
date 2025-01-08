package entities

import (
	"asteroids/internal/constants"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Bullet struct {
	Start rl.Vector2 // Start coordinates of the bullet.
	End   rl.Vector2 // End coordinates of the bullet.
	Vel   rl.Vector2 // Velocity of the bullet.
	Dir   rl.Vector2 // The direction that the bullet is traveling in.
}

func NewBullet(pos rl.Vector2, rotation float32) Bullet {
	// Calculate the direction that the bullet is traveling towards. Rotation
	// is based on the direction that it was fired.
	direction := rl.Vector2{
		X: float32(-math.Sin(float64(rotation))),
		Y: float32(math.Cos(float64(rotation))),
	}

	return Bullet{
		Start: pos,
		End:   rl.Vector2Add(pos, rl.Vector2Scale(direction, constants.BULLET_LENGTH)),
		Vel:   rl.Vector2{X: 4, Y: 4},
		Dir:   direction,
	}
}

func DrawBullet(bullet Bullet) {
	rl.DrawLineV(bullet.Start, bullet.End, rl.RayWhite)
}
