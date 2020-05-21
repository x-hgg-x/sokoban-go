package systems

import (
	"fmt"

	gc "github.com/x-hgg-x/sokoban-go/lib/components"
	"github.com/x-hgg-x/sokoban-go/lib/resources"

	ecs "github.com/x-hgg-x/goecs/v2"
	ec "github.com/x-hgg-x/goecsengine/components"
	w "github.com/x-hgg-x/goecsengine/world"

	"golang.org/x/image/font"
)

type levelTextData struct {
	Texts    []string
	FontFace font.Face
	Pivot    string
	YOffsets []int
}

func (l *levelTextData) computeDotOffsets() {
	for _, text := range l.Texts {
		_, y, _ := ec.ComputeDotOffset(text, l.FontFace, l.Pivot)
		l.YOffsets = append(l.YOffsets, y)
	}
}

var levelText *levelTextData

// InfoSystem sets game info
func InfoSystem(world w.World) {
	gameComponents := world.Components.Game.(*gc.Components)
	gameResources := world.Resources.Game.(*resources.Game)

	// Check the number of box on goal
	boxSet := world.Manager.Join(gameComponents.Box, gameComponents.GridElement)
	boxCount := boxSet.Size()
	boxOnGoalCount := 0

	boxSet.Visit(ecs.Visit(func(entity ecs.Entity) {
		boxGridElement := gameComponents.GridElement.Get(entity).(*gc.GridElement)
		if gameResources.Level.Grid[boxGridElement.Line][boxGridElement.Col].Goal != nil {
			boxOnGoalCount++
		}
	}))

	// Set text info
	world.Manager.Join(world.Components.Engine.Text, world.Components.Engine.UITransform).Visit(ecs.Visit(func(entity ecs.Entity) {
		text := world.Components.Engine.Text.Get(entity).(*ec.Text)
		uiTransform := world.Components.Engine.UITransform.Get(entity).(*ec.UITransform)
		if text.ID == "level" {
			text.Text = fmt.Sprintf("LEVEL %d/%d", gameResources.Level.CurrentNum+1, gameResources.LevelCount)
			text.OffsetY = 0
			if gameResources.Level.Modified {
				if levelText == nil {
					levelText = &levelTextData{Texts: []string{"LEVEL 1/1", "LEVEL 1/1(*)"}, FontFace: text.FontFace, Pivot: uiTransform.Pivot}
					levelText.computeDotOffsets()
				}
				text.Text += "(*)"
				// Readjust text position
				text.OffsetY = levelText.YOffsets[0] - levelText.YOffsets[1]
			}
		}
	}))

	world.Manager.Join(world.Components.Engine.Text, world.Components.Engine.UITransform).Visit(ecs.Visit(func(entity ecs.Entity) {
		text := world.Components.Engine.Text.Get(entity).(*ec.Text)
		if text.ID == "box" {
			text.Text = fmt.Sprintf("BOX: %d/%d", boxOnGoalCount, boxCount)
		}
	}))

	world.Manager.Join(world.Components.Engine.Text, world.Components.Engine.UITransform).Visit(ecs.Visit(func(entity ecs.Entity) {
		text := world.Components.Engine.Text.Get(entity).(*ec.Text)
		if text.ID == "step" {
			text.Text = fmt.Sprintf("STEPS: %d", len(gameResources.Level.Movements))
		}
	}))

	// Finish level if all boxes are on goals
	if boxOnGoalCount == boxCount {
		gameResources.StateEvent = resources.StateEventLevelComplete
	}
}
