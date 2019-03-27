package physics

import (
	"errors"
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/nickbryan/voxel/engine"
)

type TransformationComponent struct {
	Position mgl32.Vec3
}

type TransformationManager struct {
	components map[*engine.Entity]*TransformationComponent
}

func (tm *TransformationManager) Register(e *engine.Entity, tc *TransformationComponent) error {
	// TODO: is there a better way to handle this accross managers and systems?
	if tm.components == nil {
		tm.components = make(map[*engine.Entity]*TransformationComponent)
	}

	if _, ok := tm.components[e]; ok {
		return errors.New(
			fmt.Sprintf("an entity is already registered within the TransformationManager with id: %v", e.Id()),
		)
	}

	tm.components[e] = tc

	return nil
}

func (tm *TransformationManager) Translate(e *engine.Entity, v mgl32.Vec3) {
	t := mgl32.Translate3D(v.X(), v.Y(), v.Z())
	p := tm.components[e].Position
	tm.components[e].Position = mgl32.TransformCoordinate(p, t)
}
