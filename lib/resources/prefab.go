package resources

import "github.com/x-hgg-x/goecsengine/loader"

// MenuPrefabs contains menu prefabs
type MenuPrefabs struct {
	LevelCompleteMenu loader.EntityComponentList
}

// GamePrefabs contains game prefabs
type GamePrefabs struct {
	Levels    []loader.EntityComponentList
	LevelInfo loader.EntityComponentList
	BoxInfo   loader.EntityComponentList
	StepInfo  loader.EntityComponentList
}

// Prefabs contains game prefabs
type Prefabs struct {
	Menu MenuPrefabs
	Game GamePrefabs
}
