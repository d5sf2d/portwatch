package sampler

import (
	"testing"
	"time"
)

func defaultSampler() *Sampler {
	return New(Options{
		Min:      10 * time.Second,
		Max:      100 * time.Second,
		Initial:  100 * time.Second,
		StepDown: 0.5,
		StepUp:   2.0,
	})
}

func TestCurrent_InitialEqualsMax(t *testing.T) {
	s := defaultSampler()
	if s.Current() != 100*time.Second {
		t.Fatalf("expected 100s, got %v", s.Current())
	}
}

func TestRecordChange_HalvesInterval(t *testing.T) {
	s := defaultSampler()
	s.RecordChange()
	if s.Current() != 50*time.Second {
		t.Fatalf("expected 50s, got %v", s.Current())
	}
}

func TestRecordChange_DoesNotGoBelowMin(t *testing.T) {
	s := defaultSampler()
	for i := 0; i < 20; i++ {
		s.RecordChange()
	}
	if s.Current() < s.min {
		t.Fatalf("interval %v went below min %v", s.Current(), s.min)
	}
	if s.Current() != s.min {
		t.Fatalf("expected min %v, got %v", s.min, s.Current())
	}
}

func TestRecordStable_DoublesInterval(t *testing.T) {
	s := defaultSampler()
	s.RecordChange() // bring to 50s
	s.RecordStable() // should go to 100s
	if s.Current() != 100*time.Second {
		t.Fatalf("expected 100s, got %v", s.Current())
	}
}

func TestRecordStable_DoesNotExceedMax(t *testing.T) {
	s := defaultSampler()
	for i := 0; i < 10; i++ {
		s.RecordStable()
	}
	if s.Current() > s.max {
		t.Fatalf("interval %v exceeded max %v", s.Current(), s.max)
	}
}

func TestReset_RestoresMax(t *testing.T) {
	s := defaultSampler()
	s.RecordChange()
	s.RecordChange()
	s.Reset()
	if s.Current() != s.max {
		t.Fatalf("expected max %v after reset, got %v", s.max, s.Current())
	}
}

func TestNew_DefaultsApplied(t *testing.T) {
	s := New(Options{})
	if s.min != defaultMinInterval {
		t.Fatalf("expected default min %v, got %v", defaultMinInterval, s.min)
	}
	if s.max != defaultMaxInterval {
		t.Fatalf("expected default max %v, got %v", defaultMaxInterval, s.max)
	}
	if s.Current() != defaultMaxInterval {
		t.Fatalf("expected initial == max, got %v", s.Current())
	}
}
