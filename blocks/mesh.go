package blocks

import (
	"math"

	"github.com/go-gl/mathgl/mgl64"
)

type FaceSide = int

const (
	NorthFace FaceSide = iota
	SoutFace
	EastFace
	WestFace
	TopFace
	BottomFace
)

type Face struct {
	side FaceSide
}

type NaiveMesher struct {
	Mesh   MeshRenderer
	Blocks BlockContainer
}

func (nm *NaiveMesher) Update() {
	for x := 0.0; x < ChunkSize; x += 1 {
		for y := 0.0; y < ChunkSize; y += 1 {
			for z := 0.0; z < ChunkSize; z += 1 {
				if nm.Blocks.Lookup(int(x), int(y), int(z)) != Empty {
					nm.createCube(x, y, z)
				}
			}
		}
	}
}

func (nm *NaiveMesher) createCube(x, y, z float64) {
	p1 := mgl64.Vec3{x - 1, y - 1, z + 1}
	p2 := mgl64.Vec3{x + 1, y - 1, z + 1}
	p3 := mgl64.Vec3{x + 1, y + 1, z + 1}
	p4 := mgl64.Vec3{x - 1, y + 1, z + 1}
	p5 := mgl64.Vec3{x + 1, y - 1, z - 1}
	p6 := mgl64.Vec3{x - 1, y - 1, z - 1}
	p7 := mgl64.Vec3{x - 1, y + 1, z - 1}
	p8 := mgl64.Vec3{x + 1, y + 1, z - 1}

	// Front
	v1 := nm.Mesh.AddVertex(p1)
	v2 := nm.Mesh.AddVertex(p2)
	v3 := nm.Mesh.AddVertex(p3)
	v4 := nm.Mesh.AddVertex(p4)

	nm.Mesh.AddTriangle(v1, v2, v3)
	nm.Mesh.AddTriangle(v1, v3, v4)

	// Back
	v5 := nm.Mesh.AddVertex(p5)
	v6 := nm.Mesh.AddVertex(p6)
	v7 := nm.Mesh.AddVertex(p7)
	v8 := nm.Mesh.AddVertex(p8)

	nm.Mesh.AddTriangle(v5, v6, v7)
	nm.Mesh.AddTriangle(v5, v7, v8)

	// Right
	v2 = nm.Mesh.AddVertex(p2)
	v5 = nm.Mesh.AddVertex(p5)
	v8 = nm.Mesh.AddVertex(p8)
	v3 = nm.Mesh.AddVertex(p3)

	nm.Mesh.AddTriangle(v2, v5, v8)
	nm.Mesh.AddTriangle(v2, v8, v3)

	// Left
	v6 = nm.Mesh.AddVertex(p6)
	v1 = nm.Mesh.AddVertex(p1)
	v4 = nm.Mesh.AddVertex(p4)
	v7 = nm.Mesh.AddVertex(p7)

	nm.Mesh.AddTriangle(v6, v1, v4)
	nm.Mesh.AddTriangle(v6, v4, v7)

	// Top
	v4 = nm.Mesh.AddVertex(p4)
	v3 = nm.Mesh.AddVertex(p3)
	v8 = nm.Mesh.AddVertex(p8)
	v7 = nm.Mesh.AddVertex(p7)

	nm.Mesh.AddTriangle(v4, v3, v8)
	nm.Mesh.AddTriangle(v4, v8, v7)

	// Bottom
	v6 = nm.Mesh.AddVertex(p6)
	v5 = nm.Mesh.AddVertex(p5)
	v2 = nm.Mesh.AddVertex(p2)
	v1 = nm.Mesh.AddVertex(p1)

	nm.Mesh.AddTriangle(v6, v5, v2)
	nm.Mesh.AddTriangle(v6, v2, v1)
}

type CulledMesher struct {
	Mesh     MeshRenderer
	Blocks   BlockContainer
	LightMap LightMapContainer
	Offset   mgl64.Vec3
}

func (cm *CulledMesher) GetBlockColor(x, y, z int, bt BlockType, o mgl64.Vec3) mgl64.Vec3 {
	base := 0.86

	tl := 15.0
	if x > 0 && y > 0 && z > 0 && x < ChunkSize && y < ChunkSize && z < ChunkSize {
		tl = float64(cm.LightMap.Torchlight(x, y, z))
	}
	lightColor := math.Pow(tl/16.0, 1.4) + base

	if o.X() == 0 && o.Y() == 0 && o.Z() == 0 {
		return mgl64.Vec3{0, 0, 0}
	}

	switch bt {
	case Grass:
		return mgl64.Vec3{0.094 * lightColor, 0.568 * lightColor, 0.109 * lightColor}
	case Stone:
		return mgl64.Vec3{0.423, 0.478, 0.537}
	default:
		return mgl64.Vec3{} // TODO: handle error
	}
}

func (cm *CulledMesher) Update() {
	for x := 0.0; x < ChunkSize; x += 1 {
		for y := 0.0; y < ChunkSize; y += 1 {
			for z := 0.0; z < ChunkSize; z += 1 {
				if bt := cm.Blocks.Lookup(int(x), int(y), int(z)); bt != Empty {
					cm.createCube(x, y, z, bt)
				}
			}
		}
	}
}

func (cm *CulledMesher) createCube(x, y, z float64, bt BlockType) {
	p1 := mgl64.Vec3{cm.Offset.X() + x - BlockRenderSize, cm.Offset.Y() + y - BlockRenderSize, cm.Offset.Z() + z + BlockRenderSize}
	p2 := mgl64.Vec3{cm.Offset.X() + x + BlockRenderSize, cm.Offset.Y() + y - BlockRenderSize, cm.Offset.Z() + z + BlockRenderSize}
	p3 := mgl64.Vec3{cm.Offset.X() + x + BlockRenderSize, cm.Offset.Y() + y + BlockRenderSize, cm.Offset.Z() + z + BlockRenderSize}
	p4 := mgl64.Vec3{cm.Offset.X() + x - BlockRenderSize, cm.Offset.Y() + y + BlockRenderSize, cm.Offset.Z() + z + BlockRenderSize}
	p5 := mgl64.Vec3{cm.Offset.X() + x + BlockRenderSize, cm.Offset.Y() + y - BlockRenderSize, cm.Offset.Z() + z - BlockRenderSize}
	p6 := mgl64.Vec3{cm.Offset.X() + x - BlockRenderSize, cm.Offset.Y() + y - BlockRenderSize, cm.Offset.Z() + z - BlockRenderSize}
	p7 := mgl64.Vec3{cm.Offset.X() + x - BlockRenderSize, cm.Offset.Y() + y + BlockRenderSize, cm.Offset.Z() + z - BlockRenderSize}
	p8 := mgl64.Vec3{cm.Offset.X() + x + BlockRenderSize, cm.Offset.Y() + y + BlockRenderSize, cm.Offset.Z() + z - BlockRenderSize}

	var v1, v2, v3, v4, v5, v6, v7, v8 uint32

	// Front
	if (z == ChunkSize-1) || (z < ChunkSize-1 && cm.Blocks.Lookup(int(x), int(y), int(z)+1) == Empty) {
		cm.Mesh.SetColor(cm.GetBlockColor(int(x), int(y), int(z)+1, bt, cm.Offset))

		v1 = cm.Mesh.AddVertex(p1)
		v2 = cm.Mesh.AddVertex(p2)
		v3 = cm.Mesh.AddVertex(p3)
		v4 = cm.Mesh.AddVertex(p4)

		cm.Mesh.AddTriangle(v1, v2, v3)
		cm.Mesh.AddTriangle(v1, v3, v4)
	}

	// Back
	if z == 0 || (z > 0 && cm.Blocks.Lookup(int(x), int(y), int(z)-1) == Empty) {
		cm.Mesh.SetColor(cm.GetBlockColor(int(x), int(y), int(z)-1, bt, cm.Offset))

		v5 = cm.Mesh.AddVertex(p5)
		v6 = cm.Mesh.AddVertex(p6)
		v7 = cm.Mesh.AddVertex(p7)
		v8 = cm.Mesh.AddVertex(p8)

		cm.Mesh.AddTriangle(v5, v6, v7)
		cm.Mesh.AddTriangle(v5, v7, v8)
	}

	// Right
	if (x == ChunkSize-1) || (x < ChunkSize-1 && cm.Blocks.Lookup(int(x)+1, int(y), int(z)) == Empty) {
		cm.Mesh.SetColor(cm.GetBlockColor(int(x)+1, int(y), int(z), bt, cm.Offset))

		v2 = cm.Mesh.AddVertex(p2)
		v5 = cm.Mesh.AddVertex(p5)
		v8 = cm.Mesh.AddVertex(p8)
		v3 = cm.Mesh.AddVertex(p3)

		cm.Mesh.AddTriangle(v2, v5, v8)
		cm.Mesh.AddTriangle(v2, v8, v3)
	}

	// Left
	if x == 0 || (x > 0 && cm.Blocks.Lookup(int(x)-1, int(y), int(z)) == Empty) {
		cm.Mesh.SetColor(cm.GetBlockColor(int(x)-1, int(y), int(z), bt, cm.Offset))

		v6 = cm.Mesh.AddVertex(p6)
		v1 = cm.Mesh.AddVertex(p1)
		v4 = cm.Mesh.AddVertex(p4)
		v7 = cm.Mesh.AddVertex(p7)

		cm.Mesh.AddTriangle(v6, v1, v4)
		cm.Mesh.AddTriangle(v6, v4, v7)
	}

	// Top
	if (y == ChunkSize-1) || (y < ChunkSize-1 && cm.Blocks.Lookup(int(x), int(y)+1, int(z)) == Empty) {
		cm.Mesh.SetColor(cm.GetBlockColor(int(x), int(y)+1, int(z), bt, cm.Offset))

		v4 = cm.Mesh.AddVertex(p4)
		v3 = cm.Mesh.AddVertex(p3)
		v8 = cm.Mesh.AddVertex(p8)
		v7 = cm.Mesh.AddVertex(p7)

		cm.Mesh.AddTriangle(v4, v3, v8)
		cm.Mesh.AddTriangle(v4, v8, v7)
	}

	// Bottom
	if y == 0 || (y > 0 && cm.Blocks.Lookup(int(x), int(y)-1, int(z)) == Empty) {
		cm.Mesh.SetColor(cm.GetBlockColor(int(x), int(y)-1, int(z), bt, cm.Offset))

		v6 = cm.Mesh.AddVertex(p6)
		v5 = cm.Mesh.AddVertex(p5)
		v2 = cm.Mesh.AddVertex(p2)
		v1 = cm.Mesh.AddVertex(p1)

		cm.Mesh.AddTriangle(v6, v5, v2)
		cm.Mesh.AddTriangle(v6, v2, v1)
	}
}
