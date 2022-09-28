package resources

import (
	"fmt"
	"os"
	"strings"

	"github.com/x-hgg-x/goecsengine/utils"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/BurntSushi/toml"
)

// SaveConfig contains the save configuration
type SaveConfig struct {
	CurrentLevel   int
	Package        string
	LevelMovements map[string]string
}

// Encode encodes the save configuration
func (sc *SaveConfig) Encode() EncodedSaveConfig {
	esc := make(EncodedSaveConfig, len(sc.LevelMovements)+2)

	esc["CurrentLevel"] = sc.CurrentLevel
	esc["Package"] = sc.Package

	for k, v := range sc.LevelMovements {
		esc[k] = map[string]string{"Movements": v}
	}

	return esc
}

// EncodedSaveConfig contains the encoded save configuration
type EncodedSaveConfig map[string]interface{}

// Decode decodes the encoded save configuration
func (esc *EncodedSaveConfig) Decode() (sc SaveConfig, err error) {
	data := *esc

	if len(data) == 0 {
		return SaveConfig{LevelMovements: make(map[string]string)}, nil
	}

	var currentLevel int
	if v, ok := data["CurrentLevel"]; !ok {
		return sc, fmt.Errorf("invalid TOML file")
	} else if currentLevelField, ok := v.(int64); !ok {
		return sc, fmt.Errorf("invalid TOML file")
	} else {
		currentLevel = int(currentLevelField)
		delete(data, "CurrentLevel")
	}

	var packageName string
	if v, ok := data["Package"]; !ok {
		return sc, fmt.Errorf("invalid TOML file")
	} else if packageNameField, ok := v.(string); !ok {
		return sc, fmt.Errorf("invalid TOML file")
	} else {
		packageName = packageNameField
		delete(data, "Package")
	}

	levelMovements := make(map[string]string, len(data))
	for k, v := range data {
		if vm, ok := v.(map[string]interface{}); !ok {
			return sc, fmt.Errorf("invalid TOML file")
		} else if mv, ok := vm["Movements"].(string); !ok {
			return sc, fmt.Errorf("invalid TOML file")
		} else {
			levelMovements[k] = mv
		}
	}

	sc = SaveConfig{
		CurrentLevel:   currentLevel,
		Package:        packageName,
		LevelMovements: levelMovements,
	}

	return sc, nil
}

// SaveLevel saves level
func SaveLevel(world w.World) {
	gameResources := world.Resources.Game.(*Game)
	if !(world.Resources.InputHandler.Actions[SaveAction] || gameResources.Level.Modified) {
		return
	}

	saveFilePath := fmt.Sprintf("config/%s/save.toml", gameResources.Package.Name)

	// Load existing save file
	var encodedSaveConfig EncodedSaveConfig
	if saveFile, err := os.ReadFile(saveFilePath); err == nil {
		utils.Try(toml.Decode(string(saveFile), &encodedSaveConfig))
	}
	saveConfig := utils.Try(encodedSaveConfig.Decode())

	// Encode movements
	var movements strings.Builder
	for _, movement := range gameResources.Level.Movements {
		utils.LogError(movements.WriteByte(movementChars[movement]))
	}

	// Update save config
	saveConfig.CurrentLevel = gameResources.Level.CurrentNum + 1
	saveConfig.Package = gameResources.Package.Name
	saveConfig.LevelMovements[fmt.Sprintf("Level%04d", gameResources.Level.CurrentNum+1)] = movements.String()

	// Write to save file
	var encoded strings.Builder
	encoder := toml.NewEncoder(&encoded)
	encoder.Indent = ""
	utils.LogError(encoder.Encode(saveConfig.Encode()))
	utils.LogError(os.WriteFile(saveFilePath, []byte(encoded.String()), 0o666))

	gameResources.Level.Modified = false
}

// LoadSave loads save for a level
func LoadSave(world w.World) {
	gameResources := world.Resources.Game.(*Game)

	var encodedSaveConfig EncodedSaveConfig
	_, err := toml.DecodeFile(fmt.Sprintf("config/%s/save.toml", gameResources.Package.Name), &encodedSaveConfig)
	if err != nil {
		return
	}

	saveConfig, err := encodedSaveConfig.Decode()
	if err != nil {
		return
	}

	// Decode movements
	movements := []MovementType{}
	for _, char := range []byte(saveConfig.LevelMovements[fmt.Sprintf("Level%04d", gameResources.Level.CurrentNum+1)]) {
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
