package states

import (
	"fmt"

	gloader "github.com/x-hgg-x/sokoban-go/lib/loader"
	"github.com/x-hgg-x/sokoban-go/lib/resources"
	g "github.com/x-hgg-x/sokoban-go/lib/systems"

	"github.com/x-hgg-x/goecsengine/states"
	"github.com/x-hgg-x/goecsengine/utils"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// ViewSolutionState is the view solution state
type ViewSolutionState struct {
	levelNum       int
	movements      []resources.MovementType
	invalidLevel   bool
	exitTransition states.Transition
}

// OnStart method
func (st *ViewSolutionState) OnStart(world w.World) {
	// Load game
	packageName := resources.GetLastPackageName()
	packageData := utils.Try(gloader.LoadPackage(packageName))

	if 1 <= st.levelNum && st.levelNum <= len(packageData.Levels) {
		world.Resources.Game = &resources.Game{Package: packageData}
		resources.InitLevel(world, st.levelNum-1)
	} else {
		st.invalidLevel = true
	}
}

// OnPause method
func (st *ViewSolutionState) OnPause(world w.World) {}

// OnResume method
func (st *ViewSolutionState) OnResume(world w.World) {}

// OnStop method
func (st *ViewSolutionState) OnStop(world w.World) {
	world.Manager.DeleteAllEntities()
	world.Resources.Game = nil
}

// Update method
func (st *ViewSolutionState) Update(world w.World) states.Transition {
	if st.invalidLevel {
		fmt.Printf("invalid level number: %d\n", st.levelNum)
		return st.exitTransition
	}

	g.MoveSolutionSystem(world, st.movements)
	g.InfoSystem(world, true)
	g.GridUpdateSystem(world)
	g.GridTransformSystem(world)

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return st.exitTransition
	}
	return states.Transition{}
}
