package states

import (
	"fmt"
	"image/color"
	"regexp"
	"strconv"

	"github.com/x-hgg-x/sokoban-go/lib/resources"
	g "github.com/x-hgg-x/sokoban-go/lib/systems"

	ecs "github.com/x-hgg-x/goecs/v2"
	ec "github.com/x-hgg-x/goecsengine/components"
	"github.com/x-hgg-x/goecsengine/states"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var regexpLevelForbiddenChars = regexp.MustCompile("[^0-9]")

// GoToLevelState is the main game state
type GoToLevelState struct {
	levelText  *ec.Text
	inputChars []rune
	maxLen     int
}

// OnStart method
func (st *GoToLevelState) OnStart(world w.World) {
	gameResources := world.Resources.Game.(*resources.Game)

	// Find level text component
	world.Manager.Join(world.Components.Engine.Text, world.Components.Engine.UITransform).Visit(ecs.Visit(func(entity ecs.Entity) {
		text := world.Components.Engine.Text.Get(entity).(*ec.Text)
		if text.ID == "level" {
			text.Color = color.RGBA{255, 0, 0, 255}
			st.levelText = text
			st.maxLen = len(strconv.Itoa(len(gameResources.Package.Levels)))
		}
	}))
}

// OnPause method
func (st *GoToLevelState) OnPause(world w.World) {}

// OnResume method
func (st *GoToLevelState) OnResume(world w.World) {}

// OnStop method
func (st *GoToLevelState) OnStop(world w.World) {
	st.levelText.Color = color.RGBA{255, 255, 255, 255}
}

// Update method
func (st *GoToLevelState) Update(world w.World) states.Transition {
	gameResources := world.Resources.Game.(*resources.Game)

	// Get user input
	st.inputChars = append(st.inputChars, []rune(regexpLevelForbiddenChars.ReplaceAllLiteralString(string(ebiten.AppendInputChars(nil)), ""))...)
	if len(st.inputChars) > st.maxLen {
		st.inputChars = st.inputChars[:st.maxLen]
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) && len(st.inputChars) > 0 {
		st.inputChars = st.inputChars[:len(st.inputChars)-1]
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		newLevel, err := strconv.Atoi(string(st.inputChars))
		if err == nil && 1 <= newLevel && newLevel <= len(gameResources.Package.Levels) {
			world.Manager.DeleteAllEntities()
			resources.SaveLevel(world)
			resources.InitLevel(world, newLevel-1)

			g.InfoSystem(world, false)
			g.GridUpdateSystem(world)
			g.GridTransformSystem(world)
		}
		return states.Transition{Type: states.TransPop}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return states.Transition{Type: states.TransPop}
	}

	// Update level text
	levelString := "_"
	if len(st.inputChars) > 0 {
		levelString = string(st.inputChars)
	}

	st.levelText.Text = fmt.Sprintf("LEVEL %s/%d", levelString, len(gameResources.Package.Levels))
	if gameResources.Level.Modified {
		st.levelText.Text += "(*)"
	}

	return states.Transition{}
}
