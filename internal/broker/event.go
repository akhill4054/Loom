package broker

// Action represents the type of operation in the broker.
type Action byte

const (
	// Publish represents the action of publishing a event.
	Publish Action = 1
	// Subscribe represents the action of subscribing a event.
	Subscribe Action = 2
)

// Event represents a broker message with associated data.
type Event struct {
	Action  Action
	Topic   string
	Payload []byte
}
