package renderer

import (
	"github.com/faiface/mainthread"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl64"
)

type Mesh struct {
	vao, vbo, ebo uint32
	vertices      []float32
	indices       []uint32
	vertexCount   uint32
	activeColor   mgl64.Vec3
	active        bool
}

func (m *Mesh) Setup() {
	m.active = true

	mainthread.Call(func() {
		gl.GenVertexArrays(1, &m.vao)
		gl.GenBuffers(1, &m.vbo)
		gl.GenBuffers(1, &m.ebo)
	})
}

func (m *Mesh) TearDown() {
	m.active = false

	mainthread.Call(func() {
		gl.DeleteVertexArrays(1, &m.vao)
		gl.DeleteBuffers(1, &m.vbo)
		gl.DeleteBuffers(1, &m.ebo)
	})
}

func (m *Mesh) Finish() {
	if m.vertexCount == 0 {
		return
	}

	mainthread.Call(func() {
		gl.BindVertexArray(m.vao)

		gl.BindBuffer(gl.ARRAY_BUFFER, m.vbo)
		gl.BufferData(gl.ARRAY_BUFFER, len(m.vertices)*4, gl.Ptr(m.vertices), gl.STATIC_DRAW)

		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, m.ebo)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(m.indices)*4, gl.Ptr(m.indices), gl.STATIC_DRAW)

		// Position
		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 24, gl.PtrOffset(0))
		gl.EnableVertexAttribArray(0)

		//Color
		// TODO: 24 = 6 * sizeof(float) and 12 = 3 * sizeof(float)
		gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 24, gl.PtrOffset(12))
		gl.EnableVertexAttribArray(1)

		// Ensure we unbind the VAO after so other VAO calls won't accidentally modify it.
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		gl.BindVertexArray(0)
	})
}

func (m *Mesh) AddVertex(p mgl64.Vec3) uint32 {
	m.vertices = append(
		m.vertices,
		float32(p.X()),
		float32(p.Y()),
		float32(p.Z()),
		float32(m.activeColor.X()),
		float32(m.activeColor.Y()),
		float32(m.activeColor.Z()),
	)

	vCount := m.vertexCount
	m.vertexCount += 1

	return vCount
}

func (m *Mesh) AddTriangle(v1, v2, v3 uint32) {
	m.indices = append(m.indices, v1, v2, v3)
}

func (m *Mesh) SetColor(c mgl64.Vec3) {
	m.activeColor = c
}
