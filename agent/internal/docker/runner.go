package docker

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/strslice"
	"github.com/yyhuni/lunafox/agent/internal/domain"
)

const workerImagePrefix = "yyhuni/lunafox-worker:"

// StartWorker starts a worker container for a task and returns the container ID.
func (c *Client) StartWorker(ctx context.Context, t *domain.Task, serverURL, serverToken, agentVersion string) (string, error) {
	if t == nil {
		return "", fmt.Errorf("task is nil")
	}
	if err := os.MkdirAll(t.WorkspaceDir, 0755); err != nil {
		return "", fmt.Errorf("prepare workspace: %w", err)
	}

	image, err := resolveWorkerImage(agentVersion)
	if err != nil {
		return "", err
	}
	env := buildWorkerEnv(t, serverURL, serverToken)

	config := &container.Config{
		Image: image,
		Env:   env,
		Cmd:   strslice.StrSlice{},
	}

	hostConfig := &container.HostConfig{
		Binds:       []string{"/opt/lunafox:/opt/lunafox"},
		AutoRemove:  false,
		OomScoreAdj: 500,
	}

	resp, err := c.cli.ContainerCreate(ctx, config, hostConfig, &network.NetworkingConfig{}, nil, "")
	if err != nil {
		return "", err
	}

	if err := c.cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", err
	}

	return resp.ID, nil
}

func resolveWorkerImage(version string) (string, error) {
	version = strings.TrimSpace(version)
	if version == "" {
		return "", fmt.Errorf("worker version is required")
	}
	return workerImagePrefix + version, nil
}

func buildWorkerEnv(t *domain.Task, serverURL, serverToken string) []string {
	return []string{
		fmt.Sprintf("SERVER_URL=%s", serverURL),
		fmt.Sprintf("SERVER_TOKEN=%s", serverToken),
		fmt.Sprintf("SCAN_ID=%d", t.ScanID),
		fmt.Sprintf("TARGET_ID=%d", t.TargetID),
		fmt.Sprintf("TARGET_NAME=%s", t.TargetName),
		fmt.Sprintf("TARGET_TYPE=%s", t.TargetType),
		fmt.Sprintf("WORKFLOW_NAME=%s", t.WorkflowName),
		fmt.Sprintf("WORKSPACE_DIR=%s", t.WorkspaceDir),
		fmt.Sprintf("CONFIG=%s", t.Config),
	}
}
