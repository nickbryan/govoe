package entity

import (
	"github.com/go-gl/mathgl/mgl32"
)

const (
	defaultYaw         float32 = 90
	defaultPitch       float32 = 0
	defaultSpeed       float32 = 10 //2.5
	defaultSensitivity float32 = 0.1
	defaultZoom        float32 = 45.0 // FOV
)

type attachable interface {
	Pos() mgl32.Vec3
	Front() mgl32.Vec3
	Up() mgl32.Vec3
}

type Camera struct {
	*Movable

	attachedEntity attachable
}

func NewCamera() *Camera {
	c := &Camera{
		Movable: &Movable{
			pos:             mgl32.Vec3{0, 0, 10},
			worldUp:         mgl32.Vec3{0, 1, 0},
			yaw:             defaultYaw,
			pitch:           defaultPitch,
			front:           mgl32.Vec3{0, 0, -1},
			speed:           defaultSpeed,
			lookSensitivity: defaultSensitivity,
			zoom:            defaultZoom,
		},
	}

	c.updateVectors()

	return c
}

func (c *Camera) CreateViewMatrix() mgl32.Mat4 {
	return mgl32.LookAtV(c.Movable.pos, c.Movable.pos.Add(c.Movable.front), c.Movable.up)
}

func (c *Camera) Attach(a attachable) {
	c.attachedEntity = a
}

func (c *Camera) Update() {
	if c.attachedEntity == nil {
		return
	}

	// TODO: maybe move this to getters on the camera where it gets the attached or its own if no attached?
	c.Movable.pos = c.attachedEntity.Pos()
	c.Movable.front = c.attachedEntity.Front()
	c.Movable.up = c.attachedEntity.Up()
}
