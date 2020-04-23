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
	case (previousLevelAction || previousLevelFastAction) && gameResources.CurrentLevel > 0:
		newLevel = gameResources.CurrentLevel - 1
	case (nextLevelAction || nextLevelFastAction) && gameResources.CurrentLevel < gameResources.LevelCount-1:
		newLevel = gameResources.CurrentLevel + 1
	case restartAction:
		newLevel = gameResources.CurrentLevel
	default:
		return
	}

	world.Manager.DeleteAllEntities()
	resources.InitLevel(world, newLevel)
}
