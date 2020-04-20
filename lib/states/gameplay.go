package states

import (
	gloader "github.com/x-hgg-x/sokoban-go/lib/loader"
	"github.com/x-hgg-x/sokoban-go/lib/resources"
	g "github.com/x-hgg-x/sokoban-go/lib/systems"

	"github.com/x-hgg-x/goecsengine/loader"
	"github.com/x-hgg-x/goecsengine/states"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

// GameplayState is the main game state
type GameplayState struct{}

// OnStart method
func (st *GameplayState) OnStart(world w.World) {
	// Load level
	gloader.LoadPackage("levels/xsokoban/levels.txt", world)

	// Load ui entities
	prefabs := world.Resources.Prefabs.(*resources.Prefabs)
	loader.AddEntities(world, prefabs.Game.LevelInfo)
	loader.AddEntities(world, prefabs.Game.BoxInfo)
	loader.AddEntities(world, prefabs.Game.StepInfo)
}

// OnPause method
func (st *GameplayState) OnPause(world w.World) {}

// OnResume method
func (st *GameplayState) OnResume(world w.World) {}

// OnStop method
func (st *GameplayState) OnStop(world w.World) {}

// Update method
func (st *GameplayState) Update(world w.World, screen *ebiten.Image) states.Transition {
	g.SetTransformSystem(world)

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return states.Transition{Type: states.TransQuit}
	}
	return states.Transition{}
}
