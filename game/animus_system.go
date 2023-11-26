package game

import (
	"github.com/mieubrisse/ai-game/utils"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

// This is a System that animates the entities in the world
func AnimateActors(theEcs *ecs.ECS) {
	query := donburi.NewQuery(filter.Contains(Mob))
	world := theEcs.World

	query.Each(world, func(entry *donburi.Entry) {
		animateActor(world, entry)
	})
}

// TODO split out the position-calculating from the actual movement
func animateActor(world donburi.World, self *donburi.Entry) {
	mobData := Mob.Get(self)

	// If the player arrived in the last tick, set them to no longer be on the move
	if IsMobOnTheMove(mobData) && mobData.completedTicks >= mobData.requiredTicks {
		mobData.locationX = mobData.destinationX
		mobData.locationY = mobData.destinationY
	}

	// The player is arrived; they can move again
	if !IsMobOnTheMove(mobData) {
		// Determine the player's target and set it, if so
		dX, dY := mobData.getDesiredDirection(world, self)
		if dX != 0 {
			mobData.facingX = dX
		}
		// TODO Y-facing

		targetX := utils.Coerce(mobData.locationX+dX, 0, WidthCells)
		targetY := utils.Coerce(mobData.locationY+dY, 0, HeightCells)

		// If the player has a target != their current location, set them up to be on the move
		if mobData.locationX != targetX || mobData.locationY != targetY {
			mobData.destinationX = targetX
			mobData.destinationY = targetY
			mobData.requiredTicks = mobData.slowness
			mobData.completedTicks = 0
		}
	}

	// If the player is now on the move, start moving them
	if IsMobOnTheMove(mobData) {
		mobData.completedTicks++
	}
}
