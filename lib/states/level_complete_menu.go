package states

import (
	"fmt"
	"image/color"
	"os"
	"strings"

	"github.com/x-hgg-x/sokoban-go/lib/resources"
	g "github.com/x-hgg-x/sokoban-go/lib/systems"

	ecs "github.com/x-hgg-x/goecs/v2"
	ec "github.com/x-hgg-x/goecsengine/components"
	"github.com/x-hgg-x/goecsengine/loader"
	"github.com/x-hgg-x/goecsengine/states"
	"github.com/x-hgg-x/goecsengine/utils"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/BurntSushi/toml"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// LevelCompleteState is the level complete menu state
type LevelCompleteState struct {
	highscores   resources.HighscoreTable
	newHighscore *resources.Highscore
	text1        *ec.Text
	text2        *ec.Text
}

// OnStart method
func (st *LevelCompleteState) OnStart(world w.World) {
	gameResources := world.Resources.Game.(*resources.Game)

	prefabs := world.Resources.Prefabs.(*resources.Prefabs)
	loader.AddEntities(world, prefabs.Menu.LevelCompleteMenu)

	// Find text components
	world.Manager.Join(world.Components.Engine.Text, world.Components.Engine.UITransform).Visit(ecs.Visit(func(entity ecs.Entity) {
		text := world.Components.Engine.Text.Get(entity).(*ec.Text)

		switch text.ID {
		case "text1":
			st.text1 = text
		case "text2":
			st.text2 = text
		}
	}))

	// Load highscores
	st.highscores = resources.HighscoreTable{}
	toml.DecodeFile(fmt.Sprintf("config/highscores/%s.toml", gameResources.Package.Name), &st.highscores)
	resources.NormalizeHighScores(st.highscores)

	currentHighscore, ok := st.highscores[fmt.Sprintf("Level%04d", gameResources.Level.CurrentNum+1)]
	if !ok || len(gameResources.Level.Movements) < len(currentHighscore.Movements) {
		st.newHighscore = &resources.Highscore{Movements: resources.EncodeMovements(gameResources.Level.Movements)}
		st.text1.Text = "NEW RECORD !"
		st.text2.Text = "NAME: ______"
		st.text2.Color = color.RGBA{R: 255, A: 255}
	} else {
		st.text1.Text = fmt.Sprintf("BEST SCORE: %s", currentHighscore.Author)
		st.text2.Text = fmt.Sprintf("%d STEPS", len(currentHighscore.Movements))
	}

	// Reset level movements
	gameResources.Level.Movements = []resources.MovementType{}
	gameResources.Level.Modified = true
}

// OnPause method
func (st *LevelCompleteState) OnPause(world w.World) {}

// OnResume method
func (st *LevelCompleteState) OnResume(world w.World) {}

// OnStop method
func (st *LevelCompleteState) OnStop(world w.World) {}

// Update method
func (st *LevelCompleteState) Update(world w.World) states.Transition {
	gameResources := world.Resources.Game.(*resources.Game)

	if st.newHighscore != nil {
		// Set highscore author
		// Get user input
		st.newHighscore.Author += strings.ToUpper(resources.RegexpHighscoreForbiddenChars.ReplaceAllLiteralString(string(ebiten.AppendInputChars(nil)), ""))
		if len(st.newHighscore.Author) > resources.MaxAuthorLen {
			st.newHighscore.Author = st.newHighscore.Author[:resources.MaxAuthorLen]
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) && len(st.newHighscore.Author) > 0 {
			st.newHighscore.Author = st.newHighscore.Author[:len(st.newHighscore.Author)-1]
		}

		// Set new highscore text
		padding := strings.Repeat("_", resources.MaxAuthorLen-len(st.newHighscore.Author))
		st.text2.Text = fmt.Sprintf("NAME: %s%s", st.newHighscore.Author, padding)

		// Validate highscore
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			st.text1.Text = fmt.Sprintf("BEST SCORE: %s", st.newHighscore.Author)
			st.text2.Text = fmt.Sprintf("%d STEPS", len(st.newHighscore.Movements))
			st.text2.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}

			st.highscores[fmt.Sprintf("Level%04d", gameResources.Level.CurrentNum+1)] = *st.newHighscore
			st.newHighscore = nil

			// Save highscores
			var encoded strings.Builder
			encoder := toml.NewEncoder(&encoded)
			encoder.Indent = ""
			utils.LogError(encoder.Encode(st.highscores))
			utils.LogError(os.WriteFile(fmt.Sprintf("config/highscores/%s.toml", gameResources.Package.Name), []byte(encoded.String()), 0o666))
		}
	} else {
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			return states.Transition{Type: states.TransReplace, NewStates: []states.State{&MainMenuState{}}}
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			world.Resources.InputHandler.Actions[resources.RestartAction] = true
		}

		if g.SwitchLevelSystem(world) {
			return states.Transition{Type: states.TransPop}
		}
	}
	return states.Transition{}
}
