package event

type action int

const (
	subscribe action = iota
	publish
	unsubscribe
	teardown
)

// Topic should be a unique name that is used to reference a set of events
// published on a given channel.
type Topic string

// Publisher is the interface that wraps the Publish method.
//
// Publish sends the specified message to all channels subscribed to the specified topics.
type Publisher interface {
	Publish(msg interface{}, topics ...Topic)
}

// Subscriber is the interface that wraps the Subscribe method.
//
// Subscribe returns a channel that all messages relating to the passed in topics will
// be sent on.
//
// If the channel is subscribed to multiple topics and a message is published on more
// than one of the subscribed topics then the channel will receive the same message
// multiple times.
type Subscriber interface {
	Subscribe(topics ...Topic) chan interface{}
}

// Unsubscriber is the interface that wraps the Unsubscribe method.
//
// Unsubscribe will stop the specified channel from receiving any more messages on the
// specified topics. If there are no more topic subscriptions then the channel will be closed.
type Unsubscriber interface {
	Unsubscribe(ch chan interface{}, topics ...Topic)
}

// command represents a single action and is used by the commandQueue to ensure actions
// get processed in the order they are sent.
type command struct {
	topics  []Topic
	ch      chan interface{}
	action  action
	message interface{}
}

// Manager is a publish/subscribe based event manager that is safe for concurrent use.
type Manager struct {
	subscriptions map[Topic]map[chan interface{}]struct{}
	channels      map[chan interface{}]map[Topic]struct{}
	commandQueue  chan command
	bufferSize    int
}

// NewManager creates, initialises and starts a new Manager.
func NewManager(bufferSize int) *Manager {
	ed := &Manager{
		subscriptions: make(map[Topic]map[chan interface{}]struct{}),
		channels:      make(map[chan interface{}]map[Topic]struct{}),
		commandQueue:  make(chan command),
		bufferSize:    bufferSize,
	}

	go ed.processCommands()

	return ed
}

// Publish sends the specified message to all channels subscribed to the specified topics.
func (ed *Manager) Publish(msg interface{}, topics ...Topic) {
	ed.commandQueue <- command{
		topics:  topics,
		action:  publish,
		message: msg,
	}
}

// Subscribe returns a channel that all messages relating to the passed in topics will
// be sent on.
//
// If the channel is subscribed to multiple topics and a message is published on more
// than one of the subscribed topics then the channel will receive the same message
// multiple times.
func (ed *Manager) Subscribe(topics ...Topic) chan interface{} {
	ch := make(chan interface{}, ed.bufferSize)

	ed.commandQueue <- command{
		topics: topics,
		ch:     ch,
		action: subscribe,
	}

	return ch
}

// Unsubscribe will stop the specified channel from receiving any more messages on the
// specified topics. If there are no more topic subscriptions then the channel will be closed.
func (ed *Manager) Unsubscribe(ch chan interface{}, topics ...Topic) {
	ed.commandQueue <- command{
		topics: topics,
		ch:     ch,
		action: unsubscribe,
	}
}

// Teardown exits the command processing loop and closes all remaining channels.
func (ed *Manager) Teardown() {
	ed.commandQueue <- command{
		action: teardown,
	}
}

// processCommands is responsible for keeping all dispatcher actions in sync. The running
// goroutine and all channels will be closed when the Teardown method is called.
func (ed *Manager) processCommands() {
loop:
	for cmd := range ed.commandQueue {
		switch cmd.action {
		case subscribe:
			for _, topic := range cmd.topics {
				if ed.subscriptions[topic] == nil {
					ed.subscriptions[topic] = make(map[chan interface{}]struct{})
				}

				if ed.channels[cmd.ch] == nil {
					ed.channels[cmd.ch] = make(map[Topic]struct{})
				}

				ed.subscriptions[topic][cmd.ch] = struct{}{}
				ed.channels[cmd.ch][topic] = struct{}{}
			}
		case unsubscribe:
			for _, topic := range cmd.topics {
				if _, ok := ed.subscriptions[topic]; !ok {
					continue
				}

				if _, ok := ed.subscriptions[topic][cmd.ch]; !ok {
					continue
				}

				delete(ed.subscriptions[topic], cmd.ch)
				delete(ed.channels[cmd.ch], topic)

				if len(ed.subscriptions[topic]) == 0 {
					delete(ed.subscriptions, topic)
				}

				if len(ed.channels[cmd.ch]) == 0 {
					delete(ed.channels, cmd.ch)
					close(cmd.ch)
				}
			}

		case publish:
			for _, topic := range cmd.topics {
				for ch := range ed.subscriptions[topic] {
					ch <- cmd.message
				}
			}
		case teardown:
			for topic, channels := range ed.subscriptions {
				for ch := range channels {
					delete(ed.subscriptions[topic], ch)
				}

				delete(ed.subscriptions, topic)
			}

			for ch := range ed.channels {
				close(ch)
			}

			break loop
		}
	}
}
