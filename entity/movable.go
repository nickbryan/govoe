package entity

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

const pitchCap = 89

type Movable struct {
	pos, front, right, up, worldUp           mgl32.Vec3
	speed, yaw, pitch, lookSensitivity, zoom float32
}

func (m *Movable) Stride(by float32) {
	v := m.velocity(by)

	if by > 0 {
		m.pos = m.pos.Add(m.front.Mul(v))
	}

	if by < 0 {
		m.pos = m.pos.Sub(m.front.Mul(v))
	}
}

func (m *Movable) Strafe(by float32) {
	v := m.velocity(by)

	if by > 0 {
		m.pos = m.pos.Add(m.right.Mul(v))
	}

	if by < 0 {
		m.pos = m.pos.Sub(m.right.Mul(v))
	}
}

func (m *Movable) Climb(by float32) {
	v := m.velocity(by)

	if by > 0 {
		m.pos = m.pos.Add(m.worldUp.Mul(v))
	}

	if by < 0 {
		m.pos = m.pos.Sub(m.worldUp.Mul(v))
	}
}

func (m *Movable) Look(offsetX, offsetY float32) {
	offsetX *= m.lookSensitivity
	offsetY *= m.lookSensitivity

	m.yaw += offsetX
	m.pitch += offsetY

	m.pitch = float32(math.Max(math.Min(float64(m.pitch), pitchCap), -pitchCap))

	m.updateVectors()
}

func (m *Movable) Zoom(offsetY float32) {
	m.zoom = float32(math.Max(math.Min(float64(m.zoom-(offsetY*m.lookSensitivity)), 45), 44.5))
}

func (m *Movable) velocity(by float32) float32 {
	return m.speed * float32(math.Abs(float64(by)))
}

func (m *Movable) updateVectors() {
	m.front = mgl32.Vec3{
		float32(math.Cos(float64(mgl32.DegToRad(m.yaw))) * math.Cos(float64(mgl32.DegToRad(m.pitch)))),
		float32(math.Sin(float64(mgl32.DegToRad(m.pitch)))),
		float32(math.Sin(float64(mgl32.DegToRad(m.yaw))) * math.Cos(float64(mgl32.DegToRad(m.pitch)))),
	}.Normalize()

	// Normalize the vectors, because their length gets closer to 0 the more you look up or down which results in slower movement.
	m.right = m.front.Cross(m.worldUp).Normalize()
	m.up = m.right.Cross(m.front).Normalize()
}
