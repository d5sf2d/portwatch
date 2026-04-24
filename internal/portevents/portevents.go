// Package portevents provides a typed event bus for port state change events,
// allowing decoupled components to subscribe and react to scan results.
package portevents

import "sync"

// EventType classifies the kind of port change that occurred.
type EventType string

const (
	EventPortOpened EventType = "opened"
	EventPortClosed EventType = "closed"
	EventPortChanged EventType = "changed"
)

// Event carries information about a single port state transition.
type Event struct {
	Host  string
	Port  int
	Type  EventType
	Meta  map[string]string
}

// Handler is a function that receives a port event.
type Handler func(e Event)

// Bus is a simple synchronous event bus for port events.
type Bus struct {
	mu       sync.RWMutex
	handlers map[EventType][]Handler
}

// New returns an initialised Bus.
func New() *Bus {
	return &Bus{
		handlers: make(map[EventType][]Handler),
	}
}

// Subscribe registers a handler for the given event type.
// Passing an empty EventType subscribes to all event types.
func (b *Bus) Subscribe(t EventType, h Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[t] = append(b.handlers[t], h)
}

// Publish dispatches e to all matching subscribers and to wildcard subscribers.
func (b *Bus) Publish(e Event) {
	b.mu.RLock()
	specific := append([]Handler(nil), b.handlers[e.Type]...)
	wildcard := append([]Handler(nil), b.handlers[""]...)
	b.mu.RUnlock()

	for _, h := range specific {
		h(e)
	}
	for _, h := range wildcard {
		h(e)
	}
}

// SubscriberCount returns the number of handlers registered for t.
func (b *Bus) SubscriberCount(t EventType) int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.handlers[t])
}
