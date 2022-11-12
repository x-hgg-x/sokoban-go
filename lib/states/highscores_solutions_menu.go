package states

import (
	"fmt"
	"image/color"
	"sort"
	"strconv"
	"strings"

	"github.com/x-hgg-x/sokoban-go/lib/math"
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

type highscore struct {
	levelNum  int
	author    string
	movements string
}

// HighscoresSolutionsState is the highscores and solutions state
type HighscoresSolutionsState struct {
	hasAuthor        bool
	highscores       []highscore
	currentSelection int
	scoreText        []*ec.Text
	arrowUpText      *ec.Text
	arrowDownText    *ec.Text
}

//
// State interface
//

// OnPause method
func (st *HighscoresSolutionsState) OnPause(world w.World) {}

// OnResume method
func (st *HighscoresSolutionsState) OnResume(world w.World) {}

// OnStart method
func (st *HighscoresSolutionsState) OnStart(world w.World) {
	packageName := resources.GetLastPackageName()
	prefabs := world.Resources.Prefabs.(*resources.Prefabs)

	var scoreIDPrefix string

	if st.hasAuthor {
		scoreIDPrefix = "score"
		loader.AddEntities(world, prefabs.Menu.HighscoresMenu)

		// Load highscores
		highscores := resources.HighscoreTable{}
		toml.DecodeFile(fmt.Sprintf("config/highscores/%s.toml", packageName), &highscores)
		resources.NormalizeHighScores(highscores)

		for k, v := range highscores {
			levelNum, err := strconv.Atoi(k[5:])
			if err == nil {
				st.highscores = append(st.highscores, highscore{levelNum: levelNum, author: v.Author, movements: v.Movements})
			}
		}
	} else {
		scoreIDPrefix = "steps"
		loader.AddEntities(world, prefabs.Menu.SolutionsMenu)

		// Load solutions
		solutions := map[string]string{}
		toml.DecodeFile(fmt.Sprintf("levels/solutions/%s.toml", packageName), &solutions)

		for k, v := range solutions {
			levelNum, err := strconv.Atoi(k[5:])
			if err == nil {
				st.highscores = append(st.highscores, highscore{levelNum: levelNum, movements: v})
			}
		}
	}

	packageInfoEntity := loader.AddEntities(world, prefabs.Game.PackageInfo)[0]
	world.Components.Engine.Text.Get(packageInfoEntity).(*ec.Text).Text = fmt.Sprintf("Package: %s", packageName)

	sort.Slice(st.highscores, func(i, j int) bool { return st.highscores[i].levelNum < st.highscores[j].levelNum })

	// Find text components
	world.Manager.Join(world.Components.Engine.Text, world.Components.Engine.UITransform).Visit(ecs.Visit(func(entity ecs.Entity) {
		text := world.Components.Engine.Text.Get(entity).(*ec.Text)

		if strings.HasPrefix(text.ID, scoreIDPrefix) {
			st.scoreText = append(st.scoreText, text)
		} else if text.ID == "cursor" {
			if len(st.highscores) > 0 {
				text.Color.A = 255
			}
		} else if text.ID == "arrow_up" {
			st.arrowUpText = text
		} else if text.ID == "arrow_down" {
			st.arrowDownText = text
		}
	}))

	sort.Slice(st.scoreText, func(i, j int) bool { return st.scoreText[i].ID < st.scoreText[j].ID })
}

// OnStop method
func (st *HighscoresSolutionsState) OnStop(world w.World) {
	world.Manager.DeleteAllEntities()
}

// Update method
func (st *HighscoresSolutionsState) Update(world w.World) states.Transition {
	if len(st.highscores) > 0 {
		// Process inputs
		_, mouseWheelY := ebiten.Wheel()

		switch {
		case inpututil.IsKeyJustPressed(ebiten.KeyDown) || mouseWheelY < 0:
			st.currentSelection = math.Min(st.currentSelection+1, len(st.highscores)-1)
		case inpututil.IsKeyJustPressed(ebiten.KeyUp) || mouseWheelY > 0:
			st.currentSelection = math.Max(st.currentSelection-1, 0)
		case inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft):
			currentScore := st.highscores[st.currentSelection]
			exitTransition := states.Transition{Type: states.TransSwitch, NewStates: []states.State{&HighscoresSolutionsState{hasAuthor: st.hasAuthor}}}
			return states.Transition{Type: states.TransSwitch, NewStates: []states.State{&ViewSolutionState{levelNum: currentScore.levelNum, hasAuthor: st.hasAuthor, exitTransition: exitTransition}}}
		}

		// Set text entities
		for index := 0; index < 10; index++ {
			textSelection := st.scoreText[index]
			scoreIndex := st.currentSelection - 4 + index

			if scoreIndex == st.currentSelection {
				textSelection.Color = color.RGBA{R: 255}
			} else {
				textSelection.Color = color.RGBA{R: 255, G: 255, B: 255}
			}

			if 0 <= scoreIndex && scoreIndex < len(st.highscores) {
				score := st.highscores[scoreIndex]

				if st.hasAuthor {
					textSelection.Text = fmt.Sprintf("%4d  %6s  %5d", score.levelNum, score.author, len(score.movements))
				} else {
					textSelection.Text = fmt.Sprintf("%4d     %5d", score.levelNum, len(score.movements))
				}

				textSelection.Color.A = 255
			} else {
				textSelection.Color.A = 0
			}
		}

		switch st.currentSelection {
		case 0:
			st.arrowUpText.Color.A = 0
			st.arrowDownText.Color.A = 255
		case len(st.highscores) - 1:
			st.arrowUpText.Color.A = 255
			st.arrowDownText.Color.A = 0
		default:
			st.arrowUpText.Color.A = 255
			st.arrowDownText.Color.A = 255
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return states.Transition{Type: states.TransSwitch, NewStates: []states.State{&MainMenuState{}}}
	}
	return states.Transition{}
}
