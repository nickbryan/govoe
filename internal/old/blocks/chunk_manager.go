package blocks

import (
	"time"

	"github.com/go-gl/mathgl/mgl64"

	"github.com/nickbryan/voxel/internal/oldernal/old/entity"
	"github.com/nickbryan/voxel/internal/oldernal/old/renderer"
)

type ChunkManager struct {
	Player   *entity.Player
	Renderer *renderer.Renderer //TODO: interface this?

	chunks ChunkContainer
}

func NewChunkManager(r *renderer.Renderer, p *entity.Player) *ChunkManager {
	cm := &ChunkManager{
		Renderer: r,
		Player:   p,
	}

	cm.Setup()

	return cm
}

func (cm *ChunkManager) Setup() {
	cm.chunks = make(map[string]*Chunk)

	cm.newChunk(0, 0, 0)
	go cm.watchChunkLists()
}

func (cm *ChunkManager) watchChunkLists() {
	t := time.NewTicker(time.Millisecond * 100)
	defer t.Stop()

	drawDistance := float64(4 * ChunkSize)
	for range t.C {
		for _, ch := range cm.chunks {
			if ch.NumNeighbours == 6 {
				continue
			}

			{
				xPos := ch.GridPos.X() * ChunkSize * BlockSize
				yPos := ch.GridPos.Y() * ChunkSize * BlockSize
				zPos := ch.GridPos.Z() * ChunkSize * BlockSize
				lenHlf := float64(ChunkSize * BlockRenderSize)
				cntr := mgl64.Vec3{xPos, yPos, zPos}.Add(mgl64.Vec3{lenHlf, lenHlf, lenHlf})
				distVec := cntr.Sub(mgl64.Vec3{float64(cm.Player.Pos().X()), float64(cm.Player.Pos().Y()), float64(cm.Player.Pos().Z())})

				if distVec.Len() > drawDistance {
					cm.unloadChunk(ch)
					continue
				}
			}

			if ch.XPlus == nil {
				xPos := (ch.GridPos.X() + 1) * ChunkSize * BlockSize
				yPos := ch.GridPos.Y() * ChunkSize * BlockSize
				zPos := ch.GridPos.Z() * ChunkSize * BlockSize
				lenHlf := float64(ChunkSize * BlockRenderSize)
				cntr := mgl64.Vec3{xPos, yPos, zPos}.Add(mgl64.Vec3{lenHlf, lenHlf, lenHlf})
				distVec := cntr.Sub(mgl64.Vec3{float64(cm.Player.Pos().X()), float64(cm.Player.Pos().Y()), float64(cm.Player.Pos().Z())})

				if distVec.Len() <= drawDistance {
					nch, ok := cm.chunks.Lookup(int(ch.GridPos.X()+1), int(ch.GridPos.Y()), int(ch.GridPos.Z()))
					if !ok {
						nch = cm.newChunk(ch.GridPos.X()+1, ch.GridPos.Y(), ch.GridPos.Z())
					}

					ch.XPlus = nch
					nch.XMinus = ch
					ch.NumNeighbours += 1
					nch.NumNeighbours += 1
					// TODO: When we update the neighbours we need to update the meshes.
					// TODO: We will need to clear the current mesh (possibly call teadown but i dont know what happens when deleting from the card)
					// TODO: Then we will need to gen a FRESH mesh
					// TODO: We can possible swap the current mesh to a cached mesh then create a new mesh then delete the old mesh and swap if we get artifacts in rendering
				}
			}

			if ch.XMinus == nil {
				xPos := (ch.GridPos.X() - 1) * ChunkSize * BlockSize
				yPos := ch.GridPos.Y() * ChunkSize * BlockSize
				zPos := ch.GridPos.Z() * ChunkSize * BlockSize
				lenHlf := float64(ChunkSize * BlockRenderSize)
				cntr := mgl64.Vec3{xPos, yPos, zPos}.Add(mgl64.Vec3{lenHlf, lenHlf, lenHlf})
				distVec := cntr.Sub(mgl64.Vec3{float64(cm.Player.Pos().X()), float64(cm.Player.Pos().Y()), float64(cm.Player.Pos().Z())})

				if distVec.Len() <= drawDistance {
					nch, ok := cm.chunks.Lookup(int(ch.GridPos.X()-1), int(ch.GridPos.Y()), int(ch.GridPos.Z()))
					if !ok {
						nch = cm.newChunk(ch.GridPos.X()-1, ch.GridPos.Y(), ch.GridPos.Z())
					}

					ch.XMinus = nch
					nch.XPlus = ch
					ch.NumNeighbours += 1
					nch.NumNeighbours += 1
				}
			}

			if ch.ZPlus == nil {
				xPos := ch.GridPos.X() * ChunkSize * BlockSize
				yPos := ch.GridPos.Y() * ChunkSize * BlockSize
				zPos := (ch.GridPos.Z() + 1) * ChunkSize * BlockSize
				lenHlf := float64(ChunkSize * BlockRenderSize)
				cntr := mgl64.Vec3{xPos, yPos, zPos}.Add(mgl64.Vec3{lenHlf, lenHlf, lenHlf})
				distVec := cntr.Sub(mgl64.Vec3{float64(cm.Player.Pos().X()), float64(cm.Player.Pos().Y()), float64(cm.Player.Pos().Z())})

				if distVec.Len() <= drawDistance {
					nch, ok := cm.chunks.Lookup(int(ch.GridPos.X()), int(ch.GridPos.Y()), int(ch.GridPos.Z()+1))
					if !ok {
						nch = cm.newChunk(ch.GridPos.X(), ch.GridPos.Y(), ch.GridPos.Z()+1)
					}

					ch.ZPlus = nch
					nch.ZMinus = ch
					ch.NumNeighbours += 1
					nch.NumNeighbours += 1
				}
			}

			if ch.ZMinus == nil {
				xPos := ch.GridPos.X() * ChunkSize * BlockSize
				yPos := ch.GridPos.Y() * ChunkSize * BlockSize
				zPos := (ch.GridPos.Z() - 1) * ChunkSize * BlockSize
				lenHlf := float64(ChunkSize * BlockRenderSize)
				cntr := mgl64.Vec3{xPos, yPos, zPos}.Add(mgl64.Vec3{lenHlf, lenHlf, lenHlf})
				distVec := cntr.Sub(mgl64.Vec3{float64(cm.Player.Pos().X()), float64(cm.Player.Pos().Y()), float64(cm.Player.Pos().Z())})

				if distVec.Len() <= drawDistance {
					nch, ok := cm.chunks.Lookup(int(ch.GridPos.X()), int(ch.GridPos.Y()), int(ch.GridPos.Z()-1))
					if !ok {
						nch = cm.newChunk(ch.GridPos.X(), ch.GridPos.Y(), ch.GridPos.Z()-1)
					}

					ch.ZMinus = nch
					nch.ZPlus = ch
					ch.NumNeighbours += 1
					nch.NumNeighbours += 1
				}
			}

			if ch.YPlus == nil {
				xPos := ch.GridPos.X() * ChunkSize * BlockSize
				yPos := (ch.GridPos.Y() + 1) * ChunkSize * BlockSize
				zPos := ch.GridPos.Z() * ChunkSize * BlockSize
				lenHlf := float64(ChunkSize * BlockRenderSize)
				cntr := mgl64.Vec3{xPos, yPos, zPos}.Add(mgl64.Vec3{lenHlf, lenHlf, lenHlf})
				distVec := cntr.Sub(mgl64.Vec3{float64(cm.Player.Pos().X()), float64(cm.Player.Pos().Y()), float64(cm.Player.Pos().Z())})

				if distVec.Len() <= drawDistance {
					nch, ok := cm.chunks.Lookup(int(ch.GridPos.X()), int(ch.GridPos.Y()+1), int(ch.GridPos.Z()))
					if !ok {
						nch = cm.newChunk(ch.GridPos.X(), ch.GridPos.Y()+1, ch.GridPos.Z())
					}

					ch.YPlus = nch
					nch.YMinus = ch
					ch.NumNeighbours += 1
					nch.NumNeighbours += 1
				}
			}

			if ch.YMinus == nil {
				xPos := ch.GridPos.X() * ChunkSize * BlockSize
				yPos := (ch.GridPos.Y() - 1) * ChunkSize * BlockSize
				zPos := ch.GridPos.Z() * ChunkSize * BlockSize
				lenHlf := float64(ChunkSize * BlockRenderSize)
				cntr := mgl64.Vec3{xPos, yPos, zPos}.Add(mgl64.Vec3{lenHlf, lenHlf, lenHlf})
				distVec := cntr.Sub(mgl64.Vec3{float64(cm.Player.Pos().X()), float64(cm.Player.Pos().Y()), float64(cm.Player.Pos().Z())})

				if distVec.Len() <= drawDistance {
					nch, ok := cm.chunks.Lookup(int(ch.GridPos.X()), int(ch.GridPos.Y()-1), int(ch.GridPos.Z()))
					if !ok {
						nch = cm.newChunk(ch.GridPos.X(), ch.GridPos.Y()-1, ch.GridPos.Z())
					}

					ch.YMinus = nch
					nch.YPlus = ch
					ch.NumNeighbours += 1
					nch.NumNeighbours += 1
				}
			}
		}
	}
}

func (cm *ChunkManager) newChunk(x, y, z float64) *Chunk {
	xPos := x * ChunkSize * BlockSize
	yPos := y * ChunkSize * BlockSize
	zPos := z * ChunkSize * BlockSize

	ch := &Chunk{
		Mesh:    cm.Renderer.CreateMesh(), // TODO: can we re use this mesh by passing it to all other chunks?
		Pos:     mgl64.Vec3{xPos, yPos, zPos},
		GridPos: mgl64.Vec3{x, y, z},
	}

	cm.chunks.Set(int(x), int(y), int(z), ch)

	ch.Setup()
	ch.BuildMesh()
	ch.Mesh.Finish()

	return ch
}

func (cm *ChunkManager) unloadChunk(ch *Chunk) {
	if nch, ok := cm.chunks.Lookup(int(ch.GridPos.X()+1), int(ch.GridPos.Y()), int(ch.GridPos.Z())); ok {
		if nch.XMinus != nil {
			nch.NumNeighbours -= 1
			nch.XMinus = nil
		}
	}

	if nch, ok := cm.chunks.Lookup(int(ch.GridPos.X()-1), int(ch.GridPos.Y()), int(ch.GridPos.Z())); ok {
		if nch.XPlus != nil {
			nch.NumNeighbours -= 1
			nch.XPlus = nil
		}
	}

	if nch, ok := cm.chunks.Lookup(int(ch.GridPos.X()), int(ch.GridPos.Y()), int(ch.GridPos.Z()+1)); ok {
		if nch.ZMinus != nil {
			nch.NumNeighbours -= 1
			nch.ZMinus = nil
		}
	}

	if nch, ok := cm.chunks.Lookup(int(ch.GridPos.X()), int(ch.GridPos.Y()), int(ch.GridPos.Z()-1)); ok {
		if nch.ZPlus != nil {
			nch.NumNeighbours -= 1
			nch.ZPlus = nil
		}
	}

	if nch, ok := cm.chunks.Lookup(int(ch.GridPos.X()), int(ch.GridPos.Y()+1), int(ch.GridPos.Z())); ok {
		if nch.YMinus != nil {
			nch.NumNeighbours -= 1
			nch.YMinus = nil
		}
	}

	if nch, ok := cm.chunks.Lookup(int(ch.GridPos.X()), int(ch.GridPos.Y()-1), int(ch.GridPos.Z())); ok {
		if nch.YPlus != nil {
			nch.NumNeighbours -= 1
			nch.YPlus = nil
		}
	}

	ch.TearDown()
	cm.chunks.Unset(int(ch.GridPos.X()), int(ch.GridPos.Y()), int(ch.GridPos.Z()))
	ch = nil
}
