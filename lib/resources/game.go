package resources

import (
	"fmt"

	gc "github.com/x-hgg-x/sokoban-go/lib/components"

	ecs "github.com/x-hgg-x/goecs/v2"
	ec "github.com/x-hgg-x/goecsengine/components"
	"github.com/x-hgg-x/goecsengine/loader"
	w "github.com/x-hgg-x/goecsengine/world"
)

// Game contains game resources
type Game struct {
	CurrentLevel int
	LevelCount   int
	Steps        int
	Goals        map[gc.GridElement]struct{}
}

// InitLevel inits level
func InitLevel(world w.World, levelNum int) {
	// Load ui entities
	prefabs := world.Resources.Prefabs.(*Prefabs)
	loader.AddEntities(world, prefabs.Game.BoxInfo)
	loader.AddEntities(world, prefabs.Game.StepInfo)
	levelInfoEntity := loader.AddEntities(world, prefabs.Game.LevelInfo)

	// Load level
	loader.AddEntities(world, prefabs.Game.Levels[levelNum])
	game := &Game{
		CurrentLevel: levelNum,
		LevelCount:   len(prefabs.Game.Levels),
		Goals:        make(map[gc.GridElement]struct{}),
	}

	// Set grid position of goals
	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(gameComponents.Goal, gameComponents.GridElement).Visit(ecs.Visit(func(entity ecs.Entity) {
		goalGrid := gameComponents.GridElement.Get(entity).(*gc.GridElement)
		game.Goals[*goalGrid] = struct{}{}
	}))

	// Set level info text
	for iEntity := range levelInfoEntity {
		world.Components.Engine.Text.Get(levelInfoEntity[iEntity]).(*ec.Text).Text = fmt.Sprintf("LEVEL %d/%d", levelNum+1, game.LevelCount)
	}

	world.Resources.Game = game
}
