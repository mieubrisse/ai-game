package game

import "github.com/yohamta/donburi"

var PlayerTag = donburi.NewTag()

type MobData struct {
	// Number of ticks needed for the entity to get anywhere it's going
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
	getDesiredDirection func(world donburi.World, self *donburi.Entry) (dX, dY int)
}

var Mob = donburi.NewComponentType[MobData]()

func IsMobOnTheMove(mobData *MobData) bool {
	return mobData.locationX != mobData.destinationX || mobData.locationY != mobData.destinationY
}

// Represents data about which entity is being targeted
// Used exclusively for enemies
type HunterData struct {
	// Nil if no entity is being targeted currently
	target *donburi.Entity
}

var Hunter = donburi.NewComponentType[HunterData]()

/*
// X & Y position data
type PositionData struct {
	x int
	y int
}

var Position = donburi.NewComponentType[PositionData]()

// A unit vector that the entity is facing in
type FaceDirectionData struct {
	x int
	y int
}

var FaceDirection = donburi.NewComponentType[FaceDirectionData]()

//
type DesiredDestinationData struct {
	x int
	y int
}

// A component representing the inner will of the entity to move
type AnimusData struct {
	animus func() (dX, dY int)
}

var Animus = donburi.NewComponentType[AnimusData]()


*/
