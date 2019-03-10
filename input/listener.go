package input

// State represents the current state of a key within the input system.
type State int

const (
	// Release states that the listener will fire once when the key/button is released.
	Release State = iota
	// Press states that the listener will fire once when the key/button is pressed.
	Press
	// Release states that the listener will fire each frame while the key/button is pressed.
	Pressed
)

// keyPressListener encapsulates information relating to a registered key within the input manager. The Manager
// uses this information to determine if the commands should be executed within the current simulation.
type keyPressListener struct {
	key      Key
	state    State
	commands []KeyCommandExecutor
}
