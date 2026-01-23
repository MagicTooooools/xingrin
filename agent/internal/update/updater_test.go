package update

import (
	"math/rand"
	"testing"
	"time"
)

func TestSanitizeContainerName(t *testing.T) {
	got := sanitizeContainerName("v1.0.0+TEST")
	if got == "" {
		t.Fatalf("expected sanitized name")
	}
	if got == "v1.0.0+test" {
		t.Fatalf("expected sanitized to replace invalid chars")
	}
}

func TestWithJitterRange(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	delay := 10 * time.Second
	got := withJitter(delay, rng)
	if got < delay {
		t.Fatalf("expected jitter >= delay")
	}
	if got > delay+(delay/5) {
		t.Fatalf("expected jitter <= 20%%")
	}
}
