package main

import (
	"embed"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tristanisham/cope/game"
	"github.com/tristanisham/cope/meta"
)

//go:embed assets/*
var Assets embed.FS

func main() {
	g := game.NewGame(Assets)

	ebiten.SetWindowSize(meta.ScreenWidth*2, meta.ScreenHeight*2)
	ebiten.SetWindowTitle("Cope")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
