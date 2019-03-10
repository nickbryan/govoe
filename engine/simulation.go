package engine

import (
	"fmt"
)

// SimulationStepper is responsible for handling a single step within the game simulation.
type SimulationStepper interface {
	Step(e *Engine)
}

// fixedStepSimulation
type fixedStepSimulation struct {
	initialised                                          bool
	previousTime, dt, accumulator, frameTime, frameStart float64
	updates, frames, ups                                 int
}

func (sm *fixedStepSimulation) initialise(currentTime float64) {
	sm.ups = 20
	sm.dt = 1 / float64(sm.ups)
	sm.frameStart = currentTime

	sm.initialised = true
}

func (sm *fixedStepSimulation) Step(e *Engine) {
	currentTime := e.GetTime()

	if !sm.initialised {
		sm.initialise(currentTime)
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
		sm.updates++
	}

	alpha := sm.accumulator / sm.dt

	e.World.Render(alpha)

	e.win.SwapBuffers()

	if currentTime-sm.previousTime >= 1 {
		// TODO: event
		fmt.Printf("fps: %d ups: %d \n", sm.frames, sm.updates)

		sm.updates = 0
		sm.frames = 0
		sm.previousTime = currentTime
	}

	sm.frames += 1
	sm.frameStart = currentTime
}
