package main

import (
	"fmt"
	"log"

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
