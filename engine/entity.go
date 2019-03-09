package engine

type Entity struct {
	Id uint64
}

type EntityManager struct {
	entities map[uint64]*Entity
}

func NewEnitytManager() *EntityManager {
	return &EntityManager{
		entities: make(map[uint64]*Entity),
	}
}

func (em *EntityManager) CreateEntity() *Entity {
	e := &Entity{}

	for em.IsAlive(e) {
		e.Id++
	}

	em.entities[e.Id] = e

	return e
}

func (em *EntityManager) IsAlive(e *Entity) bool {
	_, ok := em.entities[e.Id]
	return ok
}

func (em *EntityManager) Destroy(e *Entity) {
	delete(em.entities, e.Id)
}
