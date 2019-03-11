package engine

// SimulationSecondElapsedEvent will be published when a seconds worth of simulations have passed.
const SimulationSecondElapsedEvent = "simulation.second_elapsed"

// SimulationSecondElapsedMessage will be published on the SimulationSecondElapsedEvent topic. It contains information
// about the previous seconds worth of simulations.
type SimulationSecondElapsedMessage struct {
	Fps int
	Sps int
}

// SimulationStepper is responsible for handling a single step within the game simulation.
type SimulationStepper interface {
	Step(e *Engine)
}

// fixedStepSimulation allows world rendering to happen as fast as possible (when vSync is disabled) but limits the
// number of simulations to the specified limit (sps).
type fixedStepSimulation struct {
	initialised                         bool
	previousTime, frameTime, frameStart float64
	dt, accumulator                     float64
	sps, simulations, frames            int
}

// Step moves the simulation forward by one frame. If the simulation is running behind then multiple simulation frames
// may occur before rendering and ending the step.
func (sm *fixedStepSimulation) Step(e *Engine) {
	currentTime := e.GetTime()

	if !sm.initialised {
		sm.dt = 1 / float64(sm.sps)
		sm.frameStart = currentTime

		sm.initialised = true
	}

	sm.frameTime = currentTime - sm.frameStart
	sm.accumulator += sm.frameTime

	if sm.accumulator > 0.25 {
		sm.accumulator = 0.25
	}

	for sm.accumulator >= sm.dt {
		e.winMgr.PollEvents()

		e.World.Simulate(sm.dt)

		sm.accumulator -= sm.dt
		sm.simulations++
	}

	alpha := sm.accumulator / sm.dt

	e.World.Render(alpha)

	e.win.SwapBuffers()

	if currentTime-sm.previousTime >= 1 {
		e.World.EventManager.Publish(
			SimulationSecondElapsedMessage{Fps: sm.frames, Sps: sm.simulations},
			SimulationSecondElapsedEvent,
		)

		sm.simulations = 0
		sm.frames = 0
		sm.previousTime = currentTime
	}

	sm.frames += 1
	sm.frameStart = currentTime
}
