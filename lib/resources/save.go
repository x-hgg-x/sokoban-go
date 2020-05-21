package resources

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/x-hgg-x/goecsengine/math"
	"github.com/x-hgg-x/goecsengine/utils"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/pelletier/go-toml"
)

const maxLineChars = 100

// LoadSaveFile loads save file
func LoadSaveFile(world w.World) *toml.Tree {
	gameResources := world.Resources.Game.(*Game)

	var tree *toml.Tree
	packageConfig := fmt.Sprintf("config/%s", gameResources.PackageName)
	if _, err := os.Stat(packageConfig + "/save.toml"); err == nil {
		tree, err = toml.LoadFile(packageConfig + "/save.toml")
		utils.LogError(err)
	} else {
		utils.LogError(os.MkdirAll(packageConfig, os.ModePerm))
		tree, err = toml.Load("")
		utils.LogError(err)
	}
	return tree
}

// SaveLevel saves level
func SaveLevel(world w.World) {
	gameResources := world.Resources.Game.(*Game)
	if !gameResources.Level.Modified {
		return
	}

	// Load existing save file
	tree := LoadSaveFile(world)

	// Encode movements
	movements := &bytes.Buffer{}
	for _, movement := range gameResources.Level.Movements {
		movements.WriteRune(movementChars[movement])
	}

	// Movements are written as an array to avoid long lines
	movementsList := []rune(movements.String())
	movementsArray := []string{}

	for start := 0; start < len(movementsList); start += maxLineChars {
		end := math.Min(start+maxLineChars, len(movementsList))
		movementsArray = append(movementsArray, string(movementsList[start:end]))
	}

	// Write tree
	tree.Set(fmt.Sprintf("Level%d.Movements", gameResources.Level.CurrentNum+1), movementsArray)
	tree.Set("Package", gameResources.PackageName)
	tree.Set("CurrentLevel", int64(gameResources.Level.CurrentNum+1))

	// Write to save file
	saveFile, err := os.Create(fmt.Sprintf("config/%s/save.toml", gameResources.PackageName))
	utils.LogError(err)
	defer saveFile.Close()

	err = toml.NewEncoder(saveFile).Indentation("    ").ArraysWithOneElementPerLine(true).Encode(tree)
	utils.LogError(err)

	gameResources.Level.Modified = false
}

// LoadSave loads save for a level
func LoadSave(world w.World) {
	gameResources := world.Resources.Game.(*Game)

	packageConfig := fmt.Sprintf("config/%s", gameResources.PackageName)
	tree, err := toml.LoadFile(packageConfig + "/save.toml")
	if err != nil {
		return
	}

	// Read movements array
	levelArrayMovements, ok := tree.Get(fmt.Sprintf("Level%d.Movements", gameResources.Level.CurrentNum+1)).([]interface{})
	if !ok {
		return
	}

	levelMovements := strings.Builder{}
	for iMovement := range levelArrayMovements {
		s, ok := levelArrayMovements[iMovement].(string)
		if !ok {
			return
		}
		levelMovements.WriteString(s)
	}

	// Decode movements
	movements := []MovementType{}
	for _, char := range levelMovements.String() {
		if movement, ok := movementCharMap[char]; ok {
			movements = append(movements, movement)
		} else {
			fmt.Printf("unknown movement when loading save: '%c'\n", char)
			return
		}
	}

	Move(world, movements...)
	gameResources.Level.Modified = false
}
