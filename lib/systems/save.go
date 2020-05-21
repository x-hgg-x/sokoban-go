package systems

import (
	"github.com/x-hgg-x/sokoban-go/lib/resources"

	w "github.com/x-hgg-x/goecsengine/world"
)

// SaveSystem saves current level
func SaveSystem(world w.World) {
	if world.Resources.InputHandler.Actions[resources.SaveAction] {
		resources.SaveLevel(world)
	}
}
