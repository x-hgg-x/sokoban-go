package resources

import "github.com/x-hgg-x/goecsengine/loader"

// MenuPrefabs contains menu prefabs
type MenuPrefabs struct {
	MainMenu          loader.EntityComponentList
	ChoosePackageMenu loader.EntityComponentList
	PauseMenu         loader.EntityComponentList
	LevelCompleteMenu loader.EntityComponentList
	HighscoresMenu    loader.EntityComponentList
	SolutionsMenu     loader.EntityComponentList
}

// GamePrefabs contains game prefabs
type GamePrefabs struct {
	LevelInfo   loader.EntityComponentList
	BoxInfo     loader.EntityComponentList
	StepInfo    loader.EntityComponentList
	PackageInfo loader.EntityComponentList
}

// Prefabs contains game prefabs
type Prefabs struct {
	Menu MenuPrefabs
	Game GamePrefabs
}
