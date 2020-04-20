package resources

import "github.com/x-hgg-x/goecsengine/loader"

// MenuPrefabs contains menu prefabs
type MenuPrefabs struct{}

// GamePrefabs contains game prefabs
type GamePrefabs struct {
	LevelInfo loader.EntityComponentList
	BoxInfo   loader.EntityComponentList
	StepInfo  loader.EntityComponentList
}

// Prefabs contains game prefabs
type Prefabs struct {
	Menu MenuPrefabs
	Game GamePrefabs
}