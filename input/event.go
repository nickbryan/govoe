package input

import "github.com/nickbryan/voxel/event"

const (
	KeyPressedEvent  = "win.KeyPressed"
	KeyReleasedEvent = "win.KeyReleased"
)

// Subscriber is the interface that wraps the Subscribe method.
//
// Subscribe subscribed the specified callback to the specified topics.
type Subscriber interface {
	Subscribe(cb event.Callback, topics ...event.Topic) int
}

// KeyEvent encapsulates the relevant information for a keyboard event. It should be dispatched when a WindowManager
// detects a keyboard event.
//
// Upon receiving a KeyEvent the Manager will trigger the relevant command callbacks that the user has registered.
type KeyEvent struct {
	Action   Action
	Key      Key
	Modifier ModifierKey
}
