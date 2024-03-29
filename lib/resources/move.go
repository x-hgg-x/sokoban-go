package resources

import (
	w "github.com/x-hgg-x/goecsengine/world"
)

// Move executes a series of movements
func Move(world w.World, movements ...MovementType) {
	gameResources := world.Resources.Game.(*Game)

	levelWidth := gameResources.Level.Grid.NCols
	levelHeight := gameResources.Level.Grid.NRows

	playerIndex := -1
	for iTile, tile := range gameResources.Level.Grid.Data {
		if tile.Contains(TilePlayer) {
			playerIndex = iTile
			break
		}
	}

	playerTile := &gameResources.Level.Grid.Data[playerIndex]
	playerLine := playerIndex / levelWidth
	playerCol := playerIndex % levelWidth

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

		oneFrontLine := playerLine + directionLine
		oneFrontCol := playerCol + directionCol
		twoFrontLine := playerLine + 2*directionLine
		twoFrontCol := playerCol + 2*directionCol

		// Check grid edge
		if !(0 <= oneFrontLine && oneFrontLine < levelHeight && 0 <= oneFrontCol && oneFrontCol < levelWidth) {
			return
		}
		oneFrontTile := gameResources.Level.Grid.Get(oneFrontLine, oneFrontCol)

		// No move if a wall is ahead
		if oneFrontTile.Contains(TileWall) {
			return
		}

		if oneFrontTile.Contains(TileBox) {
			// Check grid edge
			if !(0 <= twoFrontLine && twoFrontLine < levelHeight && 0 <= twoFrontCol && twoFrontCol < levelWidth) {
				return
			}
			twoFrontTile := gameResources.Level.Grid.Get(twoFrontLine, twoFrontCol)

			// No move if two boxes or a box and a wall are ahead
			if twoFrontTile.ContainsAny(TileBox | TileWall) {
				return
			}

			twoFrontTile.Set(TileBox)
			oneFrontTile.Remove(TileBox)

			movement = GetPushMovement(movement)
		}

		oneFrontTile.Set(TilePlayer)
		playerTile.Remove(TilePlayer)

		playerTile = oneFrontTile
		playerLine += directionLine
		playerCol += directionCol

		gameResources.Level.Movements = append(gameResources.Level.Movements, movement)
		gameResources.Level.Modified = true
	}
}

// Undo undoes the last move
func Undo(world w.World) {
	gameResources := world.Resources.Game.(*Game)
	levelWidth := gameResources.Level.Grid.NCols

	if len(gameResources.Level.Movements) == 0 {
		return
	}

	playerIndex := -1
	for iTile, tile := range gameResources.Level.Grid.Data {
		if tile.Contains(TilePlayer) {
			playerIndex = iTile
			break
		}
	}

	playerTile := &gameResources.Level.Grid.Data[playerIndex]
	playerLine := playerIndex / levelWidth
	playerCol := playerIndex % levelWidth

	var boxPush bool
	var directionLine, directionCol int
	switch gameResources.Level.Movements[len(gameResources.Level.Movements)-1] {
	case MovementUp:
		boxPush, directionLine, directionCol = false, -1, 0
	case MovementDown:
		boxPush, directionLine, directionCol = false, 1, 0
	case MovementLeft:
		boxPush, directionLine, directionCol = false, 0, -1
	case MovementRight:
		boxPush, directionLine, directionCol = false, 0, 1
	case MovementUpPush:
		boxPush, directionLine, directionCol = true, -1, 0
	case MovementDownPush:
		boxPush, directionLine, directionCol = true, 1, 0
	case MovementLeftPush:
		boxPush, directionLine, directionCol = true, 0, -1
	case MovementRightPush:
		boxPush, directionLine, directionCol = true, 0, 1
	}

	if boxPush {
		oneFrontTile := gameResources.Level.Grid.Get(playerLine+directionLine, playerCol+directionCol)
		playerTile.Set(TileBox)
		oneFrontTile.Remove(TileBox)
	}

	oneBackTile := gameResources.Level.Grid.Get(playerLine-directionLine, playerCol-directionCol)
	oneBackTile.Set(TilePlayer)
	playerTile.Remove(TilePlayer)

	gameResources.Level.Movements = gameResources.Level.Movements[:len(gameResources.Level.Movements)-1]
	gameResources.Level.Modified = true
}
