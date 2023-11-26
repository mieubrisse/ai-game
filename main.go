package main

import (
	"github.com/mieubrisse/ai-game/game"
	"log"

	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(
		3*game.WidthCells*game.PixelsPerCellSide,
		3*game.HeightCells*game.PixelsPerCellSide,
	)
	ebiten.SetWindowTitle("AI game")

	if err := ebiten.RunGame(&game.Game{}); err != nil {
		log.Fatal(err)
	}
}
