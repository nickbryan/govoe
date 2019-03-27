package engine

import "testing"

func TestCreate(t *testing.T) {
	em := NewEnitytManager()

	e1 := em.Create()
	e2 := em.Create()
	e3 := em.Create()

	if e1.Id() != 1 {
		t.Errorf("Entity1 Id expected to be 1, but received: %v", e1.Id)
	}

	if e2.Id() != 2 {
		t.Errorf("Entity2 Id expected to be 2, but received: %v", e2.Id)
	}

	if e3.Id() != 3 {
		t.Errorf("Entity3 Id expected to be 3, but received: %v", e3.Id)
	}

	em.Destroy(e2)
	e4 := em.Create()
	if e4.Id() != 2 {
		t.Errorf("Entity4 expected to have filled slot 2 with Id expected to be 2, but received Id: %v", e4.Id)
	}
}

func TestAlive(t *testing.T) {
	em := NewEnitytManager()

	e1 := em.Create()
	e2 := em.Create()
	e3 := &Entity{id: 3}

	if em.Alive(e1) == false {
		t.Error("Entity1 was expected to be alive but alive returned false")
	}

	if em.Alive(e2) == false {
		t.Error("Entity2 was expected to be alive but alive returned false")
	}

	if em.Alive(e3) == true {
		t.Error("Entity3 was not expected to be alive but alive returned true")
	}
}

func TestDestroy(t *testing.T) {
	em := NewEnitytManager()

	e1 := em.Create()
	e2 := em.Create()

	em.Destroy(e2)

	if em.Alive(e1) == false {
		t.Error("Entity1 was expected to be alive but alive returned false")
	}

	if em.Alive(e2) == true {
		t.Error("Entity2 was not expected to be alive but alive returned true")
	}
}
