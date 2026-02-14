package task

import (
	"context"
	"testing"
	"time"

	"github.com/yyhuni/lunafox/agent/internal/domain"
)

type fakeReporter struct {
	status string
	msg    string
}

type fakeDockerRunner struct {
	startWorkerFn func(ctx context.Context, t *domain.Task, serverURL, serverToken, agentVersion string) (string, error)
	waitFn        func(ctx context.Context, containerID string) (int64, error)
	stopFn        func(ctx context.Context, containerID string) error
	removeFn      func(ctx context.Context, containerID string) error
	tailLogsFn    func(ctx context.Context, containerID string, lines int) (string, error)
}

func (fake *fakeDockerRunner) StartWorker(ctx context.Context, t *domain.Task, serverURL, serverToken, agentVersion string) (string, error) {
	if fake.startWorkerFn == nil {
		return "container-1", nil
	}
	return fake.startWorkerFn(ctx, t, serverURL, serverToken, agentVersion)
}

func (fake *fakeDockerRunner) Wait(ctx context.Context, containerID string) (int64, error) {
	if fake.waitFn == nil {
		return 0, nil
	}
	return fake.waitFn(ctx, containerID)
}

func (fake *fakeDockerRunner) Stop(ctx context.Context, containerID string) error {
	if fake.stopFn == nil {
		return nil
	}
	return fake.stopFn(ctx, containerID)
}

func (fake *fakeDockerRunner) Remove(ctx context.Context, containerID string) error {
	if fake.removeFn == nil {
		return nil
	}
	return fake.removeFn(ctx, containerID)
}

func (fake *fakeDockerRunner) TailLogs(ctx context.Context, containerID string, lines int) (string, error) {
	if fake.tailLogsFn == nil {
		return "", nil
	}
	return fake.tailLogsFn(ctx, containerID, lines)
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

func TestExecutorFailurePathUsesTimeoutContexts(t *testing.T) {
	reporter := &fakeReporter{}
	tailLogsHasDeadline := false
	removeHasDeadline := false

	fakeDocker := &fakeDockerRunner{
		startWorkerFn: func(ctx context.Context, t *domain.Task, serverURL, serverToken, agentVersion string) (string, error) {
			return "container-1", nil
		},
		waitFn: func(ctx context.Context, containerID string) (int64, error) {
			return 1, nil
		},
		tailLogsFn: func(ctx context.Context, containerID string, lines int) (string, error) {
			_, tailLogsHasDeadline = ctx.Deadline()
			return "tool failed", nil
		},
		removeFn: func(ctx context.Context, containerID string) error {
			_, removeHasDeadline = ctx.Deadline()
			return nil
		},
	}

	exec := NewExecutor(fakeDocker, reporter, nil, "https://server", "token", "v1")
	exec.execute(context.Background(), &domain.Task{ID: 10, ScanID: 20})

	if reporter.status != "failed" {
		t.Fatalf("expected failed status, got %s", reporter.status)
	}
	if !tailLogsHasDeadline {
		t.Fatalf("expected tail logs context to have deadline")
	}
	if !removeHasDeadline {
		t.Fatalf("expected remove context to have deadline")
	}
}

func TestExecutorHandleTimeoutUsesDeadlineOnStop(t *testing.T) {
	reporter := &fakeReporter{}
	stopHasDeadline := false

	fakeDocker := &fakeDockerRunner{
		stopFn: func(ctx context.Context, containerID string) error {
			_, stopHasDeadline = ctx.Deadline()
			return nil
		},
	}

	exec := NewExecutor(fakeDocker, reporter, nil, "https://server", "token", "v1")
	exec.handleTimeout(context.Background(), &domain.Task{ID: 1, ScanID: 2}, "container-1")

	if !stopHasDeadline {
		t.Fatalf("expected stop context to have deadline")
	}
	if reporter.status != "failed" {
		t.Fatalf("expected failed status, got %s", reporter.status)
	}
}
