package broker

import "sync"

// Broker represents a simple broker
type Broker struct {
	mu     sync.RWMutex
	topics map[string]*Topic
}

// NewBroker returns a new Broker instance.
func NewBroker() *Broker {
	return &Broker{
		topics: make(map[string]*Topic),
	}
}

// AddTopic adds a new topic if it doesn't already exist.
func (b *Broker) AddTopic(topic string) *Topic {
	b.mu.Lock()
	defer b.mu.Unlock()

	// If exists, return it
	if t, ok := b.topics[topic]; ok {
		return t
	}

	// Otherwise create a new topic
	t := NewTopic()
	b.topics[topic] = t
	return t
}
