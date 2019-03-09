package input

import (
	"github.com/nickbryan/voxel/event"
)

type keyState int

const (
	press keyState = iota
	release
	off
)

type Manager struct {
	initialised       bool
	keyEvents         chan interface{}
	keys              map[Key]keyState
	keyPressListeners []keyPressListener
}

func New(subscriber event.Subscriber) *Manager {
	m := &Manager{
		keys: make(map[Key]keyState),
	}

	m.Initialise(subscriber)

	return m
}

func (m *Manager) Initialise(subscriber event.Subscriber) {
	if m.initialised {
		return
	}

	m.keyEvents = subscriber.Subscribe(KeyPressedEvent, KeyReleasedEvent)

	go m.listen()

	m.initialised = true
}

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

func (m *Manager) Simulate(_ float64) {
	for _, l := range m.keyPressListeners {
		if m.keys[l.key] == off {
			continue
		}

		if l.state == Pressed && m.keys[l.key] == press {
			for _, c := range l.commands {
				c.Execute()
			}
		}

		if l.state == Press && m.keys[l.key] == press {
			for _, c := range l.commands {
				c.Execute()
			}

			m.keys[l.key] = off
		}

		if l.state == Release && m.keys[l.key] == release {
			for _, c := range l.commands {
				c.Execute()
			}

			m.keys[l.key] = off
		}
	}
}

func (m *Manager) listen() {
	for evt := range m.keyEvents {
		if e, ok := evt.(KeyEvent); ok {
			m.handleKeyEvent(e)
		}
	}
}

func (m *Manager) handleKeyEvent(e KeyEvent) {
	if e.Action == KeyPressed {
		m.keys[e.Key] = press
	}

	if e.Action == KeyReleased {
		m.keys[e.Key] = release
	}
}
