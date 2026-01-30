package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/yyhuni/lunafox/agent/internal/app"
	"github.com/yyhuni/lunafox/agent/internal/config"
)

func main() {
	cfg, err := config.Load(os.Args[1:])
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	wsURL, err := config.BuildWebSocketURL(cfg.ServerURL)
	if err != nil {
		log.Fatalf("invalid server URL: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := app.Run(ctx, *cfg, wsURL); err != nil {
		log.Fatalf("agent stopped: %v", err)
	}
}
