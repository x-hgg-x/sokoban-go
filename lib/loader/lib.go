package loader

import (
	"os"

	gc "github.com/x-hgg-x/sokoban-go/lib/components"

	"github.com/x-hgg-x/goecsengine/loader"
	"github.com/x-hgg-x/goecsengine/utils"
	w "github.com/x-hgg-x/goecsengine/world"
)

type gameComponentList struct {
	GridElement *gc.GridElement
	Player      *gc.Player
	Box         *gc.Box
	Goal        *gc.Goal
	Wall        *gc.Wall
}

// PreloadEntities preloads entities with components
func PreloadEntities(entityMetadataPath string, world w.World) loader.EntityComponentList {
	return loader.EntityComponentList{Engine: loader.LoadEngineComponents(utils.Try(os.ReadFile(entityMetadataPath)), world)}
}
