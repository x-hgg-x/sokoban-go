package main

import (
	_ "image/png"

	gc "github.com/x-hgg-x/sokoban-go/lib/components"
	gloader "github.com/x-hgg-x/sokoban-go/lib/loader"
	gr "github.com/x-hgg-x/sokoban-go/lib/resources"
	gs "github.com/x-hgg-x/sokoban-go/lib/states"

	"github.com/x-hgg-x/goecsengine/loader"
	er "github.com/x-hgg-x/goecsengine/resources"
	es "github.com/x-hgg-x/goecsengine/states"
	"github.com/x-hgg-x/goecsengine/utils"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	minGameWidth  = 960
	minGameHeight = 680
)

type mainGame struct {
	world        w.World
	stateMachine es.StateMachine
}

func (game *mainGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return minGameWidth, minGameHeight
}

func (game *mainGame) Update() error {
	game.stateMachine.Update(game.world)
	return nil
}

func (game *mainGame) Draw(screen *ebiten.Image) {
	game.stateMachine.Draw(game.world, screen)
}

func main() {
	world := w.InitWorld(&gc.Components{})

	// Init screen dimensions
	world.Resources.ScreenDimensions = &er.ScreenDimensions{Width: minGameWidth, Height: minGameHeight}

	// Load controls
	axes := []string{}
	actions := []string{
		gr.MoveUpAction, gr.MoveDownAction, gr.MoveLeftAction, gr.MoveRightAction,
		gr.MoveUpFastAction, gr.MoveDownFastAction, gr.MoveLeftFastAction, gr.MoveRightFastAction,
		gr.PreviousLevelAction, gr.PreviousLevelFastAction, gr.NextLevelAction, gr.NextLevelFastAction,
		gr.UndoAction, gr.UndoFastAction, gr.RestartAction, gr.SaveAction,
	}
	controls, inputHandler := loader.LoadControls("config/controls.toml", axes, actions)
	world.Resources.Controls = &controls
	world.Resources.InputHandler = &inputHandler

	// Load sprite sheets
	spriteSheets := loader.LoadSpriteSheets("assets/metadata/spritesheets/spritesheets.toml")
	world.Resources.SpriteSheets = &spriteSheets

	// Load fonts
	fonts := loader.LoadFonts("assets/metadata/fonts/fonts.toml")
	world.Resources.Fonts = &fonts

	// Load prefabs
	world.Resources.Prefabs = &gr.Prefabs{
		Menu: gr.MenuPrefabs{
			LevelCompleteMenu: gloader.PreloadEntities("assets/metadata/entities/ui/level_complete_menu.toml", world),
		},
		Game: gr.GamePrefabs{
			LevelInfo: gloader.PreloadEntities("assets/metadata/entities/ui/level.toml", world),
			BoxInfo:   gloader.PreloadEntities("assets/metadata/entities/ui/box.toml", world),
			StepInfo:  gloader.PreloadEntities("assets/metadata/entities/ui/step.toml", world),
		},
	}

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowSize(minGameWidth, minGameHeight)
	ebiten.SetWindowTitle("Sokoban")

	utils.LogError(ebiten.RunGame(&mainGame{world, es.Init(&gs.GameplayState{}, world)}))
}
