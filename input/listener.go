package input

type State int

const (
	// Release states that the listener will fire once when the key/button is released.
	Release State = iota
	// Press states that the listener will fire once when the key/button is pressed.
	Press
	// Release states that the listener will fire each frame while the key/button is pressed.
	Pressed
)

type keyPressListener struct {
	key      Key
	state    State
	commands []KeyCommandExecutor
}
