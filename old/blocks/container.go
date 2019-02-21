package blocks

import "fmt"

func idx(x, y, z int) int {
	//return x + y*ChunkSize + z*ChunkSizeSquared
	return x | y<<5 | z<<10
}

type BlockContainer []BlockType

func (c BlockContainer) Lookup(x, y, z int) BlockType {
	return c[idx(x, y, z)]
}

func (c BlockContainer) Set(x, y, z int, bt BlockType) {
	c[idx(x, y, z)] = bt
}

type LightMapContainer []byte

func (c LightMapContainer) Sunlight(x, y, z int) uint8 {
	return (c[idx(x, y, z)] >> 4) & 0xF
}

func (c LightMapContainer) SetSunlight(x, y, z int, val uint8) {
	i := idx(x, y, z)

	c[i] = (c[i] & 0xF) | (val << 4)
}

func (c LightMapContainer) Torchlight(x, y, z int) uint8 {
	return c[idx(x, y, z)] & 0xF
}

func (c LightMapContainer) SetTorchlight(x, y, z int, val uint8) {
	i := idx(x, y, z)

	c[i] = (c[i] & 0xF0) | val
}

// TODO: these should be locked via mutex or updated via channels
type ChunkContainer map[string]*Chunk

func (c ChunkContainer) Lookup(x, y, z int) (*Chunk, bool) {
	ch, ok := c[fmt.Sprintf("%d%d%d", x, y, z)]
	return ch, ok
}

func (c ChunkContainer) Set(x, y, z int, ch *Chunk) {
	c[fmt.Sprintf("%d%d%d", x, y, z)] = ch
}

func (c ChunkContainer) Unset(x, y, z int) {
	delete(c, fmt.Sprintf("%d%d%d", x, y, z))
}
