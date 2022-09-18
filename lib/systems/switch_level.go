package systems

import (
	"github.com/x-hgg-x/sokoban-go/lib/resources"

	w "github.com/x-hgg-x/goecsengine/world"
)

// SwitchLevelSystem switches between levels
func SwitchLevelSystem(world w.World) {
	gameResources := world.Resources.Game.(*resources.Game)

	previousLevelAction := world.Resources.InputHandler.Actions[resources.PreviousLevelAction]
	previousLevelFastAction := world.Resources.InputHandler.Actions[resources.PreviousLevelFastAction]
	nextLevelAction := world.Resources.InputHandler.Actions[resources.NextLevelAction]
	nextLevelFastAction := world.Resources.InputHandler.Actions[resources.NextLevelFastAction]
	restartAction := world.Resources.InputHandler.Actions[resources.RestartAction]

	var newLevel int
	switch {
	case (previousLevelAction || previousLevelFastAction) && gameResources.Level.CurrentNum > 0:
		newLevel = gameResources.Level.CurrentNum - 1
	case (nextLevelAction || nextLevelFastAction) && gameResources.Level.CurrentNum < len(gameResources.Package.Levels)-1:
		newLevel = gameResources.Level.CurrentNum + 1
	case restartAction:
		gameResources.Level.Movements = []resources.MovementType{}
		gameResources.Level.Modified = true
		newLevel = gameResources.Level.CurrentNum
	default:
		return
	}

	world.Manager.DeleteAllEntities()
	resources.SaveLevel(world)
	resources.InitLevel(world, newLevel)
}
