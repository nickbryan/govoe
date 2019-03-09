package engine

type Simulator interface {
	Simulate(dt float64)
}

type World struct {
	EventManager EventManager

	simulators []Simulator
}

func NewWorld(eMgr EventManager) *World {
	return &World{
		EventManager: eMgr,
	}
}

func (w *World) AddSimulator(s Simulator) {
	w.simulators = append(w.simulators, s)
}

func (w *World) RunSimulations(dt float64) {
	for _, s := range w.simulators {
		s.Simulate(dt)
	}
}
