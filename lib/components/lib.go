package components

import ecs "github.com/x-hgg-x/goecs/v2"

// Components contains references to all game components
type Components struct {
	GridElement *ecs.SliceComponent
}

// GridElement component
type GridElement struct {
	PosX int
	PosY int
}
