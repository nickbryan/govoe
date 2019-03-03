package entity

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Player struct {
	*Movable
}

func NewPlayer() *Player {
	p := &Player{
		Movable: &Movable{
			pos:             mgl32.Vec3{16, 10, 16},
			worldUp:         mgl32.Vec3{0, 1, 0},
			yaw:             defaultYaw,
			pitch:           defaultPitch,
			front:           mgl32.Vec3{0, 0, -1},
			speed:           defaultSpeed * 3, // TODO: remove this as i doubled for testing
			lookSensitivity: defaultSensitivity,
			zoom:            defaultZoom,
		},
	}

	p.updateVectors()

	return p
}

func (p *Player) Pos() mgl32.Vec3 {
	return p.Movable.pos
}

func (p *Player) Front() mgl32.Vec3 {
	return p.Movable.front
}

func (p *Player) Up() mgl32.Vec3 {
	return p.Movable.up
}
