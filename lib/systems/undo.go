package systems

import (
	"github.com/x-hgg-x/sokoban-go/lib/resources"

	w "github.com/x-hgg-x/goecsengine/world"
)

// UndoSystem undoes the last move
func UndoSystem(world w.World) {
	gameResources := world.Resources.Game.(*resources.Game)

	undoAction := world.Resources.InputHandler.Actions[resources.UndoAction]
	undoFastAction := world.Resources.InputHandler.Actions[resources.UndoFastAction]

	if (undoAction || undoFastAction) && len(gameResources.Level.Movements) > 0 {
		undo(world)
	}
}

func undo(world w.World) {
	gameResources := world.Resources.Game.(*resources.Game)
	levelWidth := gameResources.Level.Grid.NCols

	playerIndex := -1
	for iTile, tile := range gameResources.Level.Grid.Data {
		if tile.Contains(resources.TilePlayer) {
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
	case resources.MovementUp:
		boxPush, directionLine, directionCol = false, -1, 0
	case resources.MovementDown:
		boxPush, directionLine, directionCol = false, 1, 0
	case resources.MovementLeft:
		boxPush, directionLine, directionCol = false, 0, -1
	case resources.MovementRight:
		boxPush, directionLine, directionCol = false, 0, 1
	case resources.MovementUpPush:
		boxPush, directionLine, directionCol = true, -1, 0
	case resources.MovementDownPush:
		boxPush, directionLine, directionCol = true, 1, 0
	case resources.MovementLeftPush:
		boxPush, directionLine, directionCol = true, 0, -1
	case resources.MovementRightPush:
		boxPush, directionLine, directionCol = true, 0, 1
	}

	if boxPush {
		oneFrontTile := gameResources.Level.Grid.Get(playerLine+directionLine, playerCol+directionCol)
		playerTile.Set(resources.TileBox)
		oneFrontTile.Remove(resources.TileBox)
	}

	oneBackTile := gameResources.Level.Grid.Get(playerLine-directionLine, playerCol-directionCol)
	oneBackTile.Set(resources.TilePlayer)
	playerTile.Remove(resources.TilePlayer)

	gameResources.Level.Movements = gameResources.Level.Movements[:len(gameResources.Level.Movements)-1]
	gameResources.Level.Modified = true
}
