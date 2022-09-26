package systems

import (
	gc "github.com/x-hgg-x/sokoban-go/lib/components"
	"github.com/x-hgg-x/sokoban-go/lib/resources"

	ecs "github.com/x-hgg-x/goecs/v2"
	w "github.com/x-hgg-x/goecsengine/world"
)

func GridUpdateSystem(world w.World) {
	gameComponents := world.Components.Game.(*gc.Components)
	gameResources := world.Resources.Game.(*resources.Game)

	playerIndex := -1
	boxIndices := []int{}
	for iTile, tile := range gameResources.Level.Grid.Data {
		switch {
		case tile.Contains(resources.TilePlayer):
			playerIndex = iTile
		case tile.Contains(resources.TileBox):
			boxIndices = append(boxIndices, iTile)
		}
	}

	levelWidth := gameResources.Level.Grid.NCols
	levelHeight := gameResources.Level.Grid.NRows

	paddingRow := (gameResources.GridLayout.Height - levelHeight) / 2
	paddingCol := (gameResources.GridLayout.Width - levelWidth) / 2

	world.Manager.Join(gameComponents.GridElement).Visit(ecs.Visit(func(entity ecs.Entity) {
		switch {
		case entity.HasComponent(gameComponents.Player):
			gridElement := gameComponents.GridElement.Get(entity).(*gc.GridElement)
			gridElement.Line = paddingRow + playerIndex/levelWidth
			gridElement.Col = paddingCol + playerIndex%levelWidth

		case entity.HasComponent(gameComponents.Box):
			gridElement := gameComponents.GridElement.Get(entity).(*gc.GridElement)
			boxIndex := boxIndices[0]
			boxIndices = boxIndices[1:]
			gridElement.Line = paddingRow + boxIndex/levelWidth
			gridElement.Col = paddingCol + boxIndex%levelWidth
		}
	}))
}
