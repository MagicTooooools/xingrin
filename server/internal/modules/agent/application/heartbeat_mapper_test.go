package application

import (
	"testing"
	"time"

	"github.com/yyhuni/lunafox/server/internal/agentproto"
)

func TestHeartbeatMapperToDomainUpdateWithHealth(t *testing.T) {
	mapper := newHeartbeatMapper()
	now := time.Date(2026, 2, 16, 8, 0, 0, 0, time.UTC)
	since := time.Date(2026, 2, 15, 8, 0, 0, 0, time.FixedZone("UTC+8", 8*60*60))
	payload := agentproto.HeartbeatPayload{
		Version:  "v1.2.3",
		Hostname: "node-1",
		Health: &agentproto.HealthStatus{
			State:   "degraded",
			Reason:  "disk",
			Message: "disk usage high",
			Since:   &since,
		},
	}

	update := mapper.toDomainUpdate(now, payload)
	if !update.LastHeartbeat.Equal(now) {
		t.Fatalf("expected last heartbeat to equal now")
	}
	if !update.HasHealth {
		t.Fatalf("expected has health true")
	}
	if update.HealthState != "degraded" {
		t.Fatalf("expected health state degraded, got %q", update.HealthState)
	}
	if update.HealthSince == nil {
		t.Fatalf("expected health since")
	}
	if update.HealthSince.Location() != time.UTC {
		t.Fatalf("expected health since in UTC")
	}
}

func TestHeartbeatMapperToCachePayloadWithoutHealth(t *testing.T) {
	mapper := newHeartbeatMapper()
	payload := agentproto.HeartbeatPayload{
		CPU:      0.2,
		Mem:      0.3,
		Disk:     0.4,
		Tasks:    2,
		Version:  "v1.2.3",
		Hostname: "node-2",
		Uptime:   123,
	}

	data := mapper.toCachePayload(payload)
	if data.Health != nil {
		t.Fatalf("expected nil health")
	}
	if data.Version != payload.Version {
		t.Fatalf("expected version %q, got %q", payload.Version, data.Version)
	}
	if data.Tasks != payload.Tasks {
		t.Fatalf("expected tasks %d, got %d", payload.Tasks, data.Tasks)
	}
}
