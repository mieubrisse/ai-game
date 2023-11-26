package game

import (
	"bytes"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mieubrisse/ai-game/resources"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"image"
	_ "image/png"
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
	ecs *ecs.ECS
}

func NewGame() *Game {
	world := donburi.NewWorld()

	playerEntity := world.Create(Mob, PlayerTag)
	playerEntry := world.Entry(playerEntity)
	playerData := Mob.Get(playerEntry)
	playerData.locationX = 10
	playerData.locationY = 10
	playerData.facingX = 1
	playerData.slowness = playerSlowness
	playerData.getDesiredDirection = func(world donburi.World, self *donburi.Entry) (dX, dY int) {
		return getPlayerDesiredDirection()
	}

	enemyId := world.Create(Mob)
	enemyEntry := world.Entry(enemyId)
	enemyData := Mob.Get(enemyEntry)
	enemyData.locationX = 0
	enemyData.locationY = 0
	enemyData.facingX = 1
	enemyData.slowness = enemySlowness
	enemyData.getDesiredDirection = func(world donburi.World, self *donburi.Entry) (dX, dY int) {
		return huntTargetEntity(self, playerEntry)
	}

	gameEcs := ecs.NewECS(world).
		AddSystem(AnimateActors)

	return &Game{
		ecs: gameEcs,
	}
}

func (g *Game) Update() error {
	g.ecs.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	query := donburi.NewQuery(filter.Contains(Mob))
	query.Each(g.ecs.World, func(entry *donburi.Entry) {
		mobData := Mob.Get(entry)

		// If the player is on the move, they'll be in progress to their destination so we need
		// to prepare to translate the sprite appropriately
		playerImage := mageImage
		var onTheMoveXOffsetPixels, onTheMoveYOffsetPixels float64
		if IsMobOnTheMove(mobData) {
			progressPercentage := float64(mobData.completedTicks) / float64(mobData.requiredTicks)
			pixelOffset := float64(PixelsPerCellSide) * progressPercentage
			onTheMoveXOffsetPixels = float64(mobData.destinationX-mobData.locationX) * pixelOffset
			onTheMoveYOffsetPixels = float64(mobData.destinationY-mobData.locationY) * pixelOffset

			if (progressPercentage >= 0.25 && progressPercentage < 0.50) || progressPercentage >= 0.75 {
				playerImage = mage2Image
			}
		}

		drawImageOptions := &ebiten.DrawImageOptions{}
		drawImageOptions.GeoM.Translate(
			float64(mobData.locationX*PixelsPerCellSide)+onTheMoveXOffsetPixels,
			float64(mobData.locationY*PixelsPerCellSide)+onTheMoveYOffsetPixels,
		)

		// text.Draw()

		screen.DrawImage(playerImage, drawImageOptions)
	})
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

func huntTargetEntity(self, target *donburi.Entry) (dX, dY int) {
	selfMobData := Mob.Get(self)
	targetMobData := Mob.Get(target)

	if targetMobData.locationX > selfMobData.locationX {
		return 1, 0
	}
	if targetMobData.locationY > selfMobData.locationY {
		return 0, 1
	}
	if targetMobData.locationX < selfMobData.locationX {
		return -1, 0
	}
	if targetMobData.locationY < selfMobData.locationY {
		return 0, -1
	}
	return 0, 0
}
