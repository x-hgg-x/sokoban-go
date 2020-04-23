package systems

import (
	gc "github.com/x-hgg-x/sokoban-go/lib/components"
	"github.com/x-hgg-x/sokoban-go/lib/resources"

	ecs "github.com/x-hgg-x/goecs/v2"
	w "github.com/x-hgg-x/goecsengine/world"
)

// UndoSystem undoes the last move
func UndoSystem(world w.World) {
	gameComponents := world.Components.Game.(*gc.Components)
	gameResources := world.Resources.Game.(*resources.Game)

	undoAction := world.Resources.InputHandler.Actions[resources.UndoAction]
	undoFastAction := world.Resources.InputHandler.Actions[resources.UndoFastAction]

	firstPlayer := ecs.GetFirst(world.Manager.Join(gameComponents.Player, gameComponents.GridElement))
	if firstPlayer == nil {
		return
	}
	playerGridElement := gameComponents.GridElement.Get(ecs.Entity(*firstPlayer)).(*gc.GridElement)

	if (undoAction || undoFastAction) && len(gameResources.Movements) > 0 {
		switch gameResources.Movements[len(gameResources.Movements)-1] {
		case resources.MovementUp:
			undo(world, false, &playerGridElement.Line, &playerGridElement.Col, -1, 0)
		case resources.MovementDown:
			undo(world, false, &playerGridElement.Line, &playerGridElement.Col, 1, 0)
		case resources.MovementLeft:
			undo(world, false, &playerGridElement.Line, &playerGridElement.Col, 0, -1)
		case resources.MovementRight:
			undo(world, false, &playerGridElement.Line, &playerGridElement.Col, 0, 1)
		case resources.MovementUpPush:
			undo(world, true, &playerGridElement.Line, &playerGridElement.Col, -1, 0)
		case resources.MovementDownPush:
			undo(world, true, &playerGridElement.Line, &playerGridElement.Col, 1, 0)
		case resources.MovementLeftPush:
			undo(world, true, &playerGridElement.Line, &playerGridElement.Col, 0, -1)
		case resources.MovementRightPush:
			undo(world, true, &playerGridElement.Line, &playerGridElement.Col, 0, 1)
		}
		gameResources.Movements = gameResources.Movements[:len(gameResources.Movements)-1]
	}
}

func undo(world w.World, boxPush bool, playerLine, playerCol *int, directionLine, directionCol int) {
	gameComponents := world.Components.Game.(*gc.Components)
	gameResources := world.Resources.Game.(*resources.Game)

	playerTile := &gameResources.Grid[*playerLine][*playerCol]
	oneFrontTile := &gameResources.Grid[*playerLine+directionLine][*playerCol+directionCol]

	if boxPush {
		boxGridElement := gameComponents.GridElement.Get(*oneFrontTile.Box).(*gc.GridElement)
		boxGridElement.Line = *playerLine
		boxGridElement.Col = *playerCol
		playerTile.Box = oneFrontTile.Box
		oneFrontTile.Box = nil
	}

	oneBackTile := &gameResources.Grid[*playerLine-directionLine][*playerCol-directionCol]
	oneBackTile.Player = playerTile.Player
	playerTile.Player = nil
	*playerLine -= directionLine
	*playerCol -= directionCol
}
