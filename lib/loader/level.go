package loader

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	gc "github.com/x-hgg-x/sokoban-go/lib/components"

	ec "github.com/x-hgg-x/goecsengine/components"
	"github.com/x-hgg-x/goecsengine/loader"
	"github.com/x-hgg-x/goecsengine/math"
	"github.com/x-hgg-x/goecsengine/utils"
)

// MaxGameSize is the maximum game size
const MaxGameSize = 100

const (
	exteriorSpriteNumber = 0
	wallSpriteNumber     = 1
	floorSpriteNumber    = 2
	goalSpriteNumber     = 3
	boxSpriteNumber      = 4
	playerSpriteNumber   = 5
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
func LoadPackage(packageName string) (packageLevels [][][]byte, packageErr error) {
	// Load file
	file := utils.Try(os.Open(fmt.Sprintf("levels/%s/levels.txt", packageName)))
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
		if len(strings.TrimSpace(line)) == 0 && len(currentLevel) > 0 {
			levels = append(levels, currentLevel)
			currentLevel = []string{}
		} else if regexpValidChars.MatchString(line) {
			currentLevel = append(currentLevel, line)
		}
	}

	// Normalize levels
	for iLevel, level := range levels {
		if grid, err := normalizeLevel(level); err == nil {
			packageLevels = append(packageLevels, grid)
		} else {
			packageErr = fmt.Errorf("error when loading level %d: %s", iLevel+1, err.Error())
			break
		}
	}

	if len(packageLevels) == 0 {
		if packageErr != nil {
			log.Println(packageErr)
		}
		utils.LogFatalf("invalid package: no valid levels in package")
	}
	return
}

func normalizeLevel(lines []string) ([][]byte, error) {
	gridWidth := 0
	gridHeight := len(lines)
	playerCount, boxCount, goalCount := 0, 0, 0
	for _, line := range lines {
		gridWidth = math.Max(gridWidth, len(line))
		playerCount += strings.Count(line, string(charPlayer)) + strings.Count(line, string(charPlayerOnGoal))
		boxCount += strings.Count(line, string(charBox)) + strings.Count(line, string(charBoxOnGoal))
		goalCount += strings.Count(line, string(charGoal)) + strings.Count(line, string(charBoxOnGoal)) + strings.Count(line, string(charPlayerOnGoal))
	}

	if gridWidth > MaxGameSize || gridHeight > MaxGameSize {
		return nil, fmt.Errorf("level size must be less than %dx%d", MaxGameSize, MaxGameSize)
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

	grid := make([][]byte, len(lines))

	for iLine := range lines {
		chars := []byte(lines[iLine])

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

	return grid, nil
}

func fillExterior(grid [][]byte, line, col, gridWidth, gridHeight int) {
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

// LoadLevel loads a level from a grid
func LoadLevel(grid [][]byte, maxWidth, maxHeight int, gameSpriteSheet *ec.SpriteSheet) (loader.EntityComponentList, error) {
	componentList := loader.EntityComponentList{}

	horizontalPadding := maxWidth - len(grid[0])
	horizontalPaddingBefore := horizontalPadding / 2
	horizontalPaddingAfter := horizontalPadding - horizontalPaddingBefore

	verticalPadding := maxHeight - len(grid)
	verticalPaddingBefore := verticalPadding / 2
	verticalPaddingAfter := verticalPadding - verticalPaddingBefore

	for iLine := 0; iLine < verticalPaddingBefore; iLine++ {
		for iCol := 0; iCol < maxWidth; iCol++ {
			createExteriorEntity(&componentList, gameSpriteSheet, iLine, iCol)
		}
	}

	for iGridLine := range grid {
		iLine := iGridLine + verticalPaddingBefore

		for iCol := 0; iCol < horizontalPaddingBefore; iCol++ {
			createExteriorEntity(&componentList, gameSpriteSheet, iLine, iCol)
		}

		for iGridCol, char := range grid[iGridLine] {
			iCol := iGridCol + horizontalPaddingBefore

			switch char {
			case charFloor:
				createFloorEntity(&componentList, gameSpriteSheet, iLine, iCol)
			case charExterior:
				createExteriorEntity(&componentList, gameSpriteSheet, iLine, iCol)
			case charWall:
				createWallEntity(&componentList, gameSpriteSheet, iLine, iCol)
			case charGoal:
				createGoalEntity(&componentList, gameSpriteSheet, iLine, iCol)
			case charBox:
				createFloorEntity(&componentList, gameSpriteSheet, iLine, iCol)
				createBoxEntity(&componentList, gameSpriteSheet, iLine, iCol)
			case charPlayer:
				createFloorEntity(&componentList, gameSpriteSheet, iLine, iCol)
				createPlayerEntity(&componentList, gameSpriteSheet, iLine, iCol)
			case charBoxOnGoal:
				createGoalEntity(&componentList, gameSpriteSheet, iLine, iCol)
				createBoxEntity(&componentList, gameSpriteSheet, iLine, iCol)
			case charPlayerOnGoal:
				createGoalEntity(&componentList, gameSpriteSheet, iLine, iCol)
				createPlayerEntity(&componentList, gameSpriteSheet, iLine, iCol)
			default:
				return loader.EntityComponentList{}, fmt.Errorf("invalid level: invalid char '%c'", char)
			}
		}

		for iCol := maxWidth - horizontalPaddingAfter; iCol < maxWidth; iCol++ {
			createExteriorEntity(&componentList, gameSpriteSheet, iLine, iCol)
		}
	}

	for iLine := maxHeight - verticalPaddingAfter; iLine < maxHeight; iLine++ {
		for iCol := 0; iCol < maxWidth; iCol++ {
			createExteriorEntity(&componentList, gameSpriteSheet, iLine, iCol)
		}
	}

	return componentList, nil
}

func createFloorEntity(componentList *loader.EntityComponentList, gameSpriteSheet *ec.SpriteSheet, line, col int) {
	componentList.Engine = append(componentList.Engine, loader.EngineComponentList{
		SpriteRender: &ec.SpriteRender{SpriteSheet: gameSpriteSheet, SpriteNumber: floorSpriteNumber},
		Transform:    &ec.Transform{},
	})
	componentList.Game = append(componentList.Game, gameComponentList{
		GridElement: &gc.GridElement{Line: line, Col: col},
	})
}

func createExteriorEntity(componentList *loader.EntityComponentList, gameSpriteSheet *ec.SpriteSheet, line, col int) {
	componentList.Engine = append(componentList.Engine, loader.EngineComponentList{
		SpriteRender: &ec.SpriteRender{SpriteSheet: gameSpriteSheet, SpriteNumber: exteriorSpriteNumber},
		Transform:    &ec.Transform{},
	})
	componentList.Game = append(componentList.Game, gameComponentList{
		GridElement: &gc.GridElement{Line: line, Col: col},
	})
}

func createWallEntity(componentList *loader.EntityComponentList, gameSpriteSheet *ec.SpriteSheet, line, col int) {
	componentList.Engine = append(componentList.Engine, loader.EngineComponentList{
		SpriteRender: &ec.SpriteRender{SpriteSheet: gameSpriteSheet, SpriteNumber: wallSpriteNumber},
		Transform:    &ec.Transform{},
	})
	componentList.Game = append(componentList.Game, gameComponentList{
		Wall:        &gc.Wall{},
		GridElement: &gc.GridElement{Line: line, Col: col},
	})
}

func createGoalEntity(componentList *loader.EntityComponentList, gameSpriteSheet *ec.SpriteSheet, line, col int) {
	componentList.Engine = append(componentList.Engine, loader.EngineComponentList{
		SpriteRender: &ec.SpriteRender{SpriteSheet: gameSpriteSheet, SpriteNumber: goalSpriteNumber},
		Transform:    &ec.Transform{},
	})
	componentList.Game = append(componentList.Game, gameComponentList{
		Goal:        &gc.Goal{},
		GridElement: &gc.GridElement{Line: line, Col: col},
	})
}

func createBoxEntity(componentList *loader.EntityComponentList, gameSpriteSheet *ec.SpriteSheet, line, col int) {
	componentList.Engine = append(componentList.Engine, loader.EngineComponentList{
		SpriteRender: &ec.SpriteRender{SpriteSheet: gameSpriteSheet, SpriteNumber: boxSpriteNumber},
		Transform:    &ec.Transform{Depth: 1},
	})
	componentList.Game = append(componentList.Game, gameComponentList{
		Box:         &gc.Box{},
		GridElement: &gc.GridElement{Line: line, Col: col},
	})
}

func createPlayerEntity(componentList *loader.EntityComponentList, gameSpriteSheet *ec.SpriteSheet, line, col int) {
	componentList.Engine = append(componentList.Engine, loader.EngineComponentList{
		SpriteRender: &ec.SpriteRender{SpriteSheet: gameSpriteSheet, SpriteNumber: playerSpriteNumber},
		Transform:    &ec.Transform{Depth: 1},
	})
	componentList.Game = append(componentList.Game, gameComponentList{
		Player:      &gc.Player{},
		GridElement: &gc.GridElement{Line: line, Col: col},
	})
}
