package states

import (
	"fmt"
	"image/color"

	"github.com/x-hgg-x/sokoban-go/lib/resources"

	ecs "github.com/x-hgg-x/goecs/v2"
	ec "github.com/x-hgg-x/goecsengine/components"
	"github.com/x-hgg-x/goecsengine/loader"
	"github.com/x-hgg-x/goecsengine/states"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/BurntSushi/toml"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// PauseMenuState is the pause menu state
type PauseMenuState struct {
	pauseMenu        []ecs.Entity
	selection        int
	invalidHighscore bool
	invalidSolution  bool
}

//
// Menu interface
//

func (st *PauseMenuState) getSelection() int {
	return st.selection
}

func (st *PauseMenuState) setSelection(selection int) {
	st.selection = selection
}

func (st *PauseMenuState) confirmSelection(world w.World) states.Transition {
	switch st.selection {
	case 0:
		// Resume
		return states.Transition{Type: states.TransPop}
	case 1:
		// Retry
		gameResources := world.Resources.Game.(*resources.Game)
		gameResources.Level.Movements = []resources.MovementType{}
		gameResources.Level.Modified = true
		return states.Transition{Type: states.TransReplace, NewStates: []states.State{&GameplayState{}}}
	case 2:
		// View highscore
		if !st.invalidHighscore {
			world.Resources.Game.(*resources.Game).Level.Modified = true
			levelNum := world.Resources.Game.(*resources.Game).Level.CurrentNum + 1
			exitTransition := states.Transition{Type: states.TransSwitch, NewStates: []states.State{&GameplayState{}}}
			return states.Transition{Type: states.TransReplace, NewStates: []states.State{&ViewSolutionState{levelNum: levelNum, hasAuthor: true, exitTransition: exitTransition}}}
		} else {
			return states.Transition{}
		}
	case 3:
		// View solution
		if !st.invalidSolution {
			world.Resources.Game.(*resources.Game).Level.Modified = true
			levelNum := world.Resources.Game.(*resources.Game).Level.CurrentNum + 1
			exitTransition := states.Transition{Type: states.TransSwitch, NewStates: []states.State{&GameplayState{}}}
			return states.Transition{Type: states.TransReplace, NewStates: []states.State{&ViewSolutionState{levelNum: levelNum, hasAuthor: false, exitTransition: exitTransition}}}
		} else {
			return states.Transition{}
		}
	case 4:
		// Choose package
		world.Resources.Game.(*resources.Game).Level.Modified = true
		return states.Transition{Type: states.TransReplace, NewStates: []states.State{&ChoosePackageState{
			exitTransition: states.Transition{Type: states.TransSwitch, NewStates: []states.State{&GameplayState{}}},
		}}}
	case 5:
		// Main Menu
		return states.Transition{Type: states.TransReplace, NewStates: []states.State{&MainMenuState{}}}
	}
	panic(fmt.Errorf("unknown selection: %d", st.selection))
}

func (st *PauseMenuState) getMenuIDs() []string {
	return []string{"resume", "retry", "view_highscore", "view_solution", "choose_package", "main_menu"}
}

func (st *PauseMenuState) getCursorMenuIDs() []string {
	return []string{"cursor_resume", "cursor_retry", "cursor_view_highscore", "cursor_view_solution", "cursor_choose_package", "cursor_main_menu"}
}

//
// State interface
//

// OnPause method
func (st *PauseMenuState) OnPause(world w.World) {}

// OnResume method
func (st *PauseMenuState) OnResume(world w.World) {}

// OnStart method
func (st *PauseMenuState) OnStart(world w.World) {
	gameResources := world.Resources.Game.(*resources.Game)

	prefabs := world.Resources.Prefabs.(*resources.Prefabs)
	st.pauseMenu = append(st.pauseMenu, loader.AddEntities(world, prefabs.Menu.PauseMenu)...)

	key := fmt.Sprintf("Level%04d", gameResources.Level.CurrentNum+1)

	highscores := resources.HighscoreTable{}
	toml.DecodeFile(fmt.Sprintf("config/highscores/%s.toml", gameResources.Package.Name), &highscores)
	if _, ok := highscores[key]; !ok {
		st.invalidHighscore = true
	}

	solutions := map[string]string{}
	toml.DecodeFile(fmt.Sprintf("levels/solutions/%s.toml", gameResources.Package.Name), &solutions)
	if _, ok := solutions[key]; !ok {
		st.invalidSolution = true
	}

	// Update text components
	world.Manager.Join(world.Components.Engine.Text, world.Components.Engine.UITransform).Visit(ecs.Visit(func(entity ecs.Entity) {
		text := world.Components.Engine.Text.Get(entity).(*ec.Text)

		switch text.ID {
		case "view_highscore":
			if st.invalidHighscore {
				text.Color = color.RGBA{0, 0, 0, 120}
			}
		case "view_solution":
			if st.invalidSolution {
				text.Color = color.RGBA{0, 0, 0, 120}
			}
		}
	}))
}

// OnStop method
func (st *PauseMenuState) OnStop(world w.World) {
	world.Manager.DeleteEntities(st.pauseMenu...)
}

// Update method
func (st *PauseMenuState) Update(world w.World) states.Transition {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return states.Transition{Type: states.TransPop}
	}
	return updateMenu(st, world)
}
