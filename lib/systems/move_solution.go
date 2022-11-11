package systems

import (
	"github.com/x-hgg-x/sokoban-go/lib/resources"

	w "github.com/x-hgg-x/goecsengine/world"
)

// MoveSolutionSystem moves using the solution
func MoveSolutionSystem(world w.World, movements []resources.MovementType) {
	gameResources := world.Resources.Game.(*resources.Game)

	nextStepSolutionAction := world.Resources.InputHandler.Actions[resources.NextStepSolutionAction]
	nextStepSolutionFastAction := world.Resources.InputHandler.Actions[resources.NextStepSolutionFastAction]
	previousStepSolutionAction := world.Resources.InputHandler.Actions[resources.PreviousStepSolutionAction]
	previousStepSolutionFastAction := world.Resources.InputHandler.Actions[resources.PreviousStepSolutionFastAction]

	if previousStepSolutionAction || previousStepSolutionFastAction {
		resources.Undo(world)
	}

	if nextStepSolutionAction || nextStepSolutionFastAction {
		if len(gameResources.Level.Movements) < len(movements) {
			resources.Move(world, movements[len(gameResources.Level.Movements)])
		}
	}
}
