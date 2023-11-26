package game

import (
	"bytes"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mieubrisse/ai-game/resources"
	"image"
)

const (
	WidthCells        = 40
	HeightCells       = 20
	PixelsPerCellSide = 10

	// How many ticks it takes for the player character to move a cell
	playerSlowness = 20

	// How many ticks it takes for the player character to move a cell
	enemySlowness = 30
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
	entities []*Entity
}

func NewGame() *Game {
	player := &Entity{
		getDesiredDirection: func(self *Entity) (dX, dY int) {
			return getPlayerDesiredDirection()
		},
		slowness: playerSlowness,
	}
	entities := []*Entity{
		player,
		{
			slowness: enemySlowness,
			getDesiredDirection: func(self *Entity) (dX, dY int) {
				return huntTargetEntity(self, player)
			},
		},
	}
	return &Game{
		entities: entities,
	}
}

func (g *Game) Update() error {
	for _, entity := range g.entities {
		entity.Update()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, entity := range g.entities {
		entity.Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return WidthCells * PixelsPerCellSide, HeightCells * PixelsPerCellSide
}

func getPlayerDesiredDirection() (dX, dY int) {
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

func loadImage(imageBytes []byte) *ebiten.Image {
	img, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		panic(err)
	}
	return ebiten.NewImageFromImage(img)
}

func huntTargetEntity(self *Entity, target *Entity) (dX, dY int) {
	if target.locationX > self.locationX {
		return 1, 0
	}
	if target.locationY > self.locationY {
		return 0, 1
	}
	if target.locationX < self.locationX {
		return -1, 0
	}
	if target.locationY < self.locationY {
		return 0, -1
	}
	return 0, 0
}
