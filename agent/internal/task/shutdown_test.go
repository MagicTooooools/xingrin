package task

import (
	"context"
	"testing"
	"time"
)

func TestExecutorGracefulShutdownTimeout(t *testing.T) {
	exec := &Executor{
		running: map[int]context.CancelFunc{},
	}
	exec.wg.Add(1)
	defer exec.wg.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	if err := exec.GracefulShutdown(ctx); err == nil {
		t.Fatalf("expected timeout error")
	}
}
