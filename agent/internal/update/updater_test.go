package update

import (
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/yyhuni/lunafox/agent/internal/domain"
)

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

func TestUpdateOnceDockerUnavailable(t *testing.T) {
	updater := &Updater{}
	payload := domain.UpdateRequiredPayload{Version: "v1.0.0", Image: "yyhuni/lunafox-agent"}

	err := updater.updateOnce(payload)
	if err == nil {
		t.Fatalf("expected error when docker client is nil")
	}
	if !strings.Contains(err.Error(), "docker client unavailable") {
		t.Fatalf("unexpected error: %v", err)
	}
}
