package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tristanisham/cope/meta"
)

func (g *Game) handleMovement() {
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.px += 4
	}

	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.py += 4
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.px -= 4
	}

	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.py -= 4
	}

	// +1/-1 is to stop player before it reaches the border
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

	// add collision for walls
	
}