package states

import (
	gloader "github.com/x-hgg-x/sokoban-go/lib/loader"
	"github.com/x-hgg-x/sokoban-go/lib/resources"
	g "github.com/x-hgg-x/sokoban-go/lib/systems"

	"github.com/x-hgg-x/goecsengine/states"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

// GameplayState is the main game state
type GameplayState struct{}

// OnStart method
func (st *GameplayState) OnStart(world w.World) {
	// Load package
	gloader.LoadPackage("levels/xsokoban/levels.txt", world)

	// Load first level
	resources.InitLevel(world, 0)
}

// OnPause method
func (st *GameplayState) OnPause(world w.World) {}

// OnResume method
func (st *GameplayState) OnResume(world w.World) {}

// OnStop method
func (st *GameplayState) OnStop(world w.World) {
	world.Resources.Game = nil
	world.Manager.DeleteAllEntities()
}

// Update method
func (st *GameplayState) Update(world w.World, screen *ebiten.Image) states.Transition {
	g.SwitchLevelSystem(world)
	g.UndoSystem(world)
	g.MoveSystem(world)
	g.InfoSystem(world)
	g.GridTransformSystem(world)

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return states.Transition{Type: states.TransQuit}
	}

	gameResources := world.Resources.Game.(*resources.Game)
	switch gameResources.StateEvent {
	case resources.StateEventLevelComplete:
		gameResources.StateEvent = resources.StateEventNone
		return states.Transition{Type: states.TransPush, NewStates: []states.State{&LevelCompleteState{}}}
	}

	return states.Transition{}
}
