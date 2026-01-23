package websocket

import "testing"

func TestHandlersTaskAvailable(t *testing.T) {
	h := NewHandler()
	called := 0
	h.OnTaskAvailable(func() { called++ })

	h.Handle([]byte(`{"type":"task_available","payload":{},"timestamp":"2026-01-01T00:00:00Z"}`))
	if called != 1 {
		t.Fatalf("expected callback to be called")
	}
}

func TestHandlersTaskCancel(t *testing.T) {
	h := NewHandler()
	var got int
	h.OnTaskCancel(func(id int) { got = id })

	h.Handle([]byte(`{"type":"task_cancel","payload":{"taskId":123},"timestamp":"2026-01-01T00:00:00Z"}`))
	if got != 123 {
		t.Fatalf("expected taskId 123")
	}
}

func TestHandlersConfigUpdate(t *testing.T) {
	h := NewHandler()
	var maxTasks int
	h.OnConfigUpdate(func(payload ConfigUpdatePayload) {
		if payload.MaxTasks != nil {
			maxTasks = *payload.MaxTasks
		}
	})

	h.Handle([]byte(`{"type":"config_update","payload":{"maxTasks":8},"timestamp":"2026-01-01T00:00:00Z"}`))
	if maxTasks != 8 {
		t.Fatalf("expected maxTasks 8")
	}
}

func TestHandlersUpdateRequired(t *testing.T) {
	h := NewHandler()
	var version string
	h.OnUpdateRequired(func(payload UpdateRequiredPayload) { version = payload.Version })

	h.Handle([]byte(`{"type":"update_required","payload":{"version":"v1.0.1","image":"yyhuni/orbit-agent"},"timestamp":"2026-01-01T00:00:00Z"}`))
	if version != "v1.0.1" {
		t.Fatalf("expected version")
	}
}
