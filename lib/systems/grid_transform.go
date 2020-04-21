package systems

import (
	gc "github.com/x-hgg-x/sokoban-go/lib/components"

	ecs "github.com/x-hgg-x/goecs/v2"
	ec "github.com/x-hgg-x/goecsengine/components"
	w "github.com/x-hgg-x/goecsengine/world"
)

const (
	transformOffsetX = 0
	transformOffsetY = -40
)

// GridTransformSystem sets transform for grid elements
func GridTransformSystem(world w.World) {
	gameComponents := world.Components.Game.(*gc.Components)

	world.Manager.Join(gameComponents.GridElement, world.Components.Engine.SpriteRender, world.Components.Engine.Transform).Visit(ecs.Visit(func(entity ecs.Entity) {
		gridElement := gameComponents.GridElement.Get(entity).(*gc.GridElement)
		elementSpriteRender := world.Components.Engine.SpriteRender.Get(entity).(*ec.SpriteRender)
		elementTranslation := &world.Components.Engine.Transform.Get(entity).(*ec.Transform).Translation

		screenHeight := float64(world.Resources.ScreenDimensions.Height)
		elementSprite := elementSpriteRender.SpriteSheet.Sprites[elementSpriteRender.SpriteNumber]

		elementTranslation.X = float64(gridElement.PosCol*elementSprite.Width) + float64(elementSprite.Width)/2 + transformOffsetX
		elementTranslation.Y = screenHeight - float64(gridElement.PosLine*elementSprite.Height) - float64(elementSprite.Height)/2 + transformOffsetY
	}))
}
