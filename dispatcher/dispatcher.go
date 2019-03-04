// Package encapsulates functionality for a publish/subscribe based message handler that is
// safe for concurrent use.
//
// Channels are used for synchronisation and to ensure events are delivered in the specified order.
package dispatcher

type action int

const (
	subscribe action = iota
	publish
	unsubscribe
	teardown
)

// command represents a single action and is used by the commandQueue to ensure actions
// get processed in the order they are sent.
type command struct {
	topics  []string
	ch      chan interface{}
	action  action
	message interface{}
}

// Dispatcher is a publish/subscribe based message handler that is safe for concurrent use.
type Dispatcher struct {
	subscriptions map[string]map[chan interface{}]struct{}
	channels      map[chan interface{}]map[string]struct{}
	commandQueue  chan command
	bufferSize    int
}

// New creates, initialises and starts a new Dispatcher.
func New(bufferSize int) *Dispatcher {
	d := &Dispatcher{
		subscriptions: make(map[string]map[chan interface{}]struct{}),
		channels:      make(map[chan interface{}]map[string]struct{}),
		commandQueue:  make(chan command),
		bufferSize:    bufferSize,
	}

	go d.processCommands()

	return d
}

// Subscribe returns a channel that all messages relating to the passed in topics will
// be sent on.
//
// If the channel is subscribed to multiple topics and a message is published on more
// than one of the subscribed topics then the channel will receive the same message
// multiple times.
func (d *Dispatcher) Subscribe(topics ...string) chan interface{} {
	ch := make(chan interface{}, d.bufferSize)

	d.commandQueue <- command{
		topics: topics,
		ch:     ch,
		action: subscribe,
	}

	return ch
}

// Unsubscribe will stop the specified channel from receiving any more messages on the
// specified topics. If there are no more topic subscriptions then the channel will be closed.
func (d *Dispatcher) Unsubscribe(ch chan interface{}, topics ...string) {
	d.commandQueue <- command{
		topics: topics,
		ch:     ch,
		action: unsubscribe,
	}
}

// Dispatch sends the specified message to all channels subscribed to the specified topics.
func (d *Dispatcher) Dispatch(msg interface{}, topics ...string) {
	d.commandQueue <- command{
		topics:  topics,
		action:  publish,
		message: msg,
	}
}

// Teardown exits the command processing loop and closes all remaining channels.
func (d *Dispatcher) Teardown() {
	d.commandQueue <- command{
		action: teardown,
	}
}

// processCommands is responsible for keeping all dispatcher actions in sync. The running
// goroutine and all channels will be closed when the Teardown method is called.
func (d *Dispatcher) processCommands() {
loop:
	for cmd := range d.commandQueue {
		switch cmd.action {
		case subscribe:
			for _, topic := range cmd.topics {
				if d.subscriptions[topic] == nil {
					d.subscriptions[topic] = make(map[chan interface{}]struct{})
				}

				if d.channels[cmd.ch] == nil {
					d.channels[cmd.ch] = make(map[string]struct{})
				}

				d.subscriptions[topic][cmd.ch] = struct{}{}
				d.channels[cmd.ch][topic] = struct{}{}
			}
		case unsubscribe:
			for _, topic := range cmd.topics {
				if _, ok := d.subscriptions[topic]; !ok {
					continue
				}

				if _, ok := d.subscriptions[topic][cmd.ch]; !ok {
					continue
				}

				delete(d.subscriptions[topic], cmd.ch)
				delete(d.channels[cmd.ch], topic)

				if len(d.subscriptions[topic]) == 0 {
					delete(d.subscriptions, topic)
				}

				if len(d.channels[cmd.ch]) == 0 {
					delete(d.channels, cmd.ch)
					close(cmd.ch)
				}
			}

		case publish:
			for _, topic := range cmd.topics {
				for ch := range d.subscriptions[topic] {
					ch <- cmd.message
				}
			}
		case teardown:
			for topic, channels := range d.subscriptions {
				for ch := range channels {
					delete(d.subscriptions[topic], ch)
				}

				delete(d.subscriptions, topic)
			}

			for ch := range d.channels {
				close(ch)
			}

			break loop
		}
	}
}
