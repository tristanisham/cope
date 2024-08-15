package game

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tristanisham/cope/meta"
)

func (g *Game) handleMovement() {
	movement := struct {
		dx float64
		dy float64
	}{0, 0}

	// Calculate intended movement
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		movement.dx = meta.MovementSpeed
	}

	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		movement.dy = meta.MovementSpeed
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		movement.dx = -meta.MovementSpeed
	}

	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		movement.dy = -meta.MovementSpeed
	}

	// Normalize movement to prevent faster diagonal movement
	if movement.dx != 0 && movement.dy != 0 {
		movement.dx /= math.Sqrt2
		movement.dy /= math.Sqrt2
	}

	// Check collision and update position if no collision
	futureX := g.px + movement.dx
	futureY := g.py + movement.dy

	if !isColliding(futureX, g.py, g.objects) {
		g.px = futureX
	}
	if !isColliding(g.px, futureY, g.objects) {
		g.py = futureY
	}

	// Ensure player stays within screen bounds
	if g.px >= meta.ScreenWidth-meta.Padding {
		g.px = meta.ScreenWidth - meta.Padding - 1
	}
	if g.px <= meta.Padding {
		g.px = meta.Padding + 1
	}
	if g.py >= meta.ScreenHeight-meta.Padding {
		g.py = meta.ScreenHeight - meta.Padding - 1
	}
	if g.py <= meta.Padding {
		g.py = meta.Padding + 1
	}
}

func isColliding(x, y float64, objects []object) bool {
	// // Define a small padding for collision detection
	// const playerSize = 2.5 // Adjust based on your player size

	// // Define the player's bounding box
	// playerBox := line{x - playerSize, y - playerSize, x + playerSize, y + playerSize}

	// for _, obj := range objects {
	// 	for _, wall := range obj.walls {
	// 		// Check if the player's bounding box intersects with any wall
	// 		if _, _, collides := intersection(playerBox, wall); collides {
	// 			return true
	// 		}
	// 	}
	// }

	const playerHalfSize = 2.1 // Half the size of the player, adjust based on your player size

	// Define the player's bounding box using the four corners
	playerBox := []line{
		{x - playerHalfSize, y - playerHalfSize, x + playerHalfSize, y - playerHalfSize}, // Top side
		{x + playerHalfSize, y - playerHalfSize, x + playerHalfSize, y + playerHalfSize}, // Right side
		{x + playerHalfSize, y + playerHalfSize, x - playerHalfSize, y + playerHalfSize}, // Bottom side
		{x - playerHalfSize, y + playerHalfSize, x - playerHalfSize, y - playerHalfSize}, // Left side
	}

	for _, obj := range objects {
		for _, wall := range obj.walls {
			// Check each side of the player's bounding box for intersection with the wall
			for _, side := range playerBox {
				if _, _, collides := intersection(side, wall); collides {
					return true
				}
			}
		}
	}

	return false // No collision detected
}
