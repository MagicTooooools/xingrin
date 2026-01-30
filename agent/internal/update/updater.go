package update

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/strslice"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/yyhuni/lunafox/agent/internal/config"
	"github.com/yyhuni/lunafox/agent/internal/domain"
)

// Updater handles agent self-update.
type Updater struct {
	docker     dockerClient
	health     healthSetter
	puller     pullerController
	executor   executorController
	cfg        configSnapshot
	apiKey     string
	token      string
	mu         sync.Mutex
	updating   bool
	randSrc    *rand.Rand
	backoff    time.Duration
	maxBackoff time.Duration
}

type dockerClient interface {
	ImagePull(ctx context.Context, imageRef string) (io.ReadCloser, error)
	ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, platform *ocispec.Platform, name string) (container.CreateResponse, error)
	ContainerStart(ctx context.Context, containerID string, opts container.StartOptions) error
}

type healthSetter interface {
	Set(state, reason, message string)
}

type pullerController interface {
	Pause()
}

type executorController interface {
	Shutdown(ctx context.Context) error
}

type configSnapshot interface {
	Snapshot() config.Config
}

// NewUpdater creates a new updater.
func NewUpdater(dockerClient dockerClient, healthManager healthSetter, puller pullerController, executor executorController, cfg configSnapshot, apiKey, token string) *Updater {
	return &Updater{
		docker:     dockerClient,
		health:     healthManager,
		puller:     puller,
		executor:   executor,
		cfg:        cfg,
		apiKey:     apiKey,
		token:      token,
		randSrc:    rand.New(rand.NewSource(time.Now().UnixNano())),
		backoff:    30 * time.Second,
		maxBackoff: 10 * time.Minute,
	}
}

// HandleUpdateRequired triggers the update flow.
func (u *Updater) HandleUpdateRequired(payload domain.UpdateRequiredPayload) {
	u.mu.Lock()
	if u.updating {
		u.mu.Unlock()
		return
	}
	u.updating = true
	u.mu.Unlock()

	go u.run(payload)
}

func (u *Updater) run(payload domain.UpdateRequiredPayload) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("agent update panic: %v", r)
			u.health.Set("paused", "update_panic", fmt.Sprintf("%v", r))
		}
		u.mu.Lock()
		u.updating = false
		u.mu.Unlock()
	}()
	u.puller.Pause()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	_ = u.executor.Shutdown(ctx)
	cancel()

	for {
		if err := u.updateOnce(payload); err == nil {
			u.health.Set("ok", "", "")
			os.Exit(0)
		} else {
			u.health.Set("paused", "update_failed", err.Error())
		}

		delay := withJitter(u.backoff, u.randSrc)
		if u.backoff < u.maxBackoff {
			u.backoff *= 2
			if u.backoff > u.maxBackoff {
				u.backoff = u.maxBackoff
			}
		}
		time.Sleep(delay)
	}
}

func (u *Updater) updateOnce(payload domain.UpdateRequiredPayload) error {
	if u.docker == nil {
		return fmt.Errorf("docker client unavailable")
	}
	image := strings.TrimSpace(payload.Image)
	version := strings.TrimSpace(payload.Version)
	if image == "" || version == "" {
		return fmt.Errorf("invalid update payload")
	}

	// Strict validation: reject invalid data from server
	if err := validateImageName(image); err != nil {
		log.Printf("invalid image name from server: %s, error: %v", image, err)
		return fmt.Errorf("invalid image name from server: %w", err)
	}
	if err := validateVersion(version); err != nil {
		log.Printf("invalid version from server: %s, error: %v", version, err)
		return fmt.Errorf("invalid version from server: %w", err)
	}

	fullImage := fmt.Sprintf("%s:%s", image, version)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	reader, err := u.docker.ImagePull(ctx, fullImage)
	if err != nil {
		return err
	}
	_, _ = io.Copy(io.Discard, reader)
	_ = reader.Close()

	if err := u.startNewContainer(ctx, image, version); err != nil {
		return err
	}

	return nil
}

func (u *Updater) startNewContainer(ctx context.Context, image, version string) error {
	env := []string{
		fmt.Sprintf("SERVER_URL=%s", u.cfg.Snapshot().ServerURL),
		fmt.Sprintf("API_KEY=%s", u.apiKey),
		fmt.Sprintf("MAX_TASKS=%d", u.cfg.Snapshot().MaxTasks),
		fmt.Sprintf("CPU_THRESHOLD=%d", u.cfg.Snapshot().CPUThreshold),
		fmt.Sprintf("MEM_THRESHOLD=%d", u.cfg.Snapshot().MemThreshold),
		fmt.Sprintf("DISK_THRESHOLD=%d", u.cfg.Snapshot().DiskThreshold),
		fmt.Sprintf("AGENT_VERSION=%s", version),
	}
	if u.token != "" {
		env = append(env, fmt.Sprintf("WORKER_TOKEN=%s", u.token))
	}

	cfg := &container.Config{
		Image: fmt.Sprintf("%s:%s", image, version),
		Env:   env,
		Cmd:   strslice.StrSlice{},
	}

	hostConfig := &container.HostConfig{
		Binds: []string{
			"/var/run/docker.sock:/var/run/docker.sock",
			"/opt/lunafox:/opt/lunafox",
		},
		RestartPolicy: container.RestartPolicy{Name: "unless-stopped"},
		OomScoreAdj:   -500,
	}

	// Version is already validated, just normalize to lowercase for container name
	name := fmt.Sprintf("lunafox-agent-%s", strings.ToLower(version))
	resp, err := u.docker.ContainerCreate(ctx, cfg, hostConfig, &network.NetworkingConfig{}, nil, name)
	if err != nil {
		return err
	}

	if err := u.docker.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return err
	}

	log.Printf("agent update started new container: %s", resp.ID)
	return nil
}

func withJitter(delay time.Duration, src *rand.Rand) time.Duration {
	if delay <= 0 || src == nil {
		return delay
	}
	jitter := src.Float64() * 0.2
	return delay + time.Duration(float64(delay)*jitter)
}

// validateImageName validates that the image name contains only safe characters.
// Returns error if validation fails.
func validateImageName(image string) error {
	if len(image) == 0 {
		return fmt.Errorf("image name cannot be empty")
	}
	if len(image) > 255 {
		return fmt.Errorf("image name too long: %d characters", len(image))
	}

	// Allow: alphanumeric, dots, hyphens, underscores, slashes (for registry paths)
	for i, r := range image {
		if !((r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') ||
			r == '.' || r == '-' || r == '_' || r == '/') {
			return fmt.Errorf("invalid character at position %d: %c", i, r)
		}
	}

	// Must not start or end with special characters
	first := rune(image[0])
	last := rune(image[len(image)-1])
	if first == '.' || first == '-' || first == '/' {
		return fmt.Errorf("image name cannot start with special character: %c", first)
	}
	if last == '.' || last == '-' || last == '/' {
		return fmt.Errorf("image name cannot end with special character: %c", last)
	}

	return nil
}

// validateVersion validates that the version string contains only safe characters.
// Returns error if validation fails.
func validateVersion(version string) error {
	if len(version) == 0 {
		return fmt.Errorf("version cannot be empty")
	}
	if len(version) > 128 {
		return fmt.Errorf("version too long: %d characters", len(version))
	}

	// Allow: alphanumeric, dots, hyphens, underscores
	for i, r := range version {
		if !((r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') ||
			r == '.' || r == '-' || r == '_') {
			return fmt.Errorf("invalid character at position %d: %c", i, r)
		}
	}

	// Must not start or end with special characters
	first := rune(version[0])
	last := rune(version[len(version)-1])
	if first == '.' || first == '-' || first == '_' {
		return fmt.Errorf("version cannot start with special character: %c", first)
	}
	if last == '.' || last == '-' || last == '_' {
		return fmt.Errorf("version cannot end with special character: %c", last)
	}

	return nil
}
