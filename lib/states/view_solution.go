package states

import (
	"fmt"
	"os"
	"strings"

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

// ViewSolutionState is the view solution state
type ViewSolutionState struct {
	levelNum       int
	hasAuthor      bool
	invalidLevel   bool
	packageName    string
	data           interface{}
	movements      []resources.MovementType
	exitTransition states.Transition
}

// OnStart method
func (st *ViewSolutionState) OnStart(world w.World) {
	// Load game
	st.packageName = resources.GetLastPackageName()
	packageData := utils.Try(gloader.LoadPackage(st.packageName))

	if 1 <= st.levelNum && st.levelNum <= len(packageData.Levels) {
		world.Resources.Game = &resources.Game{Package: packageData}
		resources.InitLevel(world, st.levelNum-1)
	} else {
		st.invalidLevel = true
	}

	// Load solution
	if st.hasAuthor {
		highscores := resources.HighscoreTable{}
		toml.DecodeFile(fmt.Sprintf("config/highscores/%s.toml", st.packageName), &highscores)
		resources.NormalizeHighScores(highscores)
		st.data = highscores
		st.movements = resources.DecodeMovements(highscores[fmt.Sprintf("Level%04d", st.levelNum)].Movements)

	} else {
		solutions := map[string]string{}
		toml.DecodeFile(fmt.Sprintf("levels/solutions/%s.toml", st.packageName), &solutions)
		st.data = solutions
		st.movements = resources.DecodeMovements(solutions[fmt.Sprintf("Level%04d", st.levelNum)])
	}

	if !st.checkSolution(world) {
		st.invalidLevel = true
	}
}

// OnPause method
func (st *ViewSolutionState) OnPause(world w.World) {}

// OnResume method
func (st *ViewSolutionState) OnResume(world w.World) {}

// OnStop method
func (st *ViewSolutionState) OnStop(world w.World) {
	// Save normalized solution
	var filename string

	switch data := st.data.(type) {
	case resources.HighscoreTable:
		filename = fmt.Sprintf("config/highscores/%s.toml", st.packageName)
		if !st.invalidLevel {
			key := fmt.Sprintf("Level%04d", st.levelNum)
			highscore := data[key]
			highscore.Movements = resources.EncodeMovements(st.movements)
			data[key] = highscore
		} else {
			delete(data, fmt.Sprintf("Level%04d", st.levelNum))
		}
	case map[string]string:
		filename = fmt.Sprintf("levels/solutions/%s.toml", st.packageName)
		if !st.invalidLevel {
			data[fmt.Sprintf("Level%04d", st.levelNum)] = resources.EncodeMovements(st.movements)
		} else {
			delete(data, fmt.Sprintf("Level%04d", st.levelNum))
		}
	}

	var encoded strings.Builder
	encoder := toml.NewEncoder(&encoded)
	encoder.Indent = ""
	utils.LogError(encoder.Encode(st.data))
	utils.LogError(os.WriteFile(filename, []byte(encoded.String()), 0o666))

	world.Manager.DeleteAllEntities()
	world.Resources.Game = nil
}

// Update method
func (st *ViewSolutionState) Update(world w.World) states.Transition {
	if st.invalidLevel {
		fmt.Printf("invalid solution for level %d\n", st.levelNum)
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

func (st *ViewSolutionState) checkSolution(world w.World) bool {
	if st.invalidLevel {
		return false
	}

	gameResources := world.Resources.Game.(*resources.Game)
	resources.Move(world, st.movements...)

	for _, tile := range gameResources.Level.Grid.Data {
		hasBox := tile.Contains(resources.TileBox)
		hasGoal := tile.Contains(resources.TileGoal)

		if hasBox != hasGoal {
			return false
		}
	}

	st.movements = gameResources.Level.Movements

	for len(gameResources.Level.Movements) > 0 {
		resources.Undo(world)
	}

	return true
}
