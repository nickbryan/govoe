package blocks

import (
	"container/list"
	"math"

	"github.com/nickbryan/voxel/internal/oldernal/old/noise"

	"github.com/go-gl/mathgl/mgl64"
)

type Chunk struct {
	blocks   BlockContainer
	lightMap LightMapContainer
	mesher   *CulledMesher

	Mesh    MeshRenderer
	Pos     mgl64.Vec3
	GridPos mgl64.Vec3

	XMinus, XPlus *Chunk
	YMinus, YPlus *Chunk
	ZMinus, ZPlus *Chunk

	NumNeighbours int
}

func (c *Chunk) Setup() {
	c.blocks = make([]BlockType, ChunkSizeCubed)
	c.lightMap = make([]byte, ChunkSizeCubed)

	c.mesher = &CulledMesher{
		Mesh:     c.Mesh,
		Blocks:   c.blocks,
		LightMap: c.lightMap,
		Offset:   c.Pos,

		XMinus: c.XMinus,
		XPlus:  c.XPlus,
		YMinus: c.YMinus,
		YPlus:  c.YPlus,
		ZMinus: c.ZMinus,
		ZPlus:  c.ZPlus,
	}

	n1 := noise.NewCombined(noise.NewOctave(8), noise.NewOctave(8))
	n2 := noise.NewCombined(noise.NewOctave(8), noise.NewOctave(8))
	n3 := noise.NewOctave(6)

	for x := 0.0; x < ChunkSize; x++ {
		xPos := c.Pos.X() + x

		for z := 0.0; z < ChunkSize; z++ {
			zPos := c.Pos.Z() + z

			hLow := n1.Compute(xPos*1.3, zPos*1.3)/6 - 4
			height := hLow

			if n3.Compute(x, z) <= 0 {
				hHigh := n2.Compute(xPos*1.3, zPos*1.3)/5 + 6
				height = math.Max(hLow, hHigh)
			}

			height *= 0.5
			if height < 0 {
				height *= 0.8
			}

			adjHeight := height // + water level
			// TODO: cap this somehow

			for y := 0.0; y < ChunkSize; y++ {
				yPos := c.Pos.Y() + y
				if yPos < adjHeight {
					c.blocks.Set(int(x), int(y), int(z), Grass)
				} else {
					c.blocks.Set(int(x), int(y), int(z), Empty)
				}

				//c.lightMap.SetSunlight()
			}
		}
	}
}

func (c *Chunk) TearDown() {
	c.Mesh.TearDown()
}

func (c *Chunk) PlaceTorch(x, y, z int) {
	lightBfsQueue := list.New()

	c.lightMap.SetTorchlight(x, y, z, 15) // TODO: 16?
	lightBfsQueue.PushBack(&LightNode{Chunk: c, X: x, Y: y, Z: z})

	for lightBfsQueue.Len() > 0 {
		e := lightBfsQueue.Front()
		lightBfsQueue.Remove(e)
		node := e.Value.(*LightNode)

		lightLvl := node.Chunk.lightMap.Torchlight(node.X, node.Y, node.Z)

		if node.X != 0 && c.blocks.Lookup(node.X-1, node.Y, node.Z) == Empty && node.Chunk.lightMap.Torchlight(node.X-1, node.Y, node.Z)+2 <= lightLvl {
			node.Chunk.lightMap.SetTorchlight(node.X-1, node.Y, node.Z, lightLvl-1)
			lightBfsQueue.PushBack(&LightNode{Chunk: node.Chunk, X: node.X - 1, Y: node.Y, Z: node.Z})
		}

		if node.Y != 0 && c.blocks.Lookup(node.X, node.Y-1, node.Z) == Empty && node.Chunk.lightMap.Torchlight(node.X, node.Y-1, node.Z)+2 <= lightLvl {
			node.Chunk.lightMap.SetTorchlight(node.X, node.Y-1, node.Z, lightLvl-1)
			lightBfsQueue.PushBack(&LightNode{Chunk: node.Chunk, X: node.X, Y: node.Y - 1, Z: node.Z})
		}

		if node.Z != 0 && c.blocks.Lookup(node.X, node.Y, node.Z-1) == Empty && node.Chunk.lightMap.Torchlight(node.X, node.Y, node.Z-1)+2 <= lightLvl {
			node.Chunk.lightMap.SetTorchlight(node.X, node.Y, node.Z-1, lightLvl-1)
			lightBfsQueue.PushBack(&LightNode{Chunk: node.Chunk, X: node.X, Y: node.Y, Z: node.Z - 1})
		}

		if node.X < ChunkSize && c.blocks.Lookup(node.X+1, node.Y, node.Z) == Empty && node.Chunk.lightMap.Torchlight(node.X+1, node.Y, node.Z)+2 <= lightLvl {
			node.Chunk.lightMap.SetTorchlight(node.X+1, node.Y, node.Z, lightLvl-1)
			lightBfsQueue.PushBack(&LightNode{Chunk: node.Chunk, X: node.X + 1, Y: node.Y, Z: node.Z})
		}

		if node.Y < ChunkSize && c.blocks.Lookup(node.X, node.Y+1, node.Z) == Empty && node.Chunk.lightMap.Torchlight(node.X, node.Y+1, node.Z)+2 <= lightLvl {
			node.Chunk.lightMap.SetTorchlight(node.X, node.Y+1, node.Z, lightLvl-1)
			lightBfsQueue.PushBack(&LightNode{Chunk: node.Chunk, X: node.X, Y: node.Y + 1, Z: node.Z})
		}

		if node.Z < ChunkSize && c.blocks.Lookup(node.X, node.Y, node.Z+1) == Empty && node.Chunk.lightMap.Torchlight(node.X, node.Y, node.Z+1)+2 <= lightLvl {
			node.Chunk.lightMap.SetTorchlight(node.X, node.Y, node.Z+1, lightLvl-1)
			lightBfsQueue.PushBack(&LightNode{Chunk: node.Chunk, X: node.X, Y: node.Y, Z: node.Z + 1})
		}
	}

}

func (c *Chunk) BuildMesh() {
	c.mesher.Update()
}
