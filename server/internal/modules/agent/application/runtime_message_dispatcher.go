package application

import (
	"context"
	"errors"
	"fmt"
)

type runtimeMessageHandler func(service *AgentRuntimeService, ctx context.Context, agentID int, message RuntimeMessageInput) error

func newRuntimeMessageDispatcher() map[string]runtimeMessageHandler {
	builder := newRuntimeMessageDispatcherBuilder()
	registerHeartbeatMessageHandler(builder)
	return builder.handlers
}

type runtimeMessageDispatcherBuilder struct {
	handlers map[string]runtimeMessageHandler
}

func newRuntimeMessageDispatcherBuilder() *runtimeMessageDispatcherBuilder {
	return &runtimeMessageDispatcherBuilder{
		handlers: map[string]runtimeMessageHandler{},
	}
}

func (builder *runtimeMessageDispatcherBuilder) mustRegister(messageType string, handler runtimeMessageHandler) {
	// Handler registration happens at service bootstrap.
	// Panic here to fail fast on wiring mistakes (duplicate or invalid handlers).
	if err := builder.register(messageType, handler); err != nil {
		panic(err)
	}
}

func (builder *runtimeMessageDispatcherBuilder) register(messageType string, handler runtimeMessageHandler) error {
	if builder == nil {
		return errors.New("runtime message dispatcher builder is nil")
	}
	if messageType == "" {
		return errors.New("message type is required")
	}
	if handler == nil {
		return fmt.Errorf("handler is required: %s", messageType)
	}
	if _, exists := builder.handlers[messageType]; exists {
		return fmt.Errorf("duplicate runtime message handler: %s", messageType)
	}

	builder.handlers[messageType] = handler
	return nil
}
