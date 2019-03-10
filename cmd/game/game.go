package main

import (
	"fmt"
	"log"

	"github.com/nickbryan/voxel/input"

	"github.com/faiface/mainthread"
	"github.com/nickbryan/voxel/engine"
)

func main() {
	// engine.CreateWorld
	// World/Engine holds a list of Manager instances?
	// Manager has single simulate method.
	// Managers can be added by user at start up.
	// Scenes are involved some how for game logic?
	// A Manager (ComponentManager) is responsible for keep track of and running logic on its entities components. <Entitity, CompnentData>.doLogic()?
	// World/Engine main loop iterates over Managers and calls update(dt)
	// Separate renderer?
	// Event Manager for system communication
	// Manager is a interface to the underlying system or the system coulkd be in the manager if its small enough
	// Managers can look at other managers for components if needed (they would hold pointers to the manger in this case as they would be linked in the domain)
	// How do loops work? Do we split update and render at the world/engine level?
	// Word/engine has a list of renderers and simulators?

	mainthread.Run(func() {
		e, err := engine.New(&engine.Configuration{
			Title:  "Engine Demo",
			Width:  1440,
			Height: 900,
			Fps:    60,
			Ups:    20,
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

		e.Run()
	})
}
