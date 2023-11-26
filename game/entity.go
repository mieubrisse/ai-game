package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mieubrisse/ai-game/utils"
)

// Describes an entity (e.g. the player, a monster, etc.)
type Entity struct {
	// Number of ticks needed to get where the entity is going
	slowness int

	locationX int
	locationY int

	// Represents the destination of the entity
	// If this is the same as the entity's X and Y, they're not in motion
	destinationX int
	destinationY int

	requiredTicks  int // The number of ticks for the entity to get where they're going
	completedTicks int // The number of ticks the entity has already completed

	facingX int // 0 if facing along the Y axis; -1 or 1 to indicate facing left or right
	facingY int // 0 if facing along the X axis; -1 or 1 to indicate facing up or down

	// Function that's called whenever the entity is not on the move, that returns a dX, dY indicating where the entity
	// should move
	getDesiredDirection func(self *Entity) (dX, dY int)
}

func (e *Entity) IsOnTheMove() bool {
	return e.locationX != e.destinationX || e.locationY != e.destinationY
}

func (e *Entity) Update() {
	// If the player arrived in the last tick, set them to no longer be on the move
	if e.IsOnTheMove() && e.completedTicks >= e.requiredTicks {
		e.locationX = e.destinationX
		e.locationY = e.destinationY
	}

	// The player is arrived; they can move again
	if !e.IsOnTheMove() {
		// Determine the player's target and set it, if so
		dX, dY := e.getDesiredDirection(e)
		if dX != 0 {
			e.facingX = dX
		}
		// TODO Y-facing

		targetX := utils.Coerce(e.locationX+dX, 0, WidthCells)
		targetY := utils.Coerce(e.locationY+dY, 0, HeightCells)

		// If the player has a target != their current location, set them up to be on the move
		if e.locationX != targetX || e.locationY != targetY {
			e.destinationX = targetX
			e.destinationY = targetY
			e.requiredTicks = e.slowness
			e.completedTicks = 0
		}
	}

	// If the player is now on the move, start moving them
	if e.IsOnTheMove() {
		e.completedTicks++
	}
}

func (e *Entity) Draw(screen *ebiten.Image) {
	// If the player is on the move, they'll be in progress to their destination so we need
	// to prepare to translate the sprite appropriately
	playerImage := mageImage
	var onTheMoveXOffsetPixels, onTheMoveYOffsetPixels float64
	if e.IsOnTheMove() {
		progressPercentage := float64(e.completedTicks) / float64(e.requiredTicks)
		pixelOffset := float64(PixelsPerCellSide) * progressPercentage
		onTheMoveXOffsetPixels = float64(e.destinationX-e.locationX) * pixelOffset
		onTheMoveYOffsetPixels = float64(e.destinationY-e.locationY) * pixelOffset

		if (progressPercentage >= 0.25 && progressPercentage < 0.50) || progressPercentage >= 0.75 {
			playerImage = mage2Image
		}
	}

	drawImageOptions := &ebiten.DrawImageOptions{}
	drawImageOptions.GeoM.Translate(
		float64(e.locationX*PixelsPerCellSide)+onTheMoveXOffsetPixels,
		float64(e.locationY*PixelsPerCellSide)+onTheMoveYOffsetPixels,
	)

	// text.Draw()

	screen.DrawImage(playerImage, drawImageOptions)
}
