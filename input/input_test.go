package input

import (
	"fmt"
	"testing"

	"github.com/nickbryan/voxel/event"
)

type subscriber struct {
	cb event.Callback
}

func (s *subscriber) Subscribe(cb event.Callback, topics ...event.Topic) int {
	s.cb = cb
	return 1
}

func TestSimulate(t *testing.T) {
	var msgs []string

	s := &subscriber{}
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

	s.cb(KeyPressedEvent, KeyEventMessage{
		Action: KeyPressed,
		Key:    KeyA,
	})

	i.Simulate(1)
	if msgs[0] != "KeyA Press 1" {
		t.Errorf("Simulation1 expected message to be: KeyA Press 1, received: %v", msgs[0])
	}

	s.cb(KeyReleasedEvent, KeyEventMessage{
		Action: KeyReleased,
		Key:    KeyA,
	})
	i.Simulate(1)
	if msgs[1] != "KeyA Release 1" {
		t.Errorf("Simulation2 expected message to be: KeyA Release 1, received: %v", msgs[1])
	}

	s.cb(KeyPressedEvent, KeyEventMessage{
		Action: KeyPressed,
		Key:    KeyB,
	})
	for len(msgs) < 4 {
		i.Simulate(1)
	}
	if msgs[2] != "KeyB Pressed 1" {
		t.Errorf("Simulation3 expected message to be:KeyB Pressed 1, received: %v", msgs[2])
	}
	if msgs[3] != "KeyB Pressed 1" {
		t.Errorf("Simulation4 expected message to be:KeyB Pressed 1, received: %v", msgs[3])
	}

	s.cb(KeyReleasedEvent, KeyEventMessage{
		Action: KeyReleased,
		Key:    KeyB,
	})
	if len(msgs) > 4 {
		t.Errorf("Did not expect any more messages after release event: received: %v", len(msgs))
	}
}
