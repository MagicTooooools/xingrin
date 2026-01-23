package websocket

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/yyhuni/orbit/agent/internal/health"
	"github.com/yyhuni/orbit/agent/internal/metrics"
)

func TestHeartbeatSenderSendOnce(t *testing.T) {
	client := &Client{send: make(chan []byte, 1)}
	collector := metrics.NewCollector()
	healthManager := health.NewManager()
	healthManager.Set("degraded", "cpu", "high load")

	sender := NewHeartbeatSender(client, collector, healthManager, "v9.9.9", "test-host", func() int { return 7 })
	sender.startedAt = time.Now().Add(-5 * time.Second)

	sender.sendOnce()

	select {
	case payload := <-client.send:
		var env struct {
			Type      string          `json:"type"`
			Payload   json.RawMessage `json:"payload"`
			Timestamp time.Time       `json:"timestamp"`
		}
		if err := json.Unmarshal(payload, &env); err != nil {
			t.Fatalf("unmarshal envelope failed: %v", err)
		}
		if env.Type != "heartbeat" {
			t.Fatalf("expected heartbeat type, got %q", env.Type)
		}
		if env.Timestamp.IsZero() {
			t.Fatalf("expected timestamp to be set")
		}

		var hb HeartbeatPayload
		if err := json.Unmarshal(env.Payload, &hb); err != nil {
			t.Fatalf("unmarshal payload failed: %v", err)
		}
		if hb.Version != "v9.9.9" {
			t.Fatalf("expected version v9.9.9, got %q", hb.Version)
		}
		if hb.Hostname != "test-host" {
			t.Fatalf("expected hostname test-host, got %q", hb.Hostname)
		}
		if hb.Tasks != 7 {
			t.Fatalf("expected tasks 7, got %d", hb.Tasks)
		}
		if hb.Uptime <= 0 {
			t.Fatalf("expected positive uptime, got %d", hb.Uptime)
		}
		if hb.Health.State != "degraded" {
			t.Fatalf("expected health state degraded, got %q", hb.Health.State)
		}
		if hb.Health.Reason != "cpu" {
			t.Fatalf("expected health reason cpu, got %q", hb.Health.Reason)
		}
		if hb.Health.Message != "high load" {
			t.Fatalf("expected health message, got %q", hb.Health.Message)
		}
		if hb.Health.Since == nil {
			t.Fatalf("expected health since to be set")
		}
	default:
		t.Fatalf("expected heartbeat payload to be sent")
	}
}
