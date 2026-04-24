package portevents_test

import (
	"sync"
	"testing"

	"github.com/example/portwatch/internal/portevents"
)

func TestPublish_SpecificSubscriberReceivesEvent(t *testing.T) {
	bus := portevents.New()
	var got portevents.Event

	bus.Subscribe(portevents.EventPortOpened, func(e portevents.Event) {
		got = e
	})

	bus.Publish(portevents.Event{Host: "localhost", Port: 22, Type: portevents.EventPortOpened})

	if got.Port != 22 || got.Host != "localhost" {
		t.Fatalf("expected event for port 22 on localhost, got %+v", got)
	}
}

func TestPublish_WildcardSubscriberReceivesAllEvents(t *testing.T) {
	bus := portevents.New()
	var received []portevents.Event
	var mu sync.Mutex

	bus.Subscribe("", func(e portevents.Event) {
		mu.Lock()
		received = append(received, e)
		mu.Unlock()
	})

	bus.Publish(portevents.Event{Port: 80, Type: portevents.EventPortOpened})
	bus.Publish(portevents.Event{Port: 443, Type: portevents.EventPortClosed})

	if len(received) != 2 {
		t.Fatalf("expected 2 events, got %d", len(received))
	}
}

func TestPublish_NoMatchingSubscriber_NoPanic(t *testing.T) {
	bus := portevents.New()
	// Should not panic with no subscribers.
	bus.Publish(portevents.Event{Port: 8080, Type: portevents.EventPortChanged})
}

func TestPublish_SpecificDoesNotReceiveOtherTypes(t *testing.T) {
	bus := portevents.New()
	called := false

	bus.Subscribe(portevents.EventPortOpened, func(e portevents.Event) {
		called = true
	})

	bus.Publish(portevents.Event{Port: 22, Type: portevents.EventPortClosed})

	if called {
		t.Fatal("opened subscriber should not receive closed event")
	}
}

func TestSubscriberCount(t *testing.T) {
	bus := portevents.New()

	if bus.SubscriberCount(portevents.EventPortOpened) != 0 {
		t.Fatal("expected 0 subscribers initially")
	}

	bus.Subscribe(portevents.EventPortOpened, func(e portevents.Event) {})
	bus.Subscribe(portevents.EventPortOpened, func(e portevents.Event) {})

	if bus.SubscriberCount(portevents.EventPortOpened) != 2 {
		t.Fatalf("expected 2 subscribers, got %d", bus.SubscriberCount(portevents.EventPortOpened))
	}
}
