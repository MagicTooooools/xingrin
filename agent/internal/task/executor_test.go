package task

import (
	"context"
	"testing"
	"time"

	"github.com/yyhuni/orbit/agent/internal/domain"
)

type fakeReporter struct {
	status string
	msg    string
}

func (f *fakeReporter) UpdateStatus(ctx context.Context, taskID int, status, errorMessage string) error {
	f.status = status
	f.msg = errorMessage
	return nil
}

func TestExecutorMissingWorkerToken(t *testing.T) {
	reporter := &fakeReporter{}
	exec := &Executor{
		client:      reporter,
		serverURL:   "https://server",
		workerToken: "",
	}

	exec.execute(context.Background(), &domain.Task{ID: 1})
	if reporter.status != "failed" {
		t.Fatalf("expected failed status, got %s", reporter.status)
	}
	if reporter.msg == "" {
		t.Fatalf("expected error message")
	}
}

func TestExecutorDockerUnavailable(t *testing.T) {
	reporter := &fakeReporter{}
	exec := &Executor{
		client:      reporter,
		serverURL:   "https://server",
		workerToken: "token",
	}

	exec.execute(context.Background(), &domain.Task{ID: 2})
	if reporter.status != "failed" {
		t.Fatalf("expected failed status, got %s", reporter.status)
	}
	if reporter.msg == "" {
		t.Fatalf("expected error message")
	}
}

func TestExecutorCancelAll(t *testing.T) {
	exec := &Executor{
		running: map[int]context.CancelFunc{},
	}
	calls := 0
	exec.running[1] = func() { calls++ }
	exec.running[2] = func() { calls++ }

	exec.CancelAll()
	if calls != 2 {
		t.Fatalf("expected cancel calls, got %d", calls)
	}
}

func TestExecutorShutdownWaits(t *testing.T) {
	exec := &Executor{
		running: map[int]context.CancelFunc{},
	}
	calls := 0
	exec.running[1] = func() { calls++ }

	exec.wg.Add(1)
	go func() {
		time.Sleep(10 * time.Millisecond)
		exec.wg.Done()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := exec.Shutdown(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected cancel call")
	}
}

func TestExecutorShutdownTimeout(t *testing.T) {
	exec := &Executor{
		running: map[int]context.CancelFunc{},
	}
	exec.wg.Add(1)
	defer exec.wg.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	if err := exec.Shutdown(ctx); err == nil {
		t.Fatalf("expected timeout error")
	}
}
