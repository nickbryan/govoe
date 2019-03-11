package input

import "github.com/nickbryan/voxel/event"

type keyState int

const (
	press keyState = iota
	release
	off
)

// Manager is responsible for executing all input commands registered within the system.
//
// Manager will subscribe to the KeyPressedEvent and KeyReleasedEvent topics of the supplied event.AsyncSubscriber
// upon initialisation. When the Manager is notified of a KeyEventMessage it will call the KeyCommandExecutor registered
// with the specified key and action.
//
// All KeyCommandExecutor's are called within the Simulate method to ensure that they only get called once per simulation.
type Manager struct {
	initialised       bool
	keys              map[Key]keyState
	keyPressListeners []keyPressListener
}

// New will create and Initialise a new Manager.
func New(subscriber Subscriber) *Manager {
	m := &Manager{
		keys: make(map[Key]keyState),
	}

	m.Initialise(subscriber)

	return m
}

// Initialise can only be called once per manager. It will subscribe and listen to the KeyPressedEvent and
// KeyReleasedEvent topics on the given Subscriber.
func (m *Manager) Initialise(subscriber Subscriber) {
	if m.initialised {
		return
	}

	subscriber.Subscribe(m.keyCallback, KeyPressedEvent, KeyReleasedEvent)

	m.initialised = true
}

// AddKeyCommand will register the given commands, key and state within the Manager. The commands will
// be called once per simulation if the Manager has been notified of the relevant key and state changes.
func (m *Manager) AddKeyCommands(key Key, state State, commands ...KeyCommandExecutor) {
	for _, l := range m.keyPressListeners {
		if l.key == key && l.state == state {
			l.commands = append(l.commands, commands...)
			return
		}
	}

	m.keys[key] = off

	m.keyPressListeners = append(m.keyPressListeners, keyPressListener{
		key:      key,
		state:    state,
		commands: commands,
	})
}

// Simulate is responsible for triggering the registered KeyCommandExecutor's and should be called once per update.
//
// The supplied dt (delta time) will be passed into all KeyCommandExecutor's.
func (m *Manager) Simulate(dt float64) {
	for _, l := range m.keyPressListeners {
		key := m.keys[l.key]

		if key == off {
			continue
		}

		if l.state == Pressed && key == press {
			for _, c := range l.commands {
				c.Execute(dt)
			}
		}

		if l.state == Press && key == press {
			for _, c := range l.commands {
				c.Execute(dt)
			}

			m.keys[l.key] = off
		}

		if l.state == Release && key == release {
			for _, c := range l.commands {
				c.Execute(dt)
			}

			m.keys[l.key] = off
		}
	}
}

func (m *Manager) keyCallback(_ event.Topic, msg interface{}) {
	if msg, ok := msg.(KeyEventMessage); ok {
		if msg.Action == KeyPressed {
			m.keys[msg.Key] = press
		}

		if msg.Action == KeyReleased {
			m.keys[msg.Key] = release
		}
	}
}
