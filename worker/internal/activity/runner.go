package activity

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/orbit/worker/internal/pkg"
	"go.uber.org/zap"
)

// ansiRegex matches ANSI escape sequences (colors, cursor movement, etc.)
var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

// controlCharReplacer removes control characters in a single pass
var controlCharReplacer = strings.NewReplacer(
	"\x00", "", // NUL
	"\r", "", // CR
	"\b", "", // Backspace
	"\f", "", // Form feed
	"\v", "", // Vertical tab
)

const (
	DefaultDirPerm  = 0755
	ExitCodeTimeout = -1
	ExitCodeError   = -2

	// Runner concurrency control (external command processes)
	// Enforced per worker container/workflow instance.
	EnvMaxCmdConcurrency     = "WORKER_MAX_CMD_CONCURRENCY"
	DefaultMaxCmdConcurrency = 2

	// Scanner buffer sizes
	ScannerInitBufSize = 64 * 1024   // 64KB initial buffer
	ScannerMaxBufSize  = 1024 * 1024 // 1MB max buffer for long lines
)

// Result represents the result of an activity execution
type Result struct {
	Name       string
	OutputFile string
	LogFile    string
	ExitCode   int
	Duration   time.Duration
	Error      error
}

// Command represents a command to execute
type Command struct {
	Name       string
	Command    string
	OutputFile string
	LogFile    string
	Timeout    time.Duration
}

// Runner executes activities (external tools)
type Runner struct {
	workDir string
	sem     chan struct{}
}

// NewRunner creates a new activity runner
func NewRunner(workDir string) *Runner {
	maxConc := getMaxCmdConcurrency()
	return &Runner{
		workDir: workDir,
		sem:     make(chan struct{}, maxConc),
	}
}

func getMaxCmdConcurrency() int {
	v := os.Getenv(EnvMaxCmdConcurrency)
	if v == "" {
		return DefaultMaxCmdConcurrency
	}

	n, err := strconv.Atoi(v)
	if err != nil || n < 1 {
		if pkg.Logger != nil {
			pkg.Logger.Warn("Invalid max command concurrency; using default",
				zap.String("env", EnvMaxCmdConcurrency),
				zap.String("value", v),
				zap.Int("default", DefaultMaxCmdConcurrency))
		}
		return DefaultMaxCmdConcurrency
	}

	return n
}

func (r *Runner) acquire(ctx context.Context) error {
	select {
	case r.sem <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (r *Runner) release() {
	<-r.sem
}

// killProcessGroup terminates the entire process group
// When using shell=true (sh -c), the actual tool runs as a child of the shell.
// If we only kill the shell process, the child becomes an orphan and keeps running.
// By killing the process group, we ensure all child processes are terminated.
func killProcessGroup(cmd *exec.Cmd) {
	if cmd == nil || cmd.Process == nil {
		return
	}

	pid := cmd.Process.Pid

	// Try to kill the process group first
	// The negative PID signals the entire process group
	if err := syscall.Kill(-pid, syscall.SIGKILL); err != nil {
		pkg.Logger.Debug("Failed to kill process group, trying single process",
			zap.Int("pid", pid),
			zap.Error(err))
		// Fallback: kill single process
		_ = cmd.Process.Kill()
	} else {
		pkg.Logger.Debug("Killed process group", zap.Int("pgid", pid))
	}
}

// Run executes a single activity with streaming output
func (r *Runner) Run(ctx context.Context, cmd Command) *Result {
	start := time.Now()
	result := &Result{
		Name:       cmd.Name,
		OutputFile: cmd.OutputFile,
		LogFile:    cmd.LogFile,
	}

	if ctx.Err() != nil {
		result.Error = fmt.Errorf("context cancelled before execution: %w", ctx.Err())
		result.ExitCode = ExitCodeError
		return result
	}
	if cmd.Timeout <= 0 {
		result.Error = fmt.Errorf("invalid timeout %v: must be > 0", cmd.Timeout)
		result.ExitCode = ExitCodeError
		return result
	}

	// Limit external command process concurrency per worker container.
	if err := r.acquire(ctx); err != nil {
		result.Error = fmt.Errorf("context cancelled while waiting for command slot: %w", err)
		result.ExitCode = ExitCodeError
		return result
	}
	defer r.release()

	execCtx, cancel := context.WithTimeout(ctx, cmd.Timeout)
	defer cancel()

	execCmd := exec.CommandContext(execCtx, "sh", "-c", cmd.Command)
	execCmd.Dir = r.workDir
	// Create new process group so we can kill all child processes
	execCmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	// Setup pipes for streaming
	stdout, err := execCmd.StdoutPipe()
	if err != nil {
		result.Error = fmt.Errorf("failed to create stdout pipe: %w", err)
		result.ExitCode = ExitCodeError
		return result
	}

	stderr, err := execCmd.StderrPipe()
	if err != nil {
		result.Error = fmt.Errorf("failed to create stderr pipe: %w", err)
		result.ExitCode = ExitCodeError
		return result
	}

	// Prepare log file
	logFile := r.prepareLogFile(cmd)
	if logFile != nil {
		defer func() { _ = logFile.Close() }()
		r.writeLogHeader(logFile, cmd)
	}

	// Start command
	if err := execCmd.Start(); err != nil {
		result.Error = fmt.Errorf("failed to start command: %w", err)
		result.ExitCode = ExitCodeError
		return result
	}

	// Ensure process cleanup on any exit path
	defer killProcessGroup(execCmd)

	// Stream output
	var wg sync.WaitGroup
	wg.Add(2)

	go r.streamOutput(&wg, stdout, logFile, cmd.Name, "stdout")
	go r.streamOutput(&wg, stderr, logFile, cmd.Name, "stderr")

	wg.Wait()

	// Wait for command to finish
	err = execCmd.Wait()
	result.Duration = time.Since(start)

	// Write duration to log
	if logFile != nil {
		r.writeLogFooter(logFile, result)
	}

	// Handle result
	if execCtx.Err() == context.DeadlineExceeded {
		result.Error = fmt.Errorf("activity execution timeout after %v", cmd.Timeout)
		result.ExitCode = ExitCodeTimeout
		pkg.Logger.Error("Activity timeout",
			zap.String("activity", cmd.Name),
			zap.Duration("timeout", cmd.Timeout))
		return result
	}

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		} else {
			result.ExitCode = ExitCodeError
		}
		result.Error = fmt.Errorf("activity execution failed: %w", err)
		pkg.Logger.Error("Activity failed",
			zap.String("activity", cmd.Name),
			zap.Int("exitCode", result.ExitCode),
			zap.Error(err))
		return result
	}

	result.ExitCode = 0
	pkg.Logger.Info("Activity completed",
		zap.String("activity", cmd.Name),
		zap.Duration("duration", result.Duration))

	return result
}

// RunParallel executes multiple activities in parallel
func (r *Runner) RunParallel(ctx context.Context, commands []Command) []*Result {
	if len(commands) == 0 {
		return nil
	}

	results := make([]*Result, len(commands))
	var wg sync.WaitGroup

	for i, cmd := range commands {
		if ctx.Err() != nil {
			results[i] = &Result{
				Name:     cmd.Name,
				ExitCode: ExitCodeError,
				Error:    fmt.Errorf("context cancelled: %w", ctx.Err()),
			}
			continue
		}

		wg.Add(1)
		go func(idx int, c Command) {
			defer wg.Done()
			results[idx] = r.Run(ctx, c)
		}(i, cmd)
	}

	wg.Wait()
	return results
}

// RunSequential executes multiple activities sequentially (one after another)
func (r *Runner) RunSequential(ctx context.Context, commands []Command) []*Result {
	if len(commands) == 0 {
		return nil
	}

	results := make([]*Result, len(commands))

	for i, cmd := range commands {
		if ctx.Err() != nil {
			results[i] = &Result{
				Name:     cmd.Name,
				ExitCode: ExitCodeError,
				Error:    fmt.Errorf("context cancelled: %w", ctx.Err()),
			}
			continue
		}

		results[i] = r.Run(ctx, cmd)
	}

	return results
}

func (r *Runner) prepareLogFile(cmd Command) *os.File {
	if cmd.LogFile == "" {
		return nil
	}

	dir := filepath.Dir(cmd.LogFile)
	if err := os.MkdirAll(dir, DefaultDirPerm); err != nil {
		pkg.Logger.Warn("Failed to create log directory",
			zap.String("activity", cmd.Name),
			zap.Error(err))
		return nil
	}

	f, err := os.Create(cmd.LogFile)
	if err != nil {
		pkg.Logger.Warn("Failed to create log file",
			zap.String("activity", cmd.Name),
			zap.Error(err))
		return nil
	}

	return f
}

func (r *Runner) streamOutput(wg *sync.WaitGroup, reader io.Reader, logFile *os.File, activityName, streamName string) {
	defer wg.Done()
	const readChunkSize = 8 * 1024
	buf := make([]byte, readChunkSize)
	lineBuf := make([]byte, 0, ScannerInitBufSize)
	truncated := false
	discarding := false

	flushLine := func() {
		if len(lineBuf) == 0 && !truncated {
			return
		}
		line := cleanLine(string(lineBuf))
		if line != "" {
			if logFile != nil {
				_, _ = fmt.Fprintln(logFile, line)
			}
			pkg.Logger.Debug("Activity output",
				zap.String("activity", activityName),
				zap.String("stream", streamName),
				zap.String("line", line))
		}
		if truncated {
			pkg.Logger.Warn("Activity output line truncated",
				zap.String("activity", activityName),
				zap.String("stream", streamName),
				zap.Int("maxBytes", ScannerMaxBufSize))
		}
		lineBuf = lineBuf[:0]
		truncated = false
		discarding = false
	}

	appendChunk := func(chunk []byte) {
		if len(chunk) == 0 || discarding {
			return
		}
		if len(lineBuf)+len(chunk) <= ScannerMaxBufSize {
			lineBuf = append(lineBuf, chunk...)
			return
		}
		remain := ScannerMaxBufSize - len(lineBuf)
		if remain > 0 {
			lineBuf = append(lineBuf, chunk[:remain]...)
		}
		truncated = true
		discarding = true
	}

	for {
		n, err := reader.Read(buf)
		if n > 0 {
			data := buf[:n]
			for len(data) > 0 {
				i := bytes.IndexAny(data, "\r\n")
				if i == -1 {
					appendChunk(data)
					break
				}

				appendChunk(data[:i])
				delim := data[i]
				data = data[i+1:]
				if delim == '\r' && len(data) > 0 && data[0] == '\n' {
					data = data[1:]
				}
				flushLine()
			}
		}
		if err != nil {
			if err != io.EOF {
				pkg.Logger.Warn("Error reading output stream",
					zap.String("activity", activityName),
					zap.String("stream", streamName),
					zap.Error(err))
			}
			flushLine()
			return
		}
	}
}

// cleanLine removes ANSI escape sequences and control characters from output
func cleanLine(line string) string {
	line = ansiRegex.ReplaceAllString(line, "")
	line = controlCharReplacer.Replace(line)
	return strings.TrimSpace(line)
}

const logSeparator = "============================================================"

func (r *Runner) writeLogHeader(f *os.File, cmd Command) {
	_, _ = fmt.Fprintf(f, "$ %s\n", cmd.Command)
	_, _ = fmt.Fprintln(f, logSeparator)
	_, _ = fmt.Fprintf(f, "# Tool: %s\n", cmd.Name)
	_, _ = fmt.Fprintf(f, "# Started: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	_, _ = fmt.Fprintf(f, "# Timeout: %v\n", cmd.Timeout)
	_, _ = fmt.Fprintln(f, "# Status: Running...")
	_, _ = fmt.Fprintln(f, logSeparator)
	_, _ = fmt.Fprintln(f)
}

func (r *Runner) writeLogFooter(f *os.File, result *Result) {
	status := "✓ Success"
	if result.ExitCode != 0 {
		status = "✗ Failed"
	}

	_, _ = fmt.Fprintln(f)
	_, _ = fmt.Fprintln(f, logSeparator)
	_, _ = fmt.Fprintf(f, "# Finished: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	_, _ = fmt.Fprintf(f, "# Duration: %.2fs\n", result.Duration.Seconds())
	_, _ = fmt.Fprintf(f, "# Exit Code: %d\n", result.ExitCode)
	_, _ = fmt.Fprintf(f, "# Status: %s\n", status)
	_, _ = fmt.Fprintln(f, logSeparator)
}
