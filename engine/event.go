package engine

type action int

const (
	subscribe action = iota
	publish
	unsubscribe
	teardown
)

// Publisher is the interface that wraps the Publish method.
//
// Publish sends the specified message to all channels subscribed to the specified topics.
type Publisher interface {
	Publish(msg interface{}, topics ...string)
}

// Subscriber is the interface that wraps the Publish method.
//
// Subscribe returns a channel that all messages relating to the passed in topics will
// be sent on.
//
// If the channel is subscribed to multiple topics and a message is published on more
// than one of the subscribed topics then the channel will receive the same message
// multiple times.
type Subscriber interface {
	Subscribe(topics ...string) chan interface{}
}

// Unsubscriber is the interface that wraps the Publish method.
//
// Unsubscribe will stop the specified channel from receiving any more messages on the
// specified topics. If there are no more topic subscriptions then the channel will be closed.
type Unsubscriber interface {
	Unsubscribe(ch chan interface{}, topics ...string)
}

// EventManager is a basic channel based publish/subscribe event system. It is used
// for communication between systems within the engine.
type EventManager interface {
	Publisher
	Subscriber
	Unsubscriber
	Teardown()
}

// command represents a single action and is used by the commandQueue to ensure actions
// get processed in the order they are sent.
type command struct {
	topics  []string
	ch      chan interface{}
	action  action
	message interface{}
}

// EventDispatcher is a publish/subscribe based EventManager that is safe for concurrent use.
type EventDispatcher struct {
	subscriptions map[string]map[chan interface{}]struct{}
	channels      map[chan interface{}]map[string]struct{}
	commandQueue  chan command
	bufferSize    int
}

// NewEventDispatcher creates, initialises and starts a new EventDispatcher.
func NewEventDispatcher(bufferSize int) *EventDispatcher {
	ed := &EventDispatcher{
		subscriptions: make(map[string]map[chan interface{}]struct{}),
		channels:      make(map[chan interface{}]map[string]struct{}),
		commandQueue:  make(chan command),
		bufferSize:    bufferSize,
	}

	go ed.processCommands()

	return ed
}

// Publish sends the specified message to all channels subscribed to the specified topics.
func (ed *EventDispatcher) Publish(msg interface{}, topics ...string) {
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
func (ed *EventDispatcher) Subscribe(topics ...string) chan interface{} {
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
func (ed *EventDispatcher) Unsubscribe(ch chan interface{}, topics ...string) {
	ed.commandQueue <- command{
		topics: topics,
		ch:     ch,
		action: unsubscribe,
	}
}

// Teardown exits the command processing loop and closes all remaining channels.
func (ed *EventDispatcher) Teardown() {
	ed.commandQueue <- command{
		action: teardown,
	}
}

// processCommands is responsible for keeping all dispatcher actions in sync. The running
// goroutine and all channels will be closed when the Teardown method is called.
func (ed *EventDispatcher) processCommands() {
loop:
	for cmd := range ed.commandQueue {
		switch cmd.action {
		case subscribe:
			for _, topic := range cmd.topics {
				if ed.subscriptions[topic] == nil {
					ed.subscriptions[topic] = make(map[chan interface{}]struct{})
				}

				if ed.channels[cmd.ch] == nil {
					ed.channels[cmd.ch] = make(map[string]struct{})
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
