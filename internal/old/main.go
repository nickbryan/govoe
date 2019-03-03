package old

import (
	"log"
	"net/http"

	"github.com/nickbryan/voxel/internal/old/engine"

	_ "net/http/pprof"

	"github.com/faiface/mainthread"
)

const (
	WindowWidth  uint = 1440
	WindowHeight uint = 900
)

func run() {
	e := engine.New(WindowWidth, WindowHeight)
	e.Run()
}

func main() {
	go func() {
		log.Fatal(http.ListenAndServe(`localhost:4200`, nil))
	}()

	mainthread.Run(run)

	// TODO: think about scene management, will they handle input, where will rendering happen and how
	// TODO: decide on float32 vs float64 (all entitt class idealy would use 64 but the matrix only accepts f32*)

	// TODO: 1. Chunk Manager, load multiple chunks, mesh across chunks, dynamically load chunks as player moves
	// TODO: 2. Terrain generation basic but nice
	// TODO: 3. Finish lighting algorithm for sunlight (when does the loop/bfs queue stop for sl?, tl stops when expected ie when light hits 0 in all dir)
	// TODO: 4. Implement day night cycle, make it settable for debugging
	// TODO: 5. Allow user to place torch by pressing a key
	// TODO: 6. Tidy up code and document
	// TODO: 7. Push to github
	// TODO: 8. Allow for block picking
	// TODO: 9. Lock player to floor
	// TODO: 10. Allow for switch between flying and not.
	// TODO: 11. Mesh switching from culled to greedy using another thread to update mesh cache.
	// TODO: 12. Evaluate
}
