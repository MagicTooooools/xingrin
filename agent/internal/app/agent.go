package app

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/yyhuni/orbit/agent/internal/config"
	"github.com/yyhuni/orbit/agent/internal/docker"
	"github.com/yyhuni/orbit/agent/internal/domain"
	"github.com/yyhuni/orbit/agent/internal/health"
	"github.com/yyhuni/orbit/agent/internal/metrics"
	"github.com/yyhuni/orbit/agent/internal/protocol"
	"github.com/yyhuni/orbit/agent/internal/task"
	"github.com/yyhuni/orbit/agent/internal/update"
	agentws "github.com/yyhuni/orbit/agent/internal/websocket"
)

func Run(ctx context.Context, cfg config.Config, wsURL string) error {
	configUpdater := config.NewUpdater(cfg)

	version := cfg.AgentVersion
	hostname := os.Getenv("ORBIT_HOSTNAME")
	if hostname == "" {
		var err error
		hostname, err = os.Hostname()
		if err != nil || hostname == "" {
			hostname = "unknown"
		}
	}

	client := agentws.NewClient(wsURL, cfg.APIKey)
	collector := metrics.NewCollector()
	healthManager := health.NewManager()
	taskCounter := &task.Counter{}
	heartbeat := agentws.NewHeartbeatSender(client, collector, healthManager, version, hostname, taskCounter.Count)

	taskClient := task.NewClient(cfg.ServerURL, cfg.APIKey)
	puller := task.NewPuller(taskClient, collector, taskCounter, cfg.MaxTasks, cfg.CPUThreshold, cfg.MemThreshold, cfg.DiskThreshold)

	taskQueue := make(chan *domain.Task, cfg.MaxTasks)
	puller.SetOnTask(func(t *domain.Task) {
		taskQueue <- t
	})

	dockerClient, err := docker.NewClient()
	if err != nil {
		log.Printf("docker client unavailable: %v", err)
	}

	workerToken := os.Getenv("WORKER_TOKEN")
	if workerToken == "" {
		return errors.New("WORKER_TOKEN environment variable is required")
	}

	executor := task.NewExecutor(dockerClient, taskClient, taskCounter, cfg.ServerURL, workerToken, version)
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := executor.Shutdown(shutdownCtx); err != nil && !errors.Is(err, context.DeadlineExceeded) {
			log.Printf("executor shutdown error: %v", err)
		}
	}()

	updater := update.NewUpdater(dockerClient, healthManager, puller, executor, configUpdater, cfg.APIKey, workerToken)

	handler := agentws.NewHandler()
	handler.OnTaskAvailable(puller.NotifyTaskAvailable)
	handler.OnTaskCancel(executor.CancelTask)
	handler.OnConfigUpdate(func(payload protocol.ConfigUpdatePayload) {
		cfgUpdate := config.Update{
			MaxTasks:      payload.MaxTasks,
			CPUThreshold:  payload.CPUThreshold,
			MemThreshold:  payload.MemThreshold,
			DiskThreshold: payload.DiskThreshold,
		}
		configUpdater.Apply(cfgUpdate)
		puller.UpdateConfig(cfgUpdate.MaxTasks, cfgUpdate.CPUThreshold, cfgUpdate.MemThreshold, cfgUpdate.DiskThreshold)
	})
	handler.OnUpdateRequired(updater.HandleUpdateRequired)
	client.SetOnMessage(handler.Handle)

	go heartbeat.Start(ctx)
	go func() {
		_ = puller.Run(ctx)
	}()
	go executor.Start(ctx, taskQueue)

	if err := client.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}
