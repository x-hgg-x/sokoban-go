package systems

import (
	"fmt"

	"github.com/x-hgg-x/sokoban-go/lib/resources"

	ecs "github.com/x-hgg-x/goecs/v2"
	ec "github.com/x-hgg-x/goecsengine/components"
	w "github.com/x-hgg-x/goecsengine/world"
)

// InfoSystem sets game info
func InfoSystem(world w.World) {
	gameResources := world.Resources.Game.(*resources.Game)

	// Check the number of box on goal
	boxCount := 0
	boxOnGoalCount := 0

	for _, tile := range gameResources.Level.Grid.Data {
		if tile.Contains(resources.TileBox) {
			boxCount += 1

			if tile.Contains(resources.TileGoal) {
				boxOnGoalCount += 1
			}
		}
	}

	// Set text info
	world.Manager.Join(world.Components.Engine.Text, world.Components.Engine.UITransform).Visit(ecs.Visit(func(entity ecs.Entity) {
		text := world.Components.Engine.Text.Get(entity).(*ec.Text)

		switch text.ID {
		case "level":
			text.Text = fmt.Sprintf("LEVEL %d/%d", gameResources.Level.CurrentNum+1, len(gameResources.Package.Levels))
			if gameResources.Level.Modified {
				text.Text += "(*)"
			}
		case "box":
			text.Text = fmt.Sprintf("BOX: %d/%d", boxOnGoalCount, boxCount)
		case "step":
			text.Text = fmt.Sprintf("STEPS: %d", len(gameResources.Level.Movements))
		}
	}))

	// Finish level if all boxes are on goals
	if boxOnGoalCount == boxCount {
		gameResources.StateEvent = resources.StateEventLevelComplete
	}
}
