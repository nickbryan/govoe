package physics

import (
	"errors"
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/nickbryan/voxel/engine"
)

type Direction int

const (
	Forward Direction = iota
	Backward
	Left
	Right
	Up
	Down
)

type MovementComponent struct {
	Acceleration mgl32.Vec3 // TODO: figure these out
	Velocity     mgl32.Vec3
}

type MovementSystem struct {
	TransformationManager *TransformationManager
	components            map[*engine.Entity]*MovementComponent
	commands              map[*engine.Entity][]Direction
}

func NewMovementSystem(tm *TransformationManager) *MovementSystem {
	return &MovementSystem{
		TransformationManager: tm,
		components:            make(map[*engine.Entity]*MovementComponent),
		commands:              make(map[*engine.Entity][]Direction),
	}
}

func (ms *MovementSystem) Simulate(dt float64) {
	for e, ds := range ms.commands {
		//c := ms.components[e]

		for _, d := range ds {
			switch d {
			case Forward:
				// TODO: figure this out with direction vectors
				ms.TransformationManager.Translate(e, mgl32.Vec3{0, 0, 1})
			case Backward:
				ms.TransformationManager.Translate(e, mgl32.Vec3{0, 0, -1})
			case Left:
				ms.TransformationManager.Translate(e, mgl32.Vec3{1, 0, 0})
			case Right:
				ms.TransformationManager.Translate(e, mgl32.Vec3{-1, 0, 0})
			case Up:
				ms.TransformationManager.Translate(e, mgl32.Vec3{0, 1, 0})
			case Down:
				ms.TransformationManager.Translate(e, mgl32.Vec3{0, -1, 0})
			}
		}
	}

	ms.commands = make(map[*engine.Entity][]Direction)
}

func (ms *MovementSystem) Register(e *engine.Entity, c *MovementComponent) error {
	if ms.components == nil {
		ms.components = make(map[*engine.Entity]*MovementComponent)
	}

	if _, ok := ms.components[e]; ok {
		return errors.New(
			fmt.Sprintf("an entity is already registered within the MovementSystem with id: %v", e.Id()),
		)
	}

	ms.components[e] = c

	return nil
}

func (ms *MovementSystem) Move(e *engine.Entity, d Direction) {
	if ms.commands == nil {
		ms.commands = make(map[*engine.Entity][]Direction)
	}

	ms.commands[e] = append(ms.commands[e], d)
}
