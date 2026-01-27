package task

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/yyhuni/orbit/agent/internal/docker"
	"github.com/yyhuni/orbit/agent/internal/domain"
)

const defaultMaxRuntime = 7 * 24 * time.Hour

// Executor runs tasks inside worker containers.
type Executor struct {
	docker       DockerRunner
	client       statusReporter
	counter      *Counter
	serverURL    string
	workerToken  string
	agentVersion string
	maxRuntime   time.Duration

	mu        sync.Mutex
	running   map[int]context.CancelFunc
	cancelMu  sync.Mutex
	cancelled map[int]struct{}
	wg        sync.WaitGroup

	stopping atomic.Bool
}

type statusReporter interface {
	UpdateStatus(ctx context.Context, taskID int, status, errorMessage string) error
}

type DockerRunner interface {
	StartWorker(ctx context.Context, t *domain.Task, serverURL, serverToken, agentVersion string) (string, error)
	Wait(ctx context.Context, containerID string) (int64, error)
	Stop(ctx context.Context, containerID string) error
	Remove(ctx context.Context, containerID string) error
	TailLogs(ctx context.Context, containerID string, lines int) (string, error)
}

// NewExecutor creates an Executor.
func NewExecutor(dockerClient DockerRunner, taskClient statusReporter, counter *Counter, serverURL, workerToken, agentVersion string) *Executor {
	return &Executor{
		docker:       dockerClient,
		client:       taskClient,
		counter:      counter,
		serverURL:    serverURL,
		workerToken:  workerToken,
		agentVersion: agentVersion,
		maxRuntime:   defaultMaxRuntime,
		running:      map[int]context.CancelFunc{},
		cancelled:    map[int]struct{}{},
	}
}

// Start processes tasks from the queue.
func (e *Executor) Start(ctx context.Context, tasks <-chan *domain.Task) {
	for {
		select {
		case <-ctx.Done():
			return
		case t, ok := <-tasks:
			if !ok {
				return
			}
			if t == nil {
				continue
			}
			if e.stopping.Load() {
				// During shutdown/update: drain the queue but don't start new work.
				continue
			}
			if e.isCancelled(t.ID) {
				e.reportStatus(ctx, t.ID, "cancelled", "")
				e.clearCancelled(t.ID)
				continue
			}
			go e.execute(ctx, t)
		}
	}
}

// CancelTask requests cancellation of a running task.
func (e *Executor) CancelTask(taskID int) {
	e.mu.Lock()
	cancel := e.running[taskID]
	e.mu.Unlock()
	if cancel != nil {
		cancel()
	}
}

// MarkCancelled records a task as cancelled to prevent execution.
func (e *Executor) MarkCancelled(taskID int) {
	e.cancelMu.Lock()
	e.cancelled[taskID] = struct{}{}
	e.cancelMu.Unlock()
}

func (e *Executor) reportStatus(ctx context.Context, taskID int, status, errorMessage string) {
	if e.client == nil {
		return
	}
	statusCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 30*time.Second)
	defer cancel()
	_ = e.client.UpdateStatus(statusCtx, taskID, status, errorMessage)
}

func (e *Executor) execute(ctx context.Context, t *domain.Task) {
	e.wg.Add(1)
	defer e.wg.Done()
	defer e.clearCancelled(t.ID)

	if e.counter != nil {
		e.counter.Inc()
		defer e.counter.Dec()
	}

	if e.workerToken == "" {
		e.reportStatus(ctx, t.ID, "failed", "missing worker token")
		return
	}
	if e.docker == nil {
		e.reportStatus(ctx, t.ID, "failed", "docker client unavailable")
		return
	}

	runCtx, cancel := context.WithTimeout(ctx, e.maxRuntime)
	defer cancel()

	containerID, err := e.docker.StartWorker(runCtx, t, e.serverURL, e.workerToken, e.agentVersion)
	if err != nil {
		message := docker.TruncateErrorMessage(err.Error())
		e.reportStatus(ctx, t.ID, "failed", message)
		return
	}
	defer func() {
		_ = e.docker.Remove(context.Background(), containerID)
	}()

	e.trackCancel(t.ID, cancel)
	defer e.clearCancel(t.ID)

	exitCode, waitErr := e.docker.Wait(runCtx, containerID)
	if waitErr != nil {
		if errors.Is(waitErr, context.DeadlineExceeded) || errors.Is(runCtx.Err(), context.DeadlineExceeded) {
			e.handleTimeout(ctx, t, containerID)
			return
		}
		if errors.Is(waitErr, context.Canceled) || errors.Is(runCtx.Err(), context.Canceled) {
			e.handleCancel(ctx, t, containerID)
			return
		}
		message := docker.TruncateErrorMessage(waitErr.Error())
		e.reportStatus(ctx, t.ID, "failed", message)
		return
	}

	if runCtx.Err() != nil {
		if errors.Is(runCtx.Err(), context.DeadlineExceeded) {
			e.handleTimeout(ctx, t, containerID)
			return
		}
		if errors.Is(runCtx.Err(), context.Canceled) {
			e.handleCancel(ctx, t, containerID)
			return
		}
	}

	if exitCode == 0 {
		e.reportStatus(ctx, t.ID, "completed", "")
		return
	}

	logs, _ := e.docker.TailLogs(context.Background(), containerID, 100)
	message := logs
	if message == "" {
		message = fmt.Sprintf("container exited with code %d", exitCode)
	}
	message = docker.TruncateErrorMessage(message)
	e.reportStatus(ctx, t.ID, "failed", message)
}

func (e *Executor) handleCancel(ctx context.Context, t *domain.Task, containerID string) {
	_ = e.docker.Stop(context.Background(), containerID)
	e.reportStatus(ctx, t.ID, "cancelled", "")
}

func (e *Executor) handleTimeout(ctx context.Context, t *domain.Task, containerID string) {
	_ = e.docker.Stop(context.Background(), containerID)
	message := docker.TruncateErrorMessage("task timed out")
	e.reportStatus(ctx, t.ID, "failed", message)
}

func (e *Executor) trackCancel(taskID int, cancel context.CancelFunc) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.running[taskID] = cancel
}

func (e *Executor) clearCancel(taskID int) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.running, taskID)
}

func (e *Executor) isCancelled(taskID int) bool {
	e.cancelMu.Lock()
	defer e.cancelMu.Unlock()
	_, ok := e.cancelled[taskID]
	return ok
}

func (e *Executor) clearCancelled(taskID int) {
	e.cancelMu.Lock()
	delete(e.cancelled, taskID)
	e.cancelMu.Unlock()
}

// CancelAll requests cancellation for all running tasks.
func (e *Executor) CancelAll() {
	e.mu.Lock()
	cancels := make([]context.CancelFunc, 0, len(e.running))
	for _, cancel := range e.running {
		cancels = append(cancels, cancel)
	}
	e.mu.Unlock()

	for _, cancel := range cancels {
		cancel()
	}
}

// Shutdown cancels running tasks and waits for completion.
func (e *Executor) Shutdown(ctx context.Context) error {
	e.stopping.Store(true)
	e.CancelAll()

	done := make(chan struct{})
	go func() {
		e.wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return nil
	}
}
