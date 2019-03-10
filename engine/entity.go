package engine

// Entity represents a single entity/world object managed by the engine. An Entity is identified by its unique Id.
//
// Due to this engines implementation of an Entity Component System, the Entity does not know about the components
// that are associated with it. Instead, an Entity should be registered with a component manager and it should be
// up to that component manager to keep track of the data relevant to the Entity within that system.
type Entity struct {
	Id uint64
}

// EntityManager is responsible for the creation and destruction of Entities within the engine. All Entities should
// be created by the EntityManager to ensure that they are managed correctly.
type EntityManager struct {
	entities map[uint64]*Entity
}

// NewEntityManager will create a new EntityManager.
func NewEnitytManager() *EntityManager {
	return &EntityManager{
		entities: make(map[uint64]*Entity),
	}
}

// Create will iterate over the existing Entities until it finds a slot in which an entity does not exist. It will
// then add the Entity to the entities list and return the newly created Entity.
//
// Entity Ids start at 1.
//
// Managing Entities this way allows us to ensure that there are no gaps in our Entity Id's and keeps our map sequentially
// stored in memory.
func (em *EntityManager) Create() *Entity {
	e := &Entity{Id: 1}

	for em.Alive(e) {
		e.Id++
	}

	em.entities[e.Id] = e

	return e
}

// Alive will check to see if there is an Entity mapped with the given Entities Id within the EntityManager.
func (em *EntityManager) Alive(e *Entity) bool {
	_, ok := em.entities[e.Id]
	return ok
}

// Destroy will remove the given Entity from the EntityManager.
func (em *EntityManager) Destroy(e *Entity) {
	delete(em.entities, e.Id)
}
