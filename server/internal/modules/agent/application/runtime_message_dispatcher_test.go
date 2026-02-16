package application

import (
	"context"
	"testing"
)

func TestRuntimeMessageDispatcherBuilderRegisterDuplicate(t *testing.T) {
	builder := newRuntimeMessageDispatcherBuilder()
	handler := func(_ *AgentRuntimeService, _ context.Context, _ int, _ RuntimeMessageInput) error { return nil }

	if err := builder.register("heartbeat", handler); err != nil {
		t.Fatalf("expected first register success, got %v", err)
	}
	if err := builder.register("heartbeat", handler); err == nil {
		t.Fatalf("expected duplicate register error")
	}
}

func TestRuntimeMessageDispatcherBuilderRegisterValidation(t *testing.T) {
	builder := newRuntimeMessageDispatcherBuilder()
	handler := func(_ *AgentRuntimeService, _ context.Context, _ int, _ RuntimeMessageInput) error { return nil }

	if err := builder.register("", handler); err == nil {
		t.Fatalf("expected empty message type error")
	}
	if err := builder.register("heartbeat", nil); err == nil {
		t.Fatalf("expected nil handler error")
	}
}

func TestNewRuntimeMessageDispatcherIncludesHeartbeat(t *testing.T) {
	handlers := newRuntimeMessageDispatcher()
	if handlers[RuntimeMessageTypeHeartbeat] == nil {
		t.Fatalf("expected heartbeat handler to be registered")
	}
}
