package systems

import (
	"github.com/x-hgg-x/sokoban-go/lib/resources"

	w "github.com/x-hgg-x/goecsengine/world"
)

// UndoSystem undoes the last move
func UndoSystem(world w.World) {
	undoAction := world.Resources.InputHandler.Actions[resources.UndoAction]
	undoFastAction := world.Resources.InputHandler.Actions[resources.UndoFastAction]

	if undoAction || undoFastAction {
		resources.Undo(world)
	}
}
