package components

import ecs "github.com/x-hgg-x/goecs/v2"

// Components contains references to all game components
type Components struct {
	GridElement *ecs.SliceComponent
	Player      *ecs.NullComponent
	Box         *ecs.NullComponent
	Goal        *ecs.NullComponent
	Wall        *ecs.NullComponent
}

// GridElement component
type GridElement struct {
	Line int
	Col  int
}

// Player component
type Player struct{}

// Box component
type Box struct{}

// Goal component
type Goal struct{}

// Wall component
type Wall struct{}
