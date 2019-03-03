package engine

import (
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"

	"github.com/faiface/mainthread"

	"github.com/go-gl/glfw/v3.2/glfw"
)

// Window encapsulates both a top-level window and an OpenGL context.
type Window interface {
	SwapBuffers()
	ShouldClose() bool
}

// GLFWWindow wraps the GLFWwindow functionality to satisfy the Window interface.
type GLFWWindow struct {
	win *glfw.Window
}

// SwapBuffers swaps the front and back buffers of the specified window when rendering with OpenGL.
// If the swap interval is greater than zero, the GPU driver waits the specified number of screen
// updates before swapping the buffers.
func (w *GLFWWindow) SwapBuffers() {
	mainthread.Call(func() {
		w.win.SwapBuffers()
	})
}

// ShouldClose returns the value of the close flag of the specified window.
func (w *GLFWWindow) ShouldClose() bool {
	var shouldClose bool

	mainthread.Call(func() {
		shouldClose = w.win.ShouldClose()
	})

	return shouldClose
}

// WindowManager encapsulates all shared window functionality.
type WindowManager interface {
	Initialise() error
	Teardown()
	CreateWindow(width, height int, title string) (*GLFWWindow, error)
	PollEvents()
	SetSwapInterval(interval int)
}

// GLFWWindowManager wraps the shared GLFWWindow functionality.
type GLFWWindowManager struct {
}

// Initialise initialises GLFW and sets appropriate window hints.
func (wm *GLFWWindowManager) Initialise() error {
	return mainthread.CallErr(func() error {
		if err := glfw.Init(); err != nil {
			return err
		}

		glfw.WindowHint(glfw.ContextVersionMajor, 4)
		glfw.WindowHint(glfw.ContextVersionMinor, 1)
		glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
		glfw.WindowHint(glfw.Samples, 8)

		// Required for OSX.
		if runtime.GOOS == "darwin" {
			glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
		}

		return nil
	})
}

// Teardown will destroy any remaining window, monitor and cursor objects,
// restore any modified gamma ramps, re-enable the screensaver if it had
// been disabled and free any other resources allocated by GLFW.
func (wm *GLFWWindowManager) Teardown() {
	mainthread.Call(func() {
		glfw.Terminate()
	})
}

// CreateWindow creates a window, its associated OpenGL context and initialises the GLFW callbacks.
func (wm *GLFWWindowManager) CreateWindow(width, height int, title string) (*GLFWWindow, error) {
	var (
		err     error
		glfwWin *glfw.Window
		win     *GLFWWindow
	)

	mainthread.Call(func() {
		glfwWin, err = glfw.CreateWindow(width, height, title, nil, nil)
		if err != nil {
			return
		}

		win = &GLFWWindow{
			win: glfwWin,
		}

		win.win.MakeContextCurrent()

		win.win.SetFramebufferSizeCallback(func(win *glfw.Window, width int, height int) {
			// TODO: move this to event to remove dependency
			gl.Viewport(0, 0, int32(width), int32(height))
		})

	})

	return win, err
}

// PollEvents processes only those events that are already in the event
// queue and then returns immediately. Processing events will cause the
// window and input callbacks associated with those events to be called.
func (wm *GLFWWindowManager) PollEvents() {
	mainthread.Call(func() {
		glfw.PollEvents()
	})
}

// SetSwapInterval sets the swap interval for the current OpenGL context
// i.e. the number of screen updates to wait from the time glfwSwapBuffers
// was called before swapping the buffers and returning.
func (wm *GLFWWindowManager) SetSwapInterval(interval int) {
	mainthread.Call(func() {
		glfw.SwapInterval(interval)
	})
}
