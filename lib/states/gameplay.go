package states

import (
	"fmt"
	"os"

	gloader "github.com/x-hgg-x/sokoban-go/lib/loader"
	"github.com/x-hgg-x/sokoban-go/lib/resources"
	g "github.com/x-hgg-x/sokoban-go/lib/systems"

	"github.com/x-hgg-x/goecsengine/states"
	"github.com/x-hgg-x/goecsengine/utils"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/BurntSushi/toml"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// GameplayState is the main game state
type GameplayState struct{}

// OnStart method
func (st *GameplayState) OnStart(world w.World) {
	// Load game
	packageName := resources.GetLastPackageName()
	packageData := utils.Try(gloader.LoadPackage(packageName))

	// Load save configuration
	var saveConfig resources.SaveConfig

	utils.LogError(os.MkdirAll("config/saves", os.ModePerm))
	if saveFile, err := os.ReadFile(fmt.Sprintf("config/saves/%s.toml", packageName)); err == nil {
		var encodedSaveConfig resources.EncodedSaveConfig
		utils.Try(toml.Decode(string(saveFile), &encodedSaveConfig))
		saveConfig = utils.Try(encodedSaveConfig.Decode())
	} else {
		saveConfig = resources.EmptySaveConfig()
	}

	// Load last played level
	levelNum := 0
	if 1 <= saveConfig.CurrentLevel && saveConfig.CurrentLevel <= len(packageData.Levels) {
		levelNum = saveConfig.CurrentLevel - 1
	}

	world.Resources.Game = &resources.Game{Package: packageData, SaveConfig: saveConfig}
	resources.InitLevel(world, levelNum)
}

// OnPause method
func (st *GameplayState) OnPause(world w.World) {}

// OnResume method
func (st *GameplayState) OnResume(world w.World) {}

// OnStop method
func (st *GameplayState) OnStop(world w.World) {
	world.Manager.DeleteAllEntities()
	resources.SaveLevel(world)
	world.Resources.Game = nil
}

// Update method
func (st *GameplayState) Update(world w.World) states.Transition {
	g.SwitchLevelSystem(world)
	g.UndoSystem(world)
	g.MoveSystem(world)
	g.SaveSystem(world)
	g.InfoSystem(world, false)
	g.GridUpdateSystem(world)
	g.GridTransformSystem(world)

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return states.Transition{Type: states.TransPush, NewStates: []states.State{&PauseMenuState{}}}
	}

	gameResources := world.Resources.Game.(*resources.Game)
	switch gameResources.StateEvent {
	case resources.StateEventLevelComplete:
		gameResources.StateEvent = resources.StateEventNone
		return states.Transition{Type: states.TransPush, NewStates: []states.State{&LevelCompleteState{}}}
	}

	return states.Transition{}
}
