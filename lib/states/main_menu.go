package states

import (
	"fmt"

	"github.com/x-hgg-x/sokoban-go/lib/resources"

	"github.com/x-hgg-x/goecsengine/loader"
	"github.com/x-hgg-x/goecsengine/states"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// MainMenuState is the main menu state
type MainMenuState struct {
	selection int
}

//
// Menu interface
//

func (st *MainMenuState) getSelection() int {
	return st.selection
}

func (st *MainMenuState) setSelection(selection int) {
	st.selection = selection
}

func (st *MainMenuState) confirmSelection() states.Transition {
	switch st.selection {
	case 0:
		// Start
		return states.Transition{Type: states.TransSwitch, NewStates: []states.State{&GameplayState{}}}
	case 1:
		// View highscores
		return states.Transition{Type: states.TransNone}
	case 2:
		// View solutions
		return states.Transition{Type: states.TransNone}
	case 3:
		// Choose package
		return states.Transition{Type: states.TransSwitch, NewStates: []states.State{&ChoosePackageState{
			exitTransition: states.Transition{Type: states.TransSwitch, NewStates: []states.State{&MainMenuState{}}},
		}}}
	case 4:
		// Exit
		return states.Transition{Type: states.TransQuit}
	}
	panic(fmt.Errorf("unknown selection: %d", st.selection))
}

func (st *MainMenuState) getMenuIDs() []string {
	return []string{"start", "view_highscores", "view_solutions", "choose_package", "exit"}
}

func (st *MainMenuState) getCursorMenuIDs() []string {
	return []string{"cursor_start", "cursor_view_highscores", "cursor_view_solutions", "cursor_choose_package", "cursor_exit"}
}

//
// State interface
//

// OnPause method
func (st *MainMenuState) OnPause(world w.World) {}

// OnResume method
func (st *MainMenuState) OnResume(world w.World) {}

// OnStart method
func (st *MainMenuState) OnStart(world w.World) {
	prefabs := world.Resources.Prefabs.(*resources.Prefabs)
	loader.AddEntities(world, prefabs.Menu.MainMenu)
}

// OnStop method
func (st *MainMenuState) OnStop(world w.World) {
	world.Manager.DeleteAllEntities()
}

// Update method
func (st *MainMenuState) Update(world w.World) states.Transition {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return states.Transition{Type: states.TransQuit}
	}
	return updateMenu(st, world)
}
