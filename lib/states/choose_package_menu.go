package states

import (
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/x-hgg-x/sokoban-go/lib/math"
	"github.com/x-hgg-x/sokoban-go/lib/resources"

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

// ChoosePackageState is the choose package state
type ChoosePackageState struct {
	packageNames     []string
	packageSelection int
	currentSelection int
	packageText      []*ec.Text
	arrowUpText      *ec.Text
	arrowDownText    *ec.Text
	exitTransition   states.Transition
}

//
// State interface
//

// OnPause method
func (st *ChoosePackageState) OnPause(world w.World) {}

// OnResume method
func (st *ChoosePackageState) OnResume(world w.World) {}

// OnStart method
func (st *ChoosePackageState) OnStart(world w.World) {
	prefabs := world.Resources.Prefabs.(*resources.Prefabs)
	loader.AddEntities(world, prefabs.Menu.ChoosePackageMenu)

	for _, file := range utils.Try(os.ReadDir("levels")) {
		fileName := file.Name()
		if filepath.Ext(fileName) == ".xsb" && len(fileName) > 4 {
			st.packageNames = append(st.packageNames, fileName[:len(fileName)-4])
		}
	}

	packageMap := map[string]int{}
	for index, name := range st.packageNames {
		packageMap[name] = index
	}

	if len(st.packageNames) == 0 {
		utils.LogFatalf("empty package list")
	}

	// Load last used package
	packageInfo := struct{ PackageName string }{"XSokoban"}
	toml.DecodeFile("config/package.toml", &packageInfo)

	if selection, ok := packageMap[packageInfo.PackageName]; ok {
		st.packageSelection = selection
	} else {
		st.packageSelection = packageMap["XSokoban"]
	}

	st.currentSelection = st.packageSelection

	// Find text components
	world.Manager.Join(world.Components.Engine.Text, world.Components.Engine.UITransform).Visit(ecs.Visit(func(entity ecs.Entity) {
		text := world.Components.Engine.Text.Get(entity).(*ec.Text)

		if strings.HasPrefix(text.ID, "package") {
			st.packageText = append(st.packageText, text)
		} else if text.ID == "arrow_up" {
			st.arrowUpText = text
		} else if text.ID == "arrow_down" {
			st.arrowDownText = text
		}
	}))

	sort.Slice(st.packageText, func(i, j int) bool { return st.packageText[i].ID < st.packageText[j].ID })
}

// OnStop method
func (st *ChoosePackageState) OnStop(world w.World) {
	// Save selected package
	var encoded strings.Builder
	encoder := toml.NewEncoder(&encoded)
	encoder.Indent = ""
	utils.LogError(encoder.Encode(struct{ PackageName string }{st.packageNames[st.packageSelection]}))
	utils.LogError(os.WriteFile("config/package.toml", []byte(encoded.String()), 0o666))

	world.Manager.DeleteAllEntities()
}

// Update method
func (st *ChoosePackageState) Update(world w.World) states.Transition {
	// Process inputs
	_, mouseWheelY := ebiten.Wheel()

	switch {
	case inpututil.IsKeyJustPressed(ebiten.KeyDown) || mouseWheelY < 0:
		st.currentSelection = math.Min(st.currentSelection+1, len(st.packageNames)-1)
	case inpututil.IsKeyJustPressed(ebiten.KeyUp) || mouseWheelY > 0:
		st.currentSelection = math.Max(st.currentSelection-1, 0)
	case inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace):
		st.packageSelection = st.currentSelection
		return st.exitTransition
	case inpututil.IsKeyJustPressed(ebiten.KeyEscape):
		return st.exitTransition
	}

	// Set text entities
	for index := 0; index < 8; index++ {
		textSelection := st.packageText[index]
		packageIndex := st.currentSelection - 3 + index

		if packageIndex == st.packageSelection {
			textSelection.Color = color.RGBA{R: 255}
		} else {
			textSelection.Color = color.RGBA{R: 255, G: 255, B: 255}
		}

		if 0 <= packageIndex && packageIndex < len(st.packageNames) {
			textSelection.Text = st.packageNames[packageIndex]
			textSelection.Color.A = 255
		} else {
			textSelection.Color.A = 0
		}
	}

	st.packageText[3].Text = fmt.Sprintf("\u25b6 %s", st.packageText[3].Text)

	switch st.currentSelection {
	case 0:
		st.arrowUpText.Color.A = 0
		st.arrowDownText.Color.A = 255
	case len(st.packageNames) - 1:
		st.arrowUpText.Color.A = 255
		st.arrowDownText.Color.A = 0
	default:
		st.arrowUpText.Color.A = 255
		st.arrowDownText.Color.A = 255
	}

	return states.Transition{}
}
