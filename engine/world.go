package engine

type Simulator interface {
	Simulate(dt float64)
}

type Renderer interface {
	Render(interpolation float64)
}

type World struct {
	EventManager EventManager

	simulators []Simulator
	renderers  []Renderer
}

func NewWorld(eMgr EventManager) *World {
	return &World{
		EventManager: eMgr,
	}
}

func (w *World) AddSimulator(s Simulator) {
	w.simulators = append(w.simulators, s)
}

func (w *World) AddRenderer(r Renderer) {
	w.renderers = append(w.renderers, r)
}

func (w *World) Simulate(dt float64) {
	for _, s := range w.simulators {
		s.Simulate(dt)
	}
}

func (w *World) Render(interpolation float64) {
	for _, r := range w.renderers {
		r.Render(interpolation)
	}
}
