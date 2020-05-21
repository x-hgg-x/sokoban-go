package resources

import (
	gc "github.com/x-hgg-x/sokoban-go/lib/components"

	ecs "github.com/x-hgg-x/goecs/v2"
	w "github.com/x-hgg-x/goecsengine/world"
)

// Move executes a series of movements
func Move(world w.World, movements ...MovementType) {
	gameComponents := world.Components.Game.(*gc.Components)
	gameResources := world.Resources.Game.(*Game)

	firstPlayer := ecs.GetFirst(world.Manager.Join(gameComponents.Player, gameComponents.GridElement))
	if firstPlayer == nil {
		return
	}
	playerGridElement := gameComponents.GridElement.Get(ecs.Entity(*firstPlayer)).(*gc.GridElement)
	playerLine := &playerGridElement.Line
	playerCol := &playerGridElement.Col

	for _, movement := range movements {
		movement = GetSimpleMovement(movement)

		var directionLine, directionCol int
		switch movement {
		case MovementUp:
			directionLine, directionCol = -1, 0
		case MovementDown:
			directionLine, directionCol = 1, 0
		case MovementLeft:
			directionLine, directionCol = 0, -1
		case MovementRight:
			directionLine, directionCol = 0, 1
		}

		oneFrontLine := *playerLine + directionLine
		oneFrontCol := *playerCol + directionCol
		twoFrontLine := *playerLine + 2*directionLine
		twoFrontCol := *playerCol + 2*directionCol

		// Check grid edge
		if !(0 <= oneFrontLine && oneFrontLine < MaxHeight && 0 <= oneFrontCol && oneFrontCol < MaxWidth) {
			return
		}
		oneFrontTile := &gameResources.Level.Grid[oneFrontLine][oneFrontCol]

		// No move if a wall is ahead
		if oneFrontTile.Wall != nil {
			return
		}

		if box := oneFrontTile.Box; box != nil {
			// Check grid edge
			if !(0 <= twoFrontLine && twoFrontLine < MaxHeight && 0 <= twoFrontCol && twoFrontCol < MaxWidth) {
				return
			}
			twoFrontTile := &gameResources.Level.Grid[twoFrontLine][twoFrontCol]

			// No move if two boxes or a box and a wall are ahead
			if twoFrontTile.Box != nil || twoFrontTile.Wall != nil {
				return
			}
			boxGridElement := gameComponents.GridElement.Get(*box).(*gc.GridElement)
			boxGridElement.Line = twoFrontLine
			boxGridElement.Col = twoFrontCol
			twoFrontTile.Box = oneFrontTile.Box
			oneFrontTile.Box = nil
			movement = GetPushMovement(movement)
		}

		playerTile := &gameResources.Level.Grid[*playerLine][*playerCol]
		oneFrontTile.Player = playerTile.Player
		playerTile.Player = nil
		*playerLine = oneFrontLine
		*playerCol = oneFrontCol

		gameResources.Level.Movements = append(gameResources.Level.Movements, movement)
		gameResources.Level.Modified = true
	}
}
