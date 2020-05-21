package systems

import (
	"github.com/x-hgg-x/sokoban-go/lib/resources"

	w "github.com/x-hgg-x/goecsengine/world"
)

// MoveSystem moves player
func MoveSystem(world w.World) {
	moveUpAction := world.Resources.InputHandler.Actions[resources.MoveUpAction]
	moveDownAction := world.Resources.InputHandler.Actions[resources.MoveDownAction]
	moveLeftAction := world.Resources.InputHandler.Actions[resources.MoveLeftAction]
	moveRightAction := world.Resources.InputHandler.Actions[resources.MoveRightAction]

	moveUpFastAction := world.Resources.InputHandler.Actions[resources.MoveUpFastAction]
	moveDownFastAction := world.Resources.InputHandler.Actions[resources.MoveDownFastAction]
	moveLeftFastAction := world.Resources.InputHandler.Actions[resources.MoveLeftFastAction]
	moveRightFastAction := world.Resources.InputHandler.Actions[resources.MoveRightFastAction]

	switch {
	case moveUpAction || moveUpFastAction:
		resources.Move(world, resources.MovementUp)
	case moveDownAction || moveDownFastAction:
		resources.Move(world, resources.MovementDown)
	case moveLeftAction || moveLeftFastAction:
		resources.Move(world, resources.MovementLeft)
	case moveRightAction || moveRightFastAction:
		resources.Move(world, resources.MovementRight)
	}
}
