// Package engine is the heart of the system and is responsible for initialisation and orchestration of
// all sub-systems and managers.
package engine

import (
	"errors"
)

// Configuration should be passed into engine.New(c *Configuration) and will be used to initialise the engine instance.
type Configuration struct {
	Title         string        // Title is the window title.
	Width, Height int           // Width and Height will determine the initial dimensions of the window when not in fullscreen mode.
	Fps           int           // Fps is the desired frames per second. // TODO: is this needed and should it be merged with UPS? or do we just need Ups?
	Ups           int           // Ups is the desired updates per second.
	WindowManager WindowManager // WindowManager will be used to create the Window instance.
	EventManager  EventManager  // EventManager will be used by the engine to communicate between systems.
}

// engine allows us to ensure that Engine will only created once.
// One of the rare cases a singleton makes sense as we do not allow for multiple
// context in a number of instances throughout the engine.
var engine *Engine

// Engine is the main application instance. It will be create only once per application and is enforced as a singleton.
type Engine struct {
	running  bool
	closed   bool
	winMgr   WindowManager
	win      Window
	eventMgr EventManager
}

// New will create, configure and return a new Engine instance. It can only be called once per application.
func New(c *Configuration) (*Engine, error) {
	if engine != nil {
		return nil, errors.New("an instance of Engine has already been created")
	}

	if c.WindowManager == nil {
		c.WindowManager = &GLFWWindowManager{}
	}

	if c.EventManager == nil {
		c.EventManager = NewEventDispatcher(0)
	}

	e := &Engine{
		winMgr:   c.WindowManager,
		eventMgr: c.EventManager,
	}
	engine = e

	err := e.winMgr.Initialise()
	if err != nil {
		return nil, err
	}

	win, err := e.winMgr.CreateWindow(c.Width, c.Height, c.Title)
	if err != nil {
		return nil, err
	}
	e.win = win

	return e, nil
}

// Run starts the main loop and will continue to block until the window is ready to be closed.
func (e *Engine) Run() {
	// Prevent the engine from being run on multiple go routines.
	if e.running {
		return
	}
	e.running = true

	defer e.teardown()

	for !e.closed {
		e.winMgr.PollEvents()
		e.win.SwapBuffers()
		e.closed = e.win.ShouldClose()
	}

}

// teardown will be called when the main loop has finished. Any state destructuring should be triggered from here.
func (e *Engine) teardown() {
	e.winMgr.Teardown()
}
