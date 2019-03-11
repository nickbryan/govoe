package main

import (
	"fmt"
	"log"

	"github.com/nickbryan/voxel/event"

	"github.com/nickbryan/voxel/input"

	"github.com/faiface/mainthread"
	"github.com/nickbryan/voxel/engine"
)

func main() {
	mainthread.Run(func() {
		e, err := engine.New(&engine.Configuration{
			Title:  "Engine Demo",
			Width:  1440,
			Height: 900,
			Sps:    20,
		})

		if err != nil {
			log.Fatalln("failed to create engine")
		}

		inputMgr := input.New(e.World.EventManager)
		inputMgr.AddKeyCommands(input.KeyW, input.Press, input.KeyCommandExecutorFunc(func(_ float64) {
			fmt.Println("W Press")
		}))

		inputMgr.AddKeyCommands(input.KeyW, input.Release, input.KeyCommandExecutorFunc(func(_ float64) {
			fmt.Println("W Release")
		}))

		inputMgr.AddKeyCommands(input.KeyS, input.Pressed, input.KeyCommandExecutorFunc(func(_ float64) {
			fmt.Println(".")
		}))

		e.World.AddSimulator(inputMgr)

		e.World.EventManager.Subscribe(func(_ event.Topic, msg interface{}) {
			if m, ok := msg.(engine.SimulationSecondElapsedMessage); ok {
				fmt.Println(fmt.Sprintf("FPS: %v, SPS: %v", m.Fps, m.Sps))
			}
		}, engine.SimulationSecondElapsedEvent)

		e.Run()
	})
}
