package states

import (
	"github.com/x-hgg-x/sokoban-go/lib/resources"

	"github.com/x-hgg-x/goecsengine/loader"
	"github.com/x-hgg-x/goecsengine/states"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// LevelCompleteState is the level complete menu state
type LevelCompleteState struct{}

// OnStart method
func (st *LevelCompleteState) OnStart(world w.World) {
	prefabs := world.Resources.Prefabs.(*resources.Prefabs)
	loader.AddEntities(world, prefabs.Menu.LevelCompleteMenu)
}

// OnPause method
func (st *LevelCompleteState) OnPause(world w.World) {}

// OnResume method
func (st *LevelCompleteState) OnResume(world w.World) {}

// OnStop method
func (st *LevelCompleteState) OnStop(world w.World) {
	gameResources := world.Resources.Game.(*resources.Game)

	world.Manager.DeleteAllEntities()
	gameResources.Level.Movements = []resources.MovementType{}
	gameResources.Level.Modified = true
	resources.SaveLevel(world)
	resources.InitLevel(world, gameResources.Level.CurrentNum)
}

// Update method
func (st *LevelCompleteState) Update(world w.World) states.Transition {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		return states.Transition{Type: states.TransPop}
	}
	return states.Transition{}
}
