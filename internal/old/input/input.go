package input

import (
	"sync"

	"github.com/go-gl/glfw/v3.2/glfw"
)

type State int

type Input struct {
	Win                    *glfw.Window
	keys                   map[glfw.Key]bool
	keyPressListeners      []keyPressListener
	mouseLastX, mouseLastY float64
	mouseInitOnce          sync.Once
	mouseMoveCommands      []MouseMoveCommand
	mouseScrollCommands    []MouseScrollCommand
	mouseButtons           map[glfw.MouseButton]bool
	mouseButtonListeners   []mouseButtonListener
}

func New(win *glfw.Window) *Input {
	return &Input{
		Win:               win,
		keys:              make(map[glfw.Key]bool),
		keyPressListeners: make([]keyPressListener, 0),
	}
}

func (i *Input) Update() {
	for _, l := range i.keyPressListeners {
		if i.keys[l.key] && l.state == Pressed {
			for _, c := range l.commands {
				c.Execute()
			}
		}
	}

	for _, l := range i.mouseButtonListeners {
		if i.mouseButtons[l.button] && l.state == Pressed {
			for _, c := range l.commands {
				c.Execute()
			}
		}
	}
}

func (i *Input) AddKeyCommands(key glfw.Key, state State, commands ...KeyCommand) {
	for _, l := range i.keyPressListeners {
		if l.key == key && l.state == state {
			l.commands = append(l.commands, commands...)
			return
		}
	}

	i.keyPressListeners = append(i.keyPressListeners, keyPressListener{
		key:      key,
		state:    state,
		commands: commands,
	})
}

func (i *Input) AddMouseMoveCommands(commands ...MouseMoveCommand) {
	i.mouseMoveCommands = append(i.mouseMoveCommands, commands...)
}

func (i *Input) AddMouseScrollCommands(commands ...MouseScrollCommand) {
	i.mouseScrollCommands = append(i.mouseScrollCommands, commands...)
}

func (i *Input) Register() {
	i.Win.SetKeyCallback(i.keyCallback)
	i.Win.SetCursorPosCallback(i.mouseMoveCallback)
	i.Win.SetMouseButtonCallback(i.mouseButtonCallback)
	i.Win.SetScrollCallback(i.mouseScrollCallback)

	i.Win.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
}

func (i *Input) keyCallback(_ *glfw.Window, key glfw.Key, _ int, action glfw.Action, _ glfw.ModifierKey) {
	if i.keys == nil {
		i.keys = make(map[glfw.Key]bool)
	}

	isPressed := action != glfw.Release
	for _, l := range i.keyPressListeners {
		if l.key == key && ((l.state == Press && !i.keys[key] && isPressed) || (l.state == Release && i.keys[key] && !isPressed)) {
			for _, c := range l.commands {
				c.Execute()
			}
		}
	}

	i.keys[key] = isPressed
}

func (i *Input) mouseMoveCallback(_ *glfw.Window, xPos float64, yPos float64) {
	i.mouseInitOnce.Do(func() {
		i.mouseLastX = xPos
		i.mouseLastY = yPos
	})

	xOffset := xPos - i.mouseLastX
	yOffset := i.mouseLastY - yPos

	for _, c := range i.mouseMoveCommands {
		c.Execute(xOffset, yOffset)
	}

	i.mouseLastX = xPos
	i.mouseLastY = yPos
}

func (i *Input) mouseButtonCallback(_ *glfw.Window, button glfw.MouseButton, action glfw.Action, _ glfw.ModifierKey) {
	if i.mouseButtons == nil {
		i.mouseButtons = make(map[glfw.MouseButton]bool)
	}

	isPressed := action != glfw.Release
	for _, l := range i.mouseButtonListeners {
		if l.button == button && ((l.state == Press && !i.mouseButtons[button] && isPressed) || (l.state == Release && i.mouseButtons[button] && !isPressed)) {
			for _, c := range l.commands {
				c.Execute()
			}
		}
	}

	i.mouseButtons[button] = isPressed
}

func (i *Input) mouseScrollCallback(_ *glfw.Window, xOffset float64, yOffset float64) {
	for _, c := range i.mouseScrollCommands {
		c.Execute(xOffset, yOffset)
	}
}
