package ratelimit_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/ratelimit"
)

func TestAllow_WithinLimit(t *testing.T) {
	l := ratelimit.New(3, time.Second)
	for i := 0; i < 3; i++ {
		if !l.Allow() {
			t.Fatalf("expected Allow()=true on call %d", i+1)
		}
	}
}

func TestAllow_ExceedsLimit(t *testing.T) {
	l := ratelimit.New(2, time.Second)
	l.Allow()
	l.Allow()
	if l.Allow() {
		t.Fatal("expected Allow()=false after limit exceeded")
	}
}

func TestAllow_ResetsAfterInterval(t *testing.T) {
	l := ratelimit.New(1, 50*time.Millisecond)
	if !l.Allow() {
		t.Fatal("first Allow() should succeed")
	}
	if l.Allow() {
		t.Fatal("second Allow() should be denied")
	}
	time.Sleep(60 * time.Millisecond)
	if !l.Allow() {
		t.Fatal("Allow() should succeed after interval reset")
	}
}

func TestRemaining(t *testing.T) {
	l := ratelimit.New(5, time.Second)
	l.Allow()
	l.Allow()
	if got := l.Remaining(); got != 3 {
		t.Fatalf("expected 3 remaining, got %d", got)
	}
}

func TestReset(t *testing.T) {
	l := ratelimit.New(2, time.Second)
	l.Allow()
	l.Allow()
	if l.Allow() {
		t.Fatal("should be denied before reset")
	}
	l.Reset()
	if !l.Allow() {
		t.Fatal("should be allowed after Reset()")
	}
}

func TestAllow_ZeroMax(t *testing.T) {
	l := ratelimit.New(0, time.Second)
	if l.Allow() {
		t.Fatal("expected Allow()=false with max=0")
	}
}
