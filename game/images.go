package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tristanisham/cope/meta"
)

var (
	fgImage       = ebiten.NewImage(meta.ScreenWidth, meta.ScreenHeight)
	bgImage       = ebiten.NewImage(meta.ScreenWidth, meta.ScreenHeight)
	shadowImage   = ebiten.NewImage(meta.ScreenWidth, meta.ScreenHeight)
	triangleImage = ebiten.NewImage(meta.ScreenWidth, meta.ScreenHeight)
	fovMask       = ebiten.NewImage(meta.ScreenWidth, meta.ScreenHeight)
)
