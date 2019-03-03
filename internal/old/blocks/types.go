package blocks

import (
	"github.com/go-gl/mathgl/mgl64"
)

type BlockType int

const (
	Empty BlockType = iota
	Grass
	Stone
)

const (
	BlockRenderSize float64 = 0.5
	BlockSize               = BlockRenderSize * 2
)

const (
	ChunkSize        = 32
	ChunkSizeSquared = ChunkSize * ChunkSize
	ChunkSizeCubed   = ChunkSize * ChunkSize * ChunkSize
)

type MeshRenderer interface {
	AddVertex(p mgl64.Vec3) uint32
	AddTriangle(v1, v2, v3 uint32)
	SetColor(c mgl64.Vec3)
	Finish()
	TearDown()
}

type LightNode struct {
	Chunk   *Chunk
	X, Y, Z int
}
