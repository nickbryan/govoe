package event

import (
	"testing"
)

func TestSubscribe(t *testing.T) {
	em := NewManager(2)
	defer em.Teardown()

	ch1 := em.Subscribe("t1")
	ch2 := em.Subscribe("t1", "t2")

	em.Publish("m1", "t1")
	em.Publish("m2", "t1", "t2")

	if v := <-ch1; v != "m1" {
		t.Errorf("Ch1.1 incorrect channel value: expected %v, but received: %v", "m1", v)
	}
	if v := <-ch2; v != "m1" {
		t.Errorf("Ch2.1 incorrect channel value: expected %v, but received: %v", "m1", v)
	}
	if v := <-ch2; v != "m2" {
		t.Errorf("Ch2.2 incorrect channel value: expected %v, but received: %v", "m2", v)
	}
	if v := <-ch2; v != "m2" {
		t.Errorf("Ch2.3 incorrect channel value: expected %v, but received: %v", "m2", v)
	}
}

func TestUnsubscribe(t *testing.T) {
	em := NewManager(3)
	defer em.Teardown()

	ch1 := em.Subscribe("t1", "t2")
	ch2 := em.Subscribe("t1", "t2")

	em.Publish("m1", "t1")
	if v := <-ch1; v != "m1" {
		t.Errorf("Ch1.1 incorrect channel value: expected %v, but received: %v", "m1", v)
	}
	if v := <-ch2; v != "m1" {
		t.Errorf("Ch2.1 incorrect channel value: expected %v, but received: %v", "m1", v)
	}

	em.Unsubscribe(ch1, "t1")
	em.Publish("m2", "t1", "t2")
	if v := <-ch1; v != "m2" {
		t.Errorf("Ch1.2 incorrect channel value: expected %v, but received: %v", "m2", v)
	}
	if v := <-ch2; v != "m2" {
		t.Errorf("C2.2 incorrect channel value: expected %v, but received: %v", "m2", v)
	}
	if v := <-ch2; v != "m2" {
		t.Errorf("C2.3 incorrect channel value: expected %v, but received: %v", "m2", v)
	}

	em.Unsubscribe(ch1, "t2")
	em.Unsubscribe(ch2, "t1", "t2")
	em.Publish("m3", "t1", "t2")

	if _, ok := <-ch1; ok {
		t.Fatalf("Ch1.3 channel is still receiving values and should have been closed.")
	}
	if _, ok := <-ch2; ok {
		t.Fatalf("Ch2.4 channel is still receiving values and should have been closed.")
	}
}
