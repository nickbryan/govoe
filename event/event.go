package event

// Topic should be a unique name that is used to reference a set of events
// published on a given channel.
type Topic string

// Callback is the function that will be called when a message is published to a topic.
type Callback func(t Topic, msg interface{})

// Manager is a publish/subscribe based event manager that is safe for concurrent use.
type Manager struct {
	ids           map[int]map[Topic]struct{}
	subscriptions map[Topic]map[int]Callback
}

// NewManager creates, initialises and starts a new Manager.
func NewManager() *Manager {
	return &Manager{
		ids:           make(map[int]map[Topic]struct{}),
		subscriptions: make(map[Topic]map[int]Callback),
	}
}

// Publish calls any callbacks that are subscribed to the specified topics and passes in the message.
func (m *Manager) Publish(msg interface{}, topics ...Topic) {
	for _, t := range topics {
		for _, cb := range m.subscriptions[t] {
			cb(t, msg)
		}
	}
}

// Subscribe registers the callback with the specified topics.
//
// A unique id is returned which can be used to unsubscribe the callback with the given id from
// specific topics.
func (m *Manager) Subscribe(cb Callback, topics ...Topic) int {
	id := 0
	for {
		if _, ok := m.ids[id]; !ok {
			break
		}

		id++
	}

	if m.ids[id] == nil {
		m.ids[id] = make(map[Topic]struct{})
	}

	for _, t := range topics {
		m.ids[id][t] = struct{}{}

		if m.subscriptions[t] == nil {
			m.subscriptions[t] = make(map[int]Callback)
		}

		m.subscriptions[t][id] = cb
	}

	return id
}

// Unsubscribe will unsubscribe the callback with the given id from the listed topics. When no more topics are associated
// with the given id, the id will be freed up to be reused on in subsequent Subscribe calls.
func (m *Manager) Unsubscribe(id int, topics ...Topic) {
	for _, t := range topics {
		if _, ok := m.subscriptions[t][id]; ok {
			delete(m.subscriptions[t], id)
		}

		if _, ok := m.ids[id][t]; ok {
			delete(m.ids[id], t)
		}

		if len(m.ids[id]) == 0 {
			delete(m.ids, id)
		}
	}
}
