package input

import "github.com/nickbryan/voxel/event"

const (
	KeyPressedEvent  = event.Topic("win.KeyPressed")
	KeyReleasedEvent = event.Topic("win.KeyReleased")
)

// KeyEvent encapsulates the relevant information for a keyboard event. It should be dispatched when a WindowManager
// detects a keyboard event.
//
// Upon receiving a KeyEvent the Manager will trigger the relevant command callbacks that the user has registered.
type KeyEvent struct {
	Action   Action
	Key      Key
	Modifier ModifierKey
}
