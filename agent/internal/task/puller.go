package task

import (
	"context"
	"errors"
	"math"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/yyhuni/lunafox/agent/internal/domain"
)

// Puller coordinates task pulling with load gating and backoff.
type Puller struct {
	client        TaskPuller
	collector     MetricsSampler
	counter       *Counter
	maxTasks      int
	cpuThreshold  int
	memThreshold  int
	diskThreshold int

	onTask       func(*domain.Task)
	notifyCh     chan struct{}
	emptyBackoff []time.Duration
	emptyIdx     int
	errorBackoff time.Duration
	errorMax     time.Duration
	randSrc      *rand.Rand
	mu           sync.RWMutex
	paused       atomic.Bool
}

type MetricsSampler interface {
	Sample() (float64, float64, float64)
}

type TaskPuller interface {
	PullTask(ctx context.Context) (*domain.Task, error)
}

// NewPuller creates a new Puller.
func NewPuller(client TaskPuller, collector MetricsSampler, counter *Counter, maxTasks, cpuThreshold, memThreshold, diskThreshold int) *Puller {
	return &Puller{
		client:        client,
		collector:     collector,
		counter:       counter,
		maxTasks:      maxTasks,
		cpuThreshold:  cpuThreshold,
		memThreshold:  memThreshold,
		diskThreshold: diskThreshold,
		notifyCh:      make(chan struct{}, 1),
		emptyBackoff:  []time.Duration{5 * time.Second, 10 * time.Second, 30 * time.Second, 60 * time.Second},
		errorBackoff:  1 * time.Second,
		errorMax:      60 * time.Second,
		randSrc:       rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// SetOnTask registers a callback invoked when a task is assigned.
func (p *Puller) SetOnTask(fn func(*domain.Task)) {
	p.onTask = fn
}

// NotifyTaskAvailable triggers an immediate pull attempt.
func (p *Puller) NotifyTaskAvailable() {
	select {
	case p.notifyCh <- struct{}{}:
	default:
	}
}

// Run starts the pull loop.
func (p *Puller) Run(ctx context.Context) error {
	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		if p.paused.Load() {
			if !p.waitUntilCanceled(ctx) {
				return ctx.Err()
			}
			continue
		}

		loadInterval := p.loadInterval()
		if !p.canPull() {
			if !p.wait(ctx, loadInterval) {
				return ctx.Err()
			}
			continue
		}

		task, err := p.client.PullTask(ctx)
		if err != nil {
			delay := p.nextErrorBackoff()
			if !p.wait(ctx, delay) {
				return ctx.Err()
			}
			continue
		}

		p.resetErrorBackoff()
		if task == nil {
			delay := p.nextEmptyDelay(loadInterval)
			if !p.waitOrNotify(ctx, delay) {
				return ctx.Err()
			}
			continue
		}

		p.resetEmptyBackoff()
		if p.onTask != nil {
			p.onTask(task)
		}
	}
}

func (p *Puller) canPull() bool {
	maxTasks, cpuThreshold, memThreshold, diskThreshold := p.currentConfig()
	if p.counter != nil && p.counter.Count() >= maxTasks {
		return false
	}
	cpu, mem, disk := p.collector.Sample()
	return cpu < float64(cpuThreshold) &&
		mem < float64(memThreshold) &&
		disk < float64(diskThreshold)
}

func (p *Puller) loadInterval() time.Duration {
	cpu, mem, disk := p.collector.Sample()
	load := math.Max(cpu, math.Max(mem, disk))
	switch {
	case load < 50:
		return 1 * time.Second
	case load < 80:
		return 3 * time.Second
	default:
		return 10 * time.Second
	}
}

func (p *Puller) nextEmptyDelay(loadInterval time.Duration) time.Duration {
	var empty time.Duration
	if p.emptyIdx < len(p.emptyBackoff) {
		empty = p.emptyBackoff[p.emptyIdx]
		p.emptyIdx++
	} else {
		empty = p.emptyBackoff[len(p.emptyBackoff)-1]
	}
	if empty < loadInterval {
		return loadInterval
	}
	return empty
}

func (p *Puller) resetEmptyBackoff() {
	p.emptyIdx = 0
}

func (p *Puller) nextErrorBackoff() time.Duration {
	delay := p.errorBackoff
	next := delay * 2
	if next > p.errorMax {
		next = p.errorMax
	}
	p.errorBackoff = next
	return withJitter(delay, p.randSrc)
}

func (p *Puller) resetErrorBackoff() {
	p.errorBackoff = 1 * time.Second
}

func (p *Puller) wait(ctx context.Context, delay time.Duration) bool {
	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}

func (p *Puller) waitOrNotify(ctx context.Context, delay time.Duration) bool {
	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return false
	case <-p.notifyCh:
		return true
	case <-timer.C:
		return true
	}
}

func withJitter(delay time.Duration, src *rand.Rand) time.Duration {
	if delay <= 0 || src == nil {
		return delay
	}
	jitter := src.Float64() * 0.2
	return delay + time.Duration(float64(delay)*jitter)
}

func (p *Puller) EnsureTaskHandler() error {
	if p.onTask == nil {
		return errors.New("task handler is required")
	}
	return nil
}

// Pause stops pulling. Once paused, only context cancellation exits the loop.
func (p *Puller) Pause() {
	p.paused.Store(true)
}

// UpdateConfig updates puller thresholds and max tasks.
func (p *Puller) UpdateConfig(maxTasks, cpuThreshold, memThreshold, diskThreshold *int) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if maxTasks != nil && *maxTasks > 0 {
		p.maxTasks = *maxTasks
	}
	if cpuThreshold != nil && *cpuThreshold > 0 {
		p.cpuThreshold = *cpuThreshold
	}
	if memThreshold != nil && *memThreshold > 0 {
		p.memThreshold = *memThreshold
	}
	if diskThreshold != nil && *diskThreshold > 0 {
		p.diskThreshold = *diskThreshold
	}
}

func (p *Puller) currentConfig() (int, int, int, int) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.maxTasks, p.cpuThreshold, p.memThreshold, p.diskThreshold
}

func (p *Puller) waitUntilCanceled(ctx context.Context) bool {
	<-ctx.Done()
	return false
}
