package resources

import (
	"fmt"

	gc "github.com/x-hgg-x/sokoban-go/lib/components"

	ecs "github.com/x-hgg-x/goecs/v2"
	ec "github.com/x-hgg-x/goecsengine/components"
	"github.com/x-hgg-x/goecsengine/loader"
	w "github.com/x-hgg-x/goecsengine/world"
)

const (
	// MaxWidth is the maximum level width
	MaxWidth = 30
	// MaxHeight is the maximum level height
	MaxHeight = 20
)

// StateEvent is an event for game progression
type StateEvent int

// List of game progression events
const (
	StateEventNone StateEvent = iota
	StateEventLevelComplete
)

// MovementType is a movement type
type MovementType int

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

var movementChars = []rune("udlrUDLR")

var movementCharMap = map[rune]MovementType{
	'u': MovementUp,
	'd': MovementDown,
	'l': MovementLeft,
	'r': MovementRight,
	'U': MovementUpPush,
	'D': MovementDownPush,
	'L': MovementLeftPush,
	'R': MovementRightPush,
}

// Tile contains tile entities
type Tile struct {
	Player *ecs.Entity
	Box    *ecs.Entity
	Goal   *ecs.Entity
	Wall   *ecs.Entity
}

// Level is a game level
type Level struct {
	CurrentNum int
	Grid       [MaxHeight][MaxWidth]Tile
	Movements  []MovementType
	Modified   bool
}

// Game contains game resources
type Game struct {
	StateEvent  StateEvent
	PackageName string
	LevelCount  int
	Level       Level
}

// NewGame creates a new game
func NewGame(world w.World, packageName string) *Game {
	prefabs := world.Resources.Prefabs.(*Prefabs)
	return &Game{PackageName: packageName, LevelCount: len(prefabs.Game.Levels)}
}

// InitLevel inits level
func InitLevel(world w.World, levelNum int) {
	gameComponents := world.Components.Game.(*gc.Components)
	gameResources := world.Resources.Game.(*Game)

	// Load ui entities
	prefabs := world.Resources.Prefabs.(*Prefabs)
	loader.AddEntities(world, prefabs.Game.BoxInfo)
	loader.AddEntities(world, prefabs.Game.StepInfo)
	levelInfoEntity := loader.AddEntities(world, prefabs.Game.LevelInfo)

	// Load level
	loader.AddEntities(world, prefabs.Game.Levels[levelNum])
	gameResources.Level = Level{CurrentNum: levelNum}

	// Set grid
	world.Manager.Join(gameComponents.Player, gameComponents.GridElement).Visit(ecs.Visit(func(entity ecs.Entity) {
		gridElement := gameComponents.GridElement.Get(entity).(*gc.GridElement)
		gameResources.Level.Grid[gridElement.Line][gridElement.Col].Player = &entity
	}))
	world.Manager.Join(gameComponents.Box, gameComponents.GridElement).Visit(ecs.Visit(func(entity ecs.Entity) {
		gridElement := gameComponents.GridElement.Get(entity).(*gc.GridElement)
		gameResources.Level.Grid[gridElement.Line][gridElement.Col].Box = &entity
	}))
	world.Manager.Join(gameComponents.Goal, gameComponents.GridElement).Visit(ecs.Visit(func(entity ecs.Entity) {
		gridElement := gameComponents.GridElement.Get(entity).(*gc.GridElement)
		gameResources.Level.Grid[gridElement.Line][gridElement.Col].Goal = &entity
	}))
	world.Manager.Join(gameComponents.Wall, gameComponents.GridElement).Visit(ecs.Visit(func(entity ecs.Entity) {
		gridElement := gameComponents.GridElement.Get(entity).(*gc.GridElement)
		gameResources.Level.Grid[gridElement.Line][gridElement.Col].Wall = &entity
	}))

	// Set level info text
	for iEntity := range levelInfoEntity {
		world.Components.Engine.Text.Get(levelInfoEntity[iEntity]).(*ec.Text).Text = fmt.Sprintf("LEVEL %d/%d", levelNum+1, gameResources.LevelCount)
	}

	LoadSave(world)
}
