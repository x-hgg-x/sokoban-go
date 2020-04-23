package systems

import (
	gc "github.com/x-hgg-x/sokoban-go/lib/components"
	"github.com/x-hgg-x/sokoban-go/lib/resources"

	ecs "github.com/x-hgg-x/goecs/v2"
	w "github.com/x-hgg-x/goecsengine/world"
)

// MoveSystem moves player
func MoveSystem(world w.World) {
	gameComponents := world.Components.Game.(*gc.Components)

	moveUpAction := world.Resources.InputHandler.Actions[resources.MoveUpAction]
	moveDownAction := world.Resources.InputHandler.Actions[resources.MoveDownAction]
	moveLeftAction := world.Resources.InputHandler.Actions[resources.MoveLeftAction]
	moveRightAction := world.Resources.InputHandler.Actions[resources.MoveRightAction]

	moveUpFastAction := world.Resources.InputHandler.Actions[resources.MoveUpFastAction]
	moveDownFastAction := world.Resources.InputHandler.Actions[resources.MoveDownFastAction]
	moveLeftFastAction := world.Resources.InputHandler.Actions[resources.MoveLeftFastAction]
	moveRightFastAction := world.Resources.InputHandler.Actions[resources.MoveRightFastAction]

	firstPlayer := ecs.GetFirst(world.Manager.Join(gameComponents.Player, gameComponents.GridElement))
	if firstPlayer == nil {
		return
	}
	playerGridElement := gameComponents.GridElement.Get(ecs.Entity(*firstPlayer)).(*gc.GridElement)

	// Move up
	if moveUpAction || moveUpFastAction {
		move(world, resources.MovementUp, &playerGridElement.Line, &playerGridElement.Col, -1, 0)
	}

	// Move down
	if moveDownAction || moveDownFastAction {
		move(world, resources.MovementDown, &playerGridElement.Line, &playerGridElement.Col, 1, 0)
	}

	// Move left
	if moveLeftAction || moveLeftFastAction {
		move(world, resources.MovementLeft, &playerGridElement.Line, &playerGridElement.Col, 0, -1)
	}

	// Move right
	if moveRightAction || moveRightFastAction {
		move(world, resources.MovementRight, &playerGridElement.Line, &playerGridElement.Col, 0, 1)
	}
}

func move(world w.World, movement resources.MovementType, playerLine, playerCol *int, directionLine, directionCol int) {
	gameComponents := world.Components.Game.(*gc.Components)
	gameResources := world.Resources.Game.(*resources.Game)

	oneFrontLine := *playerLine + directionLine
	oneFrontCol := *playerCol + directionCol
	twoFrontLine := *playerLine + 2*directionLine
	twoFrontCol := *playerCol + 2*directionCol

	// Check grid edge
	if !(0 <= oneFrontLine && oneFrontLine < resources.MaxHeight && 0 <= oneFrontCol && oneFrontCol < resources.MaxWidth) {
		return
	}
	oneFrontTile := &gameResources.Grid[oneFrontLine][oneFrontCol]

	// No move if a wall is ahead
	if oneFrontTile.Wall != nil {
		return
	}

	if box := oneFrontTile.Box; box != nil {
		// Check grid edge
		if !(0 <= twoFrontLine && twoFrontLine < resources.MaxHeight && 0 <= twoFrontCol && twoFrontCol < resources.MaxWidth) {
			return
		}
		twoFrontTile := &gameResources.Grid[twoFrontLine][twoFrontCol]

		// No move if two boxes or a box and a wall are ahead
		if twoFrontTile.Box != nil || twoFrontTile.Wall != nil {
			return
		}
		twoFrontTile.Box = oneFrontTile.Box
		oneFrontTile.Box = nil
		boxGridElement := gameComponents.GridElement.Get(*box).(*gc.GridElement)
		boxGridElement.Line = twoFrontLine
		boxGridElement.Col = twoFrontCol
		movement += 4
	}

	playerTile := &gameResources.Grid[*playerLine][*playerCol]
	oneFrontTile.Player = playerTile.Player
	playerTile.Player = nil
	*playerLine = oneFrontLine
	*playerCol = oneFrontCol

	gameResources.Movements = append(gameResources.Movements, movement)
}
