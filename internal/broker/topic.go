package broker

import (
	"log"
	"sync"
)

// Topic represents a created topic.
type Topic struct {
	mu   sync.RWMutex
	subs map[chan *Event]struct{}
}

// NewTopic returns a Topic instance.
func NewTopic() *Topic {
	return &Topic{subs: make(map[chan *Event]struct{})}
}

// Publish publishes an event to it.
func (t *Topic) Publish(event *Event) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	log.Println("Message published to:", event.Topic)

	for ch := range t.subs {
		select {
		// Non-blocking send
		case ch <- event:
		default:
			// drop message if subscriber is too slow
		}
	}
}

// AddSubscriber subscribes to it.
func (t *Topic) AddSubscriber() chan *Event {
	ch := make(chan *Event, 16)

	t.mu.Lock()
	defer t.mu.Unlock()

	t.subs[ch] = struct{}{}
	return ch
}

// RemoveSubscriber removes a subscriber.
func (t *Topic) RemoveSubscriber(ch chan *Event) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if _, ok := t.subs[ch]; ok {
		delete(t.subs, ch)
		close(ch) // close channel to stop subscriber goroutine
	}
}
