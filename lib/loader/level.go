package loader

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	gc "github.com/x-hgg-x/sokoban-go/lib/components"
	"github.com/x-hgg-x/sokoban-go/lib/resources"

	ec "github.com/x-hgg-x/goecsengine/components"
	"github.com/x-hgg-x/goecsengine/loader"
	"github.com/x-hgg-x/goecsengine/math"
	"github.com/x-hgg-x/goecsengine/utils"
	w "github.com/x-hgg-x/goecsengine/world"
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

var regexpValidChars = regexp.MustCompile(`^[ \-_#\.\$@\*\+]+$`)

// LoadPackage loads level package from a text file
func LoadPackage(packagePath string, world w.World) {
	prefabs := world.Resources.Prefabs.(*resources.Prefabs)

	// Load file
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
	lines = append(lines, "")

	// Split levels
	levels := [][]string{}
	currentLevel := []string{}
	for _, line := range lines {
		if len(strings.TrimSpace(line)) > 0 && regexpValidChars.MatchString(line) {
			currentLevel = append(currentLevel, line)
		} else if len(currentLevel) > 0 {
			levels = append(levels, currentLevel)
			currentLevel = []string{}
		}
	}

	// Preload levels
	prefabs.Game.Levels = []loader.EntityComponentList{}
	for iLevel, level := range levels {
		if componentList, err := PreloadLevel(world, level); err == nil {
			prefabs.Game.Levels = append(prefabs.Game.Levels, *componentList)
		} else {
			log.Printf("error when loading level %d: %s", iLevel+1, err.Error())
		}
	}

	if len(prefabs.Game.Levels) == 0 {
		utils.LogError(fmt.Errorf("invalid package: no valid levels in package"))
	}
}

// PreloadLevel preloads a level from text lines
func PreloadLevel(world w.World, lines []string) (*loader.EntityComponentList, error) {
	gridWidth := 0
	gridHeight := len(lines)
	playerCount, boxCount, goalCount := 0, 0, 0
	for _, line := range lines {
		gridWidth = math.Max(gridWidth, len(line))
		playerCount += strings.Count(line, string(charPlayer)) + strings.Count(line, string(charPlayerOnGoal))
		boxCount += strings.Count(line, string(charBox)) + strings.Count(line, string(charBoxOnGoal))
		goalCount += strings.Count(line, string(charGoal)) + strings.Count(line, string(charBoxOnGoal)) + strings.Count(line, string(charPlayerOnGoal))
	}

	if gridWidth > resources.MaxWidth || gridHeight > resources.MaxHeight {
		return nil, fmt.Errorf("level size must be less than 30x20")
	}
	if boxCount != goalCount {
		return nil, fmt.Errorf("invalid level: box count and goal count must be the same")
	}
	if boxCount == 0 {
		return nil, fmt.Errorf("invalid level: no box")
	}
	if playerCount != 1 {
		return nil, fmt.Errorf("invalid level: level must have one player")
	}

	lines = normalizeLevel(lines, gridWidth, gridHeight)

	componentList := &loader.EntityComponentList{}
	gameSpriteSheet := (*world.Resources.SpriteSheets)["game"]

	for iLine := range lines {
		for iChar, char := range lines[iLine] {
			switch char {
			case charFloor:
				createFloorEntity(componentList, &gameSpriteSheet, iLine, iChar)
			case charExterior:
				createExteriorEntity(componentList, &gameSpriteSheet, iLine, iChar)
			case charWall:
				createWallEntity(componentList, &gameSpriteSheet, iLine, iChar)
			case charGoal:
				createGoalEntity(componentList, &gameSpriteSheet, iLine, iChar)
			case charBox:
				createFloorEntity(componentList, &gameSpriteSheet, iLine, iChar)
				createBoxEntity(componentList, &gameSpriteSheet, iLine, iChar)
			case charPlayer:
				createFloorEntity(componentList, &gameSpriteSheet, iLine, iChar)
				createPlayerEntity(componentList, &gameSpriteSheet, iLine, iChar)
			case charBoxOnGoal:
				createGoalEntity(componentList, &gameSpriteSheet, iLine, iChar)
				createBoxEntity(componentList, &gameSpriteSheet, iLine, iChar)
			case charPlayerOnGoal:
				createGoalEntity(componentList, &gameSpriteSheet, iLine, iChar)
				createPlayerEntity(componentList, &gameSpriteSheet, iLine, iChar)
			default:
				return nil, fmt.Errorf("invalid level: invalid char '%s'", string(char))
			}
		}
	}
	return componentList, nil
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
		padding := resources.MaxWidth - gridWidth
		lines[iLine] = strings.Repeat(string(charExterior), padding/2) + string(grid[iLine]) + strings.Repeat(string(charExterior), padding-padding/2)
	}

	// Center level to max height
	padding := make([]string, resources.MaxHeight-gridHeight)
	for iPadding := range padding {
		padding[iPadding] = strings.Repeat(string(charExterior), resources.MaxWidth)
	}
	lines = append(padding[:len(padding)/2], lines...)
	lines = append(lines, padding[len(padding)/2:]...)

	return lines
}

func fillExterior(grid [][]rune, line, col, gridWidth, gridHeight int) {
	if grid[line][col] != charFloor {
		return
	}

	fillQueue := &[]struct{ line, col int }{{line, col}}

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

func createFloorEntity(componentList *loader.EntityComponentList, gameSpriteSheet *ec.SpriteSheet, line, col int) {
	componentList.Engine = append(componentList.Engine, loader.EngineComponentList{
		SpriteRender: &ec.SpriteRender{SpriteSheet: gameSpriteSheet, SpriteNumber: resources.FloorSpriteNumber},
		Transform:    &ec.Transform{},
	})
	componentList.Game = append(componentList.Game, gameComponentList{
		GridElement: &gc.GridElement{Line: line, Col: col},
	})
}

func createExteriorEntity(componentList *loader.EntityComponentList, gameSpriteSheet *ec.SpriteSheet, line, col int) {
	componentList.Engine = append(componentList.Engine, loader.EngineComponentList{
		SpriteRender: &ec.SpriteRender{SpriteSheet: gameSpriteSheet, SpriteNumber: resources.ExteriorSpriteNumber},
		Transform:    &ec.Transform{},
	})
	componentList.Game = append(componentList.Game, gameComponentList{
		GridElement: &gc.GridElement{Line: line, Col: col},
	})
}

func createWallEntity(componentList *loader.EntityComponentList, gameSpriteSheet *ec.SpriteSheet, line, col int) {
	componentList.Engine = append(componentList.Engine, loader.EngineComponentList{
		SpriteRender: &ec.SpriteRender{SpriteSheet: gameSpriteSheet, SpriteNumber: resources.WallSpriteNumber},
		Transform:    &ec.Transform{},
	})
	componentList.Game = append(componentList.Game, gameComponentList{
		Wall:        &gc.Wall{},
		GridElement: &gc.GridElement{Line: line, Col: col},
	})
}

func createGoalEntity(componentList *loader.EntityComponentList, gameSpriteSheet *ec.SpriteSheet, line, col int) {
	componentList.Engine = append(componentList.Engine, loader.EngineComponentList{
		SpriteRender: &ec.SpriteRender{SpriteSheet: gameSpriteSheet, SpriteNumber: resources.GoalSpriteNumber},
		Transform:    &ec.Transform{},
	})
	componentList.Game = append(componentList.Game, gameComponentList{
		Goal:        &gc.Goal{},
		GridElement: &gc.GridElement{Line: line, Col: col},
	})
}

func createBoxEntity(componentList *loader.EntityComponentList, gameSpriteSheet *ec.SpriteSheet, line, col int) {
	componentList.Engine = append(componentList.Engine, loader.EngineComponentList{
		SpriteRender: &ec.SpriteRender{SpriteSheet: gameSpriteSheet, SpriteNumber: resources.BoxSpriteNumber},
		Transform:    &ec.Transform{Depth: 1},
	})
	componentList.Game = append(componentList.Game, gameComponentList{
		Box:         &gc.Box{},
		GridElement: &gc.GridElement{Line: line, Col: col},
	})
}

func createPlayerEntity(componentList *loader.EntityComponentList, gameSpriteSheet *ec.SpriteSheet, line, col int) {
	componentList.Engine = append(componentList.Engine, loader.EngineComponentList{
		SpriteRender: &ec.SpriteRender{SpriteSheet: gameSpriteSheet, SpriteNumber: resources.PlayerSpriteNumber},
		Transform:    &ec.Transform{Depth: 1},
	})
	componentList.Game = append(componentList.Game, gameComponentList{
		Player:      &gc.Player{},
		GridElement: &gc.GridElement{Line: line, Col: col},
	})
}
