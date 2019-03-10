package input

import (
	"fmt"
	"testing"

	"github.com/nickbryan/voxel/event"
)

type subscriber struct {
	ch chan interface{}
}

func (s *subscriber) Subscribe(...event.Topic) chan interface{} {
	return s.ch
}

func TestSimulate(t *testing.T) {
	var msgs []string

	s := &subscriber{ch: make(chan interface{}, 2)}
	i := New(s)

	i.AddKeyCommands(KeyA, Press, KeyCommandExecutorFunc(func(dt float64) {
		msgs = append(msgs, fmt.Sprintf("KeyA Press %v", dt))
	}))
	i.AddKeyCommands(KeyA, Release, KeyCommandExecutorFunc(func(dt float64) {
		msgs = append(msgs, fmt.Sprintf("KeyA Release %v", dt))
	}))

	i.AddKeyCommands(KeyB, Pressed, KeyCommandExecutorFunc(func(dt float64) {
		msgs = append(msgs, fmt.Sprintf("KeyB Pressed %v", dt))
	}))

	s.ch <- KeyEvent{
		Action: KeyPressed,
		Key:    KeyA,
	}
	// Could take multiple simulations for event to be processed.
	for len(msgs) < 1 {
		i.Simulate(1)
	}
	if msgs[0] != "KeyA Press 1" {
		t.Errorf("Simulation1 expected message to be: KeyA Press 1, received: %v", msgs[0])
	}

	s.ch <- KeyEvent{
		Action: KeyReleased,
		Key:    KeyA,
	}
	// Could take multiple simulations for event to be processed.
	for len(msgs) < 2 {
		i.Simulate(1)
	}
	if msgs[1] != "KeyA Release 1" {
		t.Errorf("Simulation2 expected message to be: KeyA Release 1, received: %v", msgs[1])
	}

	s.ch <- KeyEvent{
		Action: KeyPressed,
		Key:    KeyB,
	}
	// Could take multiple simulations for event to be processed.
	for len(msgs) < 4 {
		i.Simulate(1)
	}
	if msgs[2] != "KeyB Pressed 1" {
		t.Errorf("Simulation3 expected message to be:KeyB Pressed 1, received: %v", msgs[2])
	}
	if msgs[3] != "KeyB Pressed 1" {
		t.Errorf("Simulation4 expected message to be:KeyB Pressed 1, received: %v", msgs[3])
	}
	s.ch <- KeyEvent{
		Action: KeyReleased,
		Key:    KeyB,
	}
	if msgs[len(msgs)-1] != "KeyB Pressed 1" {
		t.Errorf("Simulation5 expected message to be:KeyB Pressed 1, received: %v", msgs[3])
	}
}
