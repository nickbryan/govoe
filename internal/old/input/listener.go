package input

import "github.com/go-gl/glfw/v3.2/glfw"

const (
	// Release states that the listener will fire once when the key/button is released.
	Release State = iota
	// Press states that the listener will fire once when the key/button is pressed.
	Press
	// Release states that the listener will fire each frame while the key/button is pressed.
	Pressed
)

type keyPressListener struct {
	key      glfw.Key
	state    State
	commands []KeyCommand
}

type mouseButtonListener struct {
	button   glfw.MouseButton
	state    State
	commands []MouseButtonCommand
}
