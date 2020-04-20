package loader

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	gc "github.com/x-hgg-x/sokoban-go/lib/components"
	"github.com/x-hgg-x/sokoban-go/lib/resources"

	ecs "github.com/x-hgg-x/goecs/v2"
	ec "github.com/x-hgg-x/goecsengine/components"
	"github.com/x-hgg-x/goecsengine/math"
	"github.com/x-hgg-x/goecsengine/utils"
	w "github.com/x-hgg-x/goecsengine/world"
)

const (
	maxWidth  = 30
	maxHeight = 20
)

const (
	charFloor1       = ' '
	charFloor2       = '-'
	charFloor3       = '_'
	charFloor        = ' '
	charExterior     = '_'
	charWall         = '#'
	charGoal         = '.'
	charBox          = '$'
	charPlayer       = '@'
	charBoxOnGoal    = '*'
	charPlayerOnGoal = '+'
)

// LoadPackage loads level package from a text file
func LoadPackage(packagePath string, world w.World) []ecs.Entity {
	file, err := os.Open(packagePath)
	utils.LogError(err)
	defer file.Close()

	lines := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	utils.LogError(scanner.Err())

	entities, err := LoadLevel(world, lines)
	utils.LogError(err)
	return entities
}

// LoadLevel loads a level from text lines
func LoadLevel(world w.World, lines []string) ([]ecs.Entity, error) {
	entities := []ecs.Entity{}

	gridWidth := 0
	gridHeight := len(lines)
	playerCount, boxCount, goalCount := 0, 0, 0
	for _, line := range lines {
		gridWidth = math.Max(gridWidth, len(line))
		playerCount += strings.Count(line, string(charPlayer)) + strings.Count(line, string(charPlayerOnGoal))
		boxCount += strings.Count(line, string(charBox)) + strings.Count(line, string(charBoxOnGoal))
		goalCount += strings.Count(line, string(charGoal)) + strings.Count(line, string(charBoxOnGoal)) + strings.Count(line, string(charPlayerOnGoal))
	}

	if gridWidth > maxWidth || gridHeight > maxHeight {
		return []ecs.Entity{}, fmt.Errorf("level size must be less than 30x20")
	}
	if boxCount != goalCount {
		return []ecs.Entity{}, fmt.Errorf("invalid level: box count and goal count must be the same")
	}
	if boxCount == 0 {
		return []ecs.Entity{}, fmt.Errorf("invalid level: no box")
	}
	if playerCount != 1 {
		return []ecs.Entity{}, fmt.Errorf("invalid level: level must have one player")
	}

	lines = normalizeLevel(lines, gridWidth, gridHeight)

	for iLine := range lines {
		for iChar, char := range lines[iLine] {
			switch char {
			case charFloor:
				entities = append(entities, createFloorEntity(world, iChar, iLine))
			case charExterior:
				entities = append(entities, createExteriorEntity(world, iChar, iLine))
			case charWall:
				entities = append(entities, createWallEntity(world, iChar, iLine))
			case charGoal:
				entities = append(entities, createGoalEntity(world, iChar, iLine))
			case charBox:
				entities = append(entities, createFloorEntity(world, iChar, iLine))
				entities = append(entities, createBoxEntity(world, iChar, iLine))
			case charPlayer:
				entities = append(entities, createFloorEntity(world, iChar, iLine))
				entities = append(entities, createPlayerEntity(world, iChar, iLine))
			case charBoxOnGoal:
				entities = append(entities, createGoalEntity(world, iChar, iLine))
				entities = append(entities, createBoxEntity(world, iChar, iLine))
			case charPlayerOnGoal:
				entities = append(entities, createGoalEntity(world, iChar, iLine))
				entities = append(entities, createPlayerEntity(world, iChar, iLine))
			default:
				world.Manager.DeleteEntities(entities...)
				return []ecs.Entity{}, fmt.Errorf("invalid level: invalid char '%s'", string(char))
			}
		}
	}
	return entities, nil
}

func normalizeLevel(lines []string, gridWidth, gridHeight int) []string {
	grid := make([][]rune, len(lines))

	for iLine := range lines {
		chars := []rune(lines[iLine])

		// Replace floor chars
		for iChar := range chars {
			if chars[iChar] == charFloor1 || chars[iChar] == charFloor2 || chars[iChar] == charFloor3 {
				chars[iChar] = charFloor
			}
		}

		// Complete line to grid width
		deltaLen := gridWidth - len(chars)
		for iSlice := 0; iSlice < deltaLen; iSlice++ {
			chars = append(chars, charFloor)
		}

		grid[iLine] = chars
	}

	// Fill exterior
	for iLine := 0; iLine < gridHeight; iLine++ {
		fillExterior(grid, iLine, 0, gridWidth, gridHeight)
		fillExterior(grid, iLine, gridWidth-1, gridWidth, gridHeight)
	}

	for iCol := 0; iCol < gridWidth; iCol++ {
		fillExterior(grid, 0, iCol, gridWidth, gridHeight)
		fillExterior(grid, gridHeight-1, iCol, gridWidth, gridHeight)
	}

	// Center level to max width
	for iLine := range lines {
		padding := maxWidth - gridWidth
		lines[iLine] = strings.Repeat(string(charExterior), padding/2) + string(grid[iLine]) + strings.Repeat(string(charExterior), padding-padding/2)
	}

	// Center level to max height
	padding := make([]string, maxHeight-gridHeight)
	for iPadding := range padding {
		padding[iPadding] = strings.Repeat(string(charExterior), maxWidth)
	}
	lines = append(padding[:len(padding)/2], lines...)
	lines = append(lines, padding[len(padding)/2:]...)

	return lines
}

func fillExterior(grid [][]rune, posLine, posCol, gridWidth, gridHeight int) {
	if grid[posLine][posCol] != charFloor {
		return
	}

	fillQueue := &[]struct{ line, col int }{{posLine, posCol}}

	for len(*fillQueue) > 0 {
		elem := (*fillQueue)[0]
		*fillQueue = (*fillQueue)[1:]

		colLeft := elem.col
		for colLeft > 0 && grid[elem.line][colLeft-1] == charFloor {
			colLeft--
		}

		colRight := elem.col
		for colRight < gridWidth-1 && grid[elem.line][colRight+1] == charFloor {
			colRight++
		}

		for iCol := colLeft; iCol <= colRight; iCol++ {
			grid[elem.line][iCol] = charExterior

			if elem.line > 0 && grid[elem.line-1][iCol] == charFloor {
				*fillQueue = append(*fillQueue, struct{ line, col int }{elem.line - 1, iCol})
			}

			if elem.line < gridHeight-1 && grid[elem.line+1][iCol] == charFloor {
				*fillQueue = append(*fillQueue, struct{ line, col int }{elem.line + 1, iCol})
			}
		}
	}
}

func createFloorEntity(world w.World, posX, posY int) ecs.Entity {
	gameComponents := world.Components.Game.(*gc.Components)
	gameSpriteSheet := (*world.Resources.SpriteSheets)["game"]

	return world.Manager.NewEntity().
		AddComponent(world.Components.Engine.SpriteRender, &ec.SpriteRender{SpriteSheet: &gameSpriteSheet, SpriteNumber: resources.FloorSpriteNumber}).
		AddComponent(world.Components.Engine.Transform, &ec.Transform{}).
		AddComponent(gameComponents.GridElement, &gc.GridElement{PosX: posX, PosY: posY})
}

func createExteriorEntity(world w.World, posX, posY int) ecs.Entity {
	gameComponents := world.Components.Game.(*gc.Components)
	gameSpriteSheet := (*world.Resources.SpriteSheets)["game"]

	return world.Manager.NewEntity().
		AddComponent(world.Components.Engine.SpriteRender, &ec.SpriteRender{SpriteSheet: &gameSpriteSheet, SpriteNumber: resources.ExteriorSpriteNumber}).
		AddComponent(world.Components.Engine.Transform, &ec.Transform{}).
		AddComponent(gameComponents.GridElement, &gc.GridElement{PosX: posX, PosY: posY})
}

func createWallEntity(world w.World, posX, posY int) ecs.Entity {
	gameComponents := world.Components.Game.(*gc.Components)
	gameSpriteSheet := (*world.Resources.SpriteSheets)["game"]

	return world.Manager.NewEntity().
		AddComponent(world.Components.Engine.SpriteRender, &ec.SpriteRender{SpriteSheet: &gameSpriteSheet, SpriteNumber: resources.WallSpriteNumber}).
		AddComponent(world.Components.Engine.Transform, &ec.Transform{}).
		AddComponent(gameComponents.GridElement, &gc.GridElement{PosX: posX, PosY: posY})
}

func createGoalEntity(world w.World, posX, posY int) ecs.Entity {
	gameComponents := world.Components.Game.(*gc.Components)
	gameSpriteSheet := (*world.Resources.SpriteSheets)["game"]

	return world.Manager.NewEntity().
		AddComponent(world.Components.Engine.SpriteRender, &ec.SpriteRender{SpriteSheet: &gameSpriteSheet, SpriteNumber: resources.GoalSpriteNumber}).
		AddComponent(world.Components.Engine.Transform, &ec.Transform{}).
		AddComponent(gameComponents.GridElement, &gc.GridElement{PosX: posX, PosY: posY})
}

func createBoxEntity(world w.World, posX, posY int) ecs.Entity {
	gameComponents := world.Components.Game.(*gc.Components)
	gameSpriteSheet := (*world.Resources.SpriteSheets)["game"]

	return world.Manager.NewEntity().
		AddComponent(world.Components.Engine.SpriteRender, &ec.SpriteRender{SpriteSheet: &gameSpriteSheet, SpriteNumber: resources.BoxSpriteNumber}).
		AddComponent(world.Components.Engine.Transform, &ec.Transform{Depth: 1}).
		AddComponent(gameComponents.GridElement, &gc.GridElement{PosX: posX, PosY: posY})
}

func createPlayerEntity(world w.World, posX, posY int) ecs.Entity {
	gameComponents := world.Components.Game.(*gc.Components)
	gameSpriteSheet := (*world.Resources.SpriteSheets)["game"]

	return world.Manager.NewEntity().
		AddComponent(world.Components.Engine.SpriteRender, &ec.SpriteRender{SpriteSheet: &gameSpriteSheet, SpriteNumber: resources.PlayerSpriteNumber}).
		AddComponent(world.Components.Engine.Transform, &ec.Transform{Depth: 1}).
		AddComponent(gameComponents.GridElement, &gc.GridElement{PosX: posX, PosY: posY})
}
