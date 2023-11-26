package main

import (
	"bytes"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/mieubrisse/ai-game/resources"
	"image"
	"log"

	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	widthCells        = 40
	heightCells       = 20
	pixelsPerCellSide = 10

	// How many ticks it takes for the player character to move a cell
	playerSlowness = 20
)

var (
	mageImage      *ebiten.Image
	mage2Image     *ebiten.Image
	magicBallImage *ebiten.Image
)

func init() {
	mageImage = loadImage(resources.MagePNG)
	mage2Image = loadImage(resources.Mage2PNG)
	magicBallImage = loadImage(resources.MagicBallPNG)
}

type Game struct {
	// Represents the destination of the player
	// If this is the same as the player's X and Y, they're not in motion
	destinationX   int
	destinationY   int
	requiredTicks  int // The number of ticks for the player to get where they're going
	completedTicks int // The number of ticks the player has already completed

	facingX int // 0 if facing along the Y axis; -1 or 1 to indicate facing left or right
	facingY int // 0 if facing along the X axis; -1 or 1 to indicate facing up or down

	// TODO make this an entity
	cells [][]bool

	playerX int
	playerY int

	mageDude *ebiten.Image

	numLeft  int
	numRight int
}

func (g *Game) Update() error {
	// If the player arrived in the last tick, set them to no longer be on the move
	if g.isPlayerOnTheMove() && g.completedTicks >= g.requiredTicks {
		g.playerX = g.destinationX
		g.playerY = g.destinationY
	}

	// The player is arrived; they can move again
	if !g.isPlayerOnTheMove() {
		// Determine the player's target and set it, if so
		dX, dY := getPlayerTranslation()
		if dX != 0 {
			g.facingX = dX
		}
		// TODO Y-facing

		targetX := coerce(g.playerX+dX, 0, widthCells)
		targetY := coerce(g.playerY+dY, 0, heightCells)

		// If the player has a target != their current location, set them up to be on the move
		if g.playerX != targetX || g.playerY != targetY {
			g.destinationX = targetX
			g.destinationY = targetY
			g.requiredTicks = playerSlowness
			g.completedTicks = 0
		}
	}

	// If the player is now on the move, start moving them
	if g.isPlayerOnTheMove() {
		g.completedTicks++
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// str := fmt.Sprintf("Left: %d, Right: %d", g.numLeft, g.numRight)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Facing %d, %d", g.facingX, g.facingY))

	// If the player is on the move, they'll be in progress to their destination so we need
	// to prepare to translate the sprite appropriately
	playerImage := mageImage
	var onTheMoveXOffsetPixels, onTheMoveYOffsetPixels float64
	if g.isPlayerOnTheMove() {
		progressPercentage := float64(g.completedTicks) / float64(g.requiredTicks)
		pixelOffset := float64(pixelsPerCellSide) * progressPercentage
		onTheMoveXOffsetPixels = float64(g.destinationX-g.playerX) * pixelOffset
		onTheMoveYOffsetPixels = float64(g.destinationY-g.playerY) * pixelOffset

		if (progressPercentage >= 0.25 && progressPercentage < 0.50) || progressPercentage >= 0.75 {
			playerImage = mage2Image
		}
	}

	drawImageOptions := &ebiten.DrawImageOptions{}
	drawImageOptions.GeoM.Translate(
		float64(g.playerX*pixelsPerCellSide)+onTheMoveXOffsetPixels,
		float64(g.playerY*pixelsPerCellSide)+onTheMoveYOffsetPixels,
	)

	screen.DrawImage(playerImage, drawImageOptions)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return widthCells * pixelsPerCellSide, heightCells * pixelsPerCellSide
}

func main() {
	ebiten.SetWindowSize(
		3*widthCells*pixelsPerCellSide,
		3*heightCells*pixelsPerCellSide,
	)
	ebiten.SetWindowTitle("AI game")

	cells := make([][]bool, widthCells)
	for i := 0; i < widthCells; i++ {
		cells[i] = make([]bool, heightCells)
	}
	cells[0][0] = true

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

func coerce(input, min, max int) int {
	if max < min {
		newMax := min
		min = max
		max = newMax
	}
	if input < min {
		return min
	}
	if input > max {
		return max
	}
	return input
}

func getPlayerTranslation() (dX, dY int) {
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		return 1, 0
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		return 0, 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		return -1, 0
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		return 0, -1
	}
	return 0, 0
}

func (g *Game) isPlayerOnTheMove() bool {
	return g.playerX != g.destinationX || g.playerY != g.destinationY
}

func loadImage(imageBytes []byte) *ebiten.Image {
	img, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		panic(err)
	}
	return ebiten.NewImageFromImage(img)
}