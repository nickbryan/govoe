package main

import (
	"fmt"
	"log"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/nickbryan/voxel/physics"

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

		em := engine.NewEnitytManager()
		p := em.Create()

		tm := &physics.TransformationManager{}
		t := &physics.TransformationComponent{}
		tm.Register(p, t)

		m := physics.NewMovementSystem(tm)
		m.Register(p, &physics.MovementComponent{
			Acceleration: mgl32.Vec3{9.18, 9.18, 9.18},
		})

		i := input.New(e.World.EventManager)
		i.AddKeyCommands(input.KeyW, input.Pressed, input.KeyCommandExecutorFunc(func(_ float64) {
			m.Move(p, physics.Forward)
		}))

		e.World.AddSimulator(i)
		e.World.AddSimulator(m)

		e.World.EventManager.Subscribe(func(_ event.Topic, msg interface{}) {
			if m, ok := msg.(engine.SimulationSecondElapsedMessage); ok {
				fmt.Println(fmt.Sprintf("FPS: %v, SPS: %v", m.Fps, m.Sps))
				fmt.Println(t)
			}
		}, engine.SimulationSecondElapsedEvent)

		e.Run()
	})
}
