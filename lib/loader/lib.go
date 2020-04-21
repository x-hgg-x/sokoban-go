package loader

import gc "github.com/x-hgg-x/sokoban-go/lib/components"

type gameComponentList struct {
	GridElement *gc.GridElement
	Player      *gc.Player
	Box         *gc.Box
	Goal        *gc.Goal
	Wall        *gc.Wall
}
