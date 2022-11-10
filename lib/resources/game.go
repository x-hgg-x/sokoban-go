package resources

import (
	"fmt"

	gloader "github.com/x-hgg-x/sokoban-go/lib/loader"
	"github.com/x-hgg-x/sokoban-go/lib/math"
	gutils "github.com/x-hgg-x/sokoban-go/lib/utils"

	ec "github.com/x-hgg-x/goecsengine/components"
	"github.com/x-hgg-x/goecsengine/loader"
	"github.com/x-hgg-x/goecsengine/utils"
	w "github.com/x-hgg-x/goecsengine/world"
)

const (
	offsetX       = 0
	offsetY       = 80
	gridBlockSize = 32

	minGridWidth  = 30
	minGridHeight = 20
)

// StateEvent is an event for game progression
type StateEvent int

// List of game progression events
const (
	StateEventNone StateEvent = iota
	StateEventLevelComplete
)

// MovementType is a movement type
type MovementType uint8

// List of movements
const (
	MovementUp MovementType = iota
	MovementDown
	MovementLeft
	MovementRight
	MovementUpPush
	MovementDownPush
	MovementLeftPush
	MovementRightPush
)

// GetSimpleMovement returns simple movement
func GetSimpleMovement(m MovementType) MovementType {
	return m % 4
}

// GetPushMovement returns push movement
func GetPushMovement(m MovementType) MovementType {
	return m%4 + 4
}

// MovementChars contains the char representation of each movement type
const MovementChars = "udlrUDLR"

// MovementCharMap contains the movement type associated to a char representation
var MovementCharMap = map[byte]MovementType{
	'u': MovementUp,
	'd': MovementDown,
	'l': MovementLeft,
	'r': MovementRight,
	'U': MovementUpPush,
	'D': MovementDownPush,
	'L': MovementLeftPush,
	'R': MovementRightPush,
}

// PackageData contains level package data
type PackageData = gloader.PackageData

// Tile is a game tile
type Tile = gloader.Tile

// List of game tiles
const (
	TilePlayer = gloader.TilePlayer
	TileBox    = gloader.TileBox
	TileGoal   = gloader.TileGoal
	TileWall   = gloader.TileWall
	TileEmpty  = gloader.TileEmpty
)

// Level is a game level
type Level struct {
	CurrentNum int
	Grid       gutils.Vec2d[Tile]
	Movements  []MovementType
	Modified   bool
}

// GridLayout is the grid layout
type GridLayout struct {
	Width  int
	Height int
}

// Game contains game resources
type Game struct {
	StateEvent StateEvent
	Package    PackageData
	Level      Level
	GridLayout GridLayout
	SaveConfig SaveConfig
}

// InitLevel inits level
func InitLevel(world w.World, levelNum int) {
	gameResources := world.Resources.Game.(*Game)

	// Load ui entities
	prefabs := world.Resources.Prefabs.(*Prefabs)
	loader.AddEntities(world, prefabs.Game.BoxInfo)
	loader.AddEntities(world, prefabs.Game.StepInfo)
	loader.AddEntities(world, prefabs.Game.PackageInfo)
	levelInfoEntity := loader.AddEntities(world, prefabs.Game.LevelInfo)[0]

	// Load level
	level := gameResources.Package.Levels[levelNum]
	gridLayout := &gameResources.GridLayout
	gridLayout.Width = math.Max(minGridWidth, level.NCols)
	gridLayout.Height = math.Max(minGridHeight, level.NRows)

	UpdateGameLayout(world, gridLayout)

	gameSpriteSheet := (*world.Resources.SpriteSheets)["game"]
	grid, levelComponentList := utils.Try2(gloader.LoadLevel(gameResources.Package, levelNum, gridLayout.Width, gridLayout.Height, &gameSpriteSheet))
	loader.AddEntities(world, levelComponentList)
	gameResources.Level = Level{CurrentNum: levelNum, Grid: grid}

	// Set level info text
	world.Components.Engine.Text.Get(levelInfoEntity).(*ec.Text).Text = fmt.Sprintf("LEVEL %d/%d", levelNum+1, len(gameResources.Package.Levels))

	LoadSave(world)
}

// UpdateGameLayout updates game layout
func UpdateGameLayout(world w.World, gridLayout *GridLayout) (int, int) {
	gridWidth, gridHeight := minGridWidth, minGridHeight

	if gridLayout != nil {
		gridWidth = gridLayout.Width
		gridHeight = gridLayout.Height
	}

	gameWidth := gridWidth*gridBlockSize + offsetX
	gameHeight := gridHeight*gridBlockSize + offsetY

	fadeOutSprite := &(*world.Resources.SpriteSheets)["background"].Sprites[0]
	fadeOutSprite.Width = gameWidth
	fadeOutSprite.Height = gameHeight

	world.Resources.ScreenDimensions.Width = gameWidth
	world.Resources.ScreenDimensions.Height = gameHeight

	return gameWidth, gameHeight
}
