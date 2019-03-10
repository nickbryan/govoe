package event

import (
	"testing"
)

func TestSubscribe(t *testing.T) {
	em := NewManager()

	cb1 := false
	em.Subscribe(func(tp Topic, msg interface{}) {
		cb1 = true
		if msg.(string) != "m1" {
			t.Errorf("Subscription to t1 was called with incorrect message: expected m1, but received: %v", msg.(string))
		}
	}, "t1")

	var msgs1 []string
	em.Subscribe(func(tp Topic, msg interface{}) {
		msgs1 = append(msgs1, msg.(string))
	}, "t1", "t2")

	var msgs2 []string
	em.Subscribe(func(tp Topic, msg interface{}) {
		msgs2 = append(msgs2, msg.(string))
	}, "t3")

	em.Publish("m1", "t1")
	em.Publish("m2", "t2")
	em.Publish("m3", "t2", "t3")

	if cb1 == false {
		t.Error("cb1 was not called")
	}

	if msgs1[0] != "m1" {
		t.Errorf("Subscription to t1, t2 message 1 incorrect: expected: m1, but received: %v", msgs1[0])
	}
	if msgs1[1] != "m2" {
		t.Errorf("Subscription to t1, t2 message 2 incorrect: expected: m2, but received: %v", msgs1[1])
	}
	if msgs1[2] != "m3" {
		t.Errorf("Subscription to t1, t2 message 3 incorrect: expected: m3, but received: %v", msgs1[2])
	}

	if msgs2[0] != "m3" {
		t.Errorf("Subscription to t2, t3 message 1 incorrect: expected: m3, but received: %v", msgs2[0])
	}
}

func TestUnsubscribe(t *testing.T) {
	em := NewManager()

	{
		s1 := em.Subscribe(func(tp Topic, msg interface{}) {
		}, "t1")

		s2 := em.Subscribe(func(tp Topic, msg interface{}) {
		}, "t1")

		em.Unsubscribe(s1, "t1")

		s3 := em.Subscribe(func(tp Topic, msg interface{}) {
		}, "t1")

		if s1 != 0 {
			t.Errorf("s1 id is expected to be 0: received: %v", s1)
		}

		if s2 != 1 {
			t.Errorf("s1 id is expected to be 1: received: %v", s2)
		}

		if s3 != 0 {
			t.Errorf("s1 id is expected to be 0: received: %v", s3)
		}
	}

	{
		var msgs1 []string
		s1 := em.Subscribe(func(tp Topic, msg interface{}) {
			msgs1 = append(msgs1, msg.(string))
		}, "t1")

		var msgs2 []string
		em.Subscribe(func(tp Topic, msg interface{}) {
			msgs2 = append(msgs2, msg.(string))
		}, "t1")

		em.Publish("m1", "t1")
		em.Unsubscribe(s1, "t1")
		em.Publish("m2", "t1")

		if len(msgs1) != 1 {
			t.Error("s1 received more messages than expected.")
		}
		if msgs1[0] != "m1" {
			t.Errorf("Subscription to t1 message 1 incorrect: expected: m1, but received: %v", msgs1[0])
		}
		if msgs2[0] != "m1" {
			t.Errorf("Subscription to t1.2 message 1 incorrect: expected: m1, but received: %v", msgs2[0])
		}
		if msgs2[1] != "m2" {
			t.Errorf("Subscription to t1.2 message 2 incorrect: expected: m2, but received: %v", msgs2[1])
		}
	}

	{
		var msgs1 []string
		s1 := em.Subscribe(func(tp Topic, msg interface{}) {
			msgs1 = append(msgs1, msg.(string))
		}, "t1", "t2")

		em.Publish("m1", "t1", "t2")
		em.Unsubscribe(s1, "t1")
		em.Publish("m2", "t1", "t2")

		if len(msgs1) != 3 {
			t.Error("s1 received more messages than expected.")
		}
		if msgs1[0] != "m1" {
			t.Errorf("Subscription to t1, t2 message 1 incorrect: expected: m1, but received: %v", msgs1[0])
		}
		if msgs1[1] != "m1" {
			t.Errorf("Subscription to t1, t2 message 2 incorrect: expected: m1, but received: %v", msgs1[1])
		}
		if msgs1[2] != "m2" {
			t.Errorf("Subscription to t1 message 3 incorrect: expected: m2, but received: %v", msgs1[2])
		}
	}
}
