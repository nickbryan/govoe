package engine

import (
	"fmt"
	"log"
	"sync"

	"github.com/nickbryan/voxel/internal/old/blocks"

	"github.com/nickbryan/voxel/internal/oldernal/old/entity"

	"github.com/nickbryan/voxel/internal/oldernal/old/input"

	"github.com/nickbryan/voxel/internal/oldernal/old/renderer"

	"github.com/faiface/mainthread"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Engine struct {
	win                 *glfw.Window
	WinWidth, WinHeight uint
	closed              bool
	running             bool

	resizeOnce sync.Once

	renderer     *renderer.Renderer
	camera       *entity.Camera
	inputManager *input.Input
	player       *entity.Player

	chunkManager *blocks.ChunkManager
}

func New(winWidth, winHeight uint) *Engine {
	return &Engine{
		WinWidth:  winWidth,
		WinHeight: winHeight,
	}
}

func (e *Engine) Run() {
	if e.running {
		return
	}

	e.setup()
	e.running = true
	defer e.tearDown()

	previousTime := 0.0
	updates := 0
	frames := 0

	fps := 20.0
	dt := 1 / fps
	accumulator := 0.0
	frameTime := 0.0

	gameStart := glfw.GetTime()
	frameStart := gameStart

	for !e.closed {
		currentTime := glfw.GetTime()

		frameTime = currentTime - frameStart
		accumulator += frameTime

		if accumulator > 0.25 {
			accumulator = 0.25
		}

		for accumulator >= dt {
			e.update(dt)

			accumulator -= dt
			updates++
		}

		alpha := accumulator / dt

		e.render(alpha)

		if currentTime-previousTime >= 1 {
			mainthread.Call(func() {
				e.win.SetTitle(fmt.Sprintf("Fps: %d UPS: %d", frames, updates))
			})

			updates = 0
			frames = 0
			previousTime = currentTime
		}

		frames += 1
		frameStart = currentTime
	}
}

func (e *Engine) setup() {
	mainthread.Call(func() {
		if err := glfw.Init(); err != nil {
			log.Fatalln("Failed to initialize glfw: ", err)
		}

		glfw.WindowHint(glfw.ContextVersionMajor, 4)
		glfw.WindowHint(glfw.ContextVersionMinor, 1)
		glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
		glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

		win, err := glfw.CreateWindow(int(e.WinWidth), int(e.WinHeight), "Voxel", nil, nil)
		if err != nil {
			panic(err)
		}
		e.win = win

		e.win.MakeContextCurrent()

		e.win.SetFramebufferSizeCallback(func(win *glfw.Window, width int, height int) {
			gl.Viewport(0, 0, int32(width), int32(height))
		})

		if err := gl.Init(); err != nil {
			panic(err)
		}

		version := gl.GoStr(gl.GetString(gl.VERSION))
		fmt.Println("OpenGL version: ", version)

		gl.Viewport(0, 0, int32(e.WinWidth), int32(e.WinHeight))

		gl.Enable(gl.DEPTH_TEST)

		e.renderer = renderer.New()

		e.player = entity.NewPlayer()

		e.camera = entity.NewCamera()
		e.camera.Attach(e.player)

		e.inputManager = input.New(e.win)
		e.inputManager.Register()
		e.inputManager.AddKeyCommands(glfw.KeyW, input.Pressed, input.KeyCommandFunc(func() {
			e.player.Stride(0.05)
			fmt.Println(e.player.Pos())
		}))
		e.inputManager.AddKeyCommands(glfw.KeyS, input.Pressed, input.KeyCommandFunc(func() {
			e.player.Stride(-0.05)
			fmt.Println(e.player.Pos())
		}))
		e.inputManager.AddKeyCommands(glfw.KeyA, input.Pressed, input.KeyCommandFunc(func() {
			e.player.Strafe(-0.05)
			fmt.Println(e.player.Pos())
		}))
		e.inputManager.AddKeyCommands(glfw.KeyD, input.Pressed, input.KeyCommandFunc(func() {
			e.player.Strafe(0.05)
			fmt.Println(e.player.Pos())
		}))
		e.inputManager.AddKeyCommands(glfw.KeyQ, input.Pressed, input.KeyCommandFunc(func() {
			e.player.Climb(0.05)
			fmt.Println(e.player.Pos())
		}))
		e.inputManager.AddKeyCommands(glfw.KeyE, input.Pressed, input.KeyCommandFunc(func() {
			e.player.Climb(-0.05)
			fmt.Println(e.player.Pos())
		}))
		e.inputManager.AddMouseMoveCommands(input.MouseMoveCommandFunc(func(offsetX, offsetY float64) {
			e.player.Look(float32(offsetX), float32(offsetY))
		}))
	})

	e.chunkManager = blocks.NewChunkManager(e.renderer, e.player)

}

func (e *Engine) tearDown() {
	e.renderer.Teardown()

	mainthread.Call(func() {
		glfw.Terminate()
	})
}

func (e *Engine) update(dt float64) {
	e.inputManager.Update()
	e.camera.Update()
}

func (e *Engine) render(alpha float64) {
	mainthread.Call(func() {
		gl.ClearColor(0.57, 0.71, 0.77, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		e.renderer.Draw(e.camera)

		e.win.SwapBuffers()
		glfw.PollEvents()

		// TODO: replace with glfw3.3 fix when released
		e.resizeOnce.Do(func() {
			x, y := e.win.GetPos()
			e.win.SetPos(x+1, y)
		})

		e.closed = e.win.ShouldClose()
	})
}
