package input

import (
	"sync"

	"github.com/nickbryan/voxel/event"
)

type keyState int

const (
	press keyState = iota
	release
	off
)

// Manager is responsible for executing all input commands registered within the system.
//
// Manager will subscribe to the KeyPressedEvent and KeyReleasedEvent topics of the supplied event.Subscriber
// upon initialisation. When the Manager is notified of a KeyEvent it will call the KeyCommandExecutor registered
// with the specified key and action.
//
// All KeyCommandExecutor's are called within the Simulate method to ensure that they only get called once per simulation.
type Manager struct {
	initialised       bool
	keyEvents         chan interface{}
	keys              map[Key]keyState
	keysMux           sync.RWMutex
	keyPressListeners []keyPressListener
}

// New will create and Initialise a new Manager.
func New(subscriber event.Subscriber) *Manager {
	m := &Manager{
		keys: make(map[Key]keyState),
	}

	m.Initialise(subscriber)

	return m
}

// Initialise can only be called once per manager. It will subscribe and listen to the KeyPressedEvent and
// KeyReleasedEvent topics on the given event.Subscriber.
func (m *Manager) Initialise(subscriber event.Subscriber) {
	if m.initialised {
		return
	}

	m.keyEvents = subscriber.Subscribe(KeyPressedEvent, KeyReleasedEvent)

	go m.listen()

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

	m.keysMux.Lock()
	m.keys[key] = off
	m.keysMux.Unlock()

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
		m.keysMux.RLock()
		key := m.keys[l.key]
		m.keysMux.RUnlock()

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

			m.keysMux.Lock()
			m.keys[l.key] = off
			m.keysMux.Unlock()
		}

		if l.state == Release && key == release {
			for _, c := range l.commands {
				c.Execute(dt)
			}

			m.keysMux.Lock()
			m.keys[l.key] = off
			m.keysMux.Unlock()
		}
	}
}

func (m *Manager) listen() {
	for evt := range m.keyEvents {
		if e, ok := evt.(KeyEvent); ok {
			if e.Action == KeyPressed {
				m.keysMux.Lock()
				m.keys[e.Key] = press
				m.keysMux.Unlock()
			}

			if e.Action == KeyReleased {
				m.keysMux.Lock()
				m.keys[e.Key] = release
				m.keysMux.Unlock()
			}
		}
	}
}
