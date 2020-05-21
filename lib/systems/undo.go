package systems

import (
	gc "github.com/x-hgg-x/sokoban-go/lib/components"
	"github.com/x-hgg-x/sokoban-go/lib/resources"

	ecs "github.com/x-hgg-x/goecs/v2"
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
	gameComponents := world.Components.Game.(*gc.Components)
	gameResources := world.Resources.Game.(*resources.Game)

	firstPlayer := ecs.GetFirst(world.Manager.Join(gameComponents.Player, gameComponents.GridElement))
	if firstPlayer == nil {
		return
	}
	playerGridElement := gameComponents.GridElement.Get(ecs.Entity(*firstPlayer)).(*gc.GridElement)
	playerLine := &playerGridElement.Line
	playerCol := &playerGridElement.Col

	var boxPush bool
	var directionLine, directionCol int
	switch gameResources.Level.Movements[len(gameResources.Level.Movements)-1] {
	case resources.MovementUp:
		boxPush = false
		directionLine, directionCol = -1, 0
	case resources.MovementDown:
		boxPush = false
		directionLine, directionCol = 1, 0
	case resources.MovementLeft:
		boxPush = false
		directionLine, directionCol = 0, -1
	case resources.MovementRight:
		boxPush = false
		directionLine, directionCol = 0, 1
	case resources.MovementUpPush:
		boxPush = true
		directionLine, directionCol = -1, 0
	case resources.MovementDownPush:
		boxPush = true
		directionLine, directionCol = 1, 0
	case resources.MovementLeftPush:
		boxPush = true
		directionLine, directionCol = 0, -1
	case resources.MovementRightPush:
		boxPush = true
		directionLine, directionCol = 0, 1
	}

	playerTile := &gameResources.Level.Grid[*playerLine][*playerCol]
	oneFrontTile := &gameResources.Level.Grid[*playerLine+directionLine][*playerCol+directionCol]

	if boxPush {
		boxGridElement := gameComponents.GridElement.Get(*oneFrontTile.Box).(*gc.GridElement)
		boxGridElement.Line = *playerLine
		boxGridElement.Col = *playerCol
		playerTile.Box = oneFrontTile.Box
		oneFrontTile.Box = nil
	}

	oneBackTile := &gameResources.Level.Grid[*playerLine-directionLine][*playerCol-directionCol]
	oneBackTile.Player = playerTile.Player
	playerTile.Player = nil
	*playerLine -= directionLine
	*playerCol -= directionCol

	gameResources.Level.Movements = gameResources.Level.Movements[:len(gameResources.Level.Movements)-1]
	gameResources.Level.Modified = true
}
