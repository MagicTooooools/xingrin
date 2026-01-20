package activity

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/orbit/worker/internal/pkg"
	"github.com/stretchr/testify/require"
)

func TestRunner_Run_InvalidTimeout(t *testing.T) {
	r := NewRunner(t.TempDir())

	res := r.Run(context.Background(), Command{
		Name:    "test",
		Command: "echo hi",
		Timeout: 0,
	})

	require.Error(t, res.Error)
	require.Equal(t, ExitCodeError, res.ExitCode)
}

func TestRunner_RunParallel_RespectsMaxCmdConcurrency(t *testing.T) {
	require.NoError(t, pkg.InitLogger("error"))
	defer pkg.Sync()

	t.Setenv(EnvMaxCmdConcurrency, "1")

	workDir := t.TempDir()
	r := NewRunner(workDir)

	lockDir := filepath.Join(workDir, "lock")
	cmdStr := fmt.Sprintf("mkdir %q || exit 99; sleep 0.3; rmdir %q", lockDir, lockDir)

	cmds := []Command{
		{Name: "c1", Command: cmdStr, Timeout: 5 * time.Second},
		{Name: "c2", Command: cmdStr, Timeout: 5 * time.Second},
		{Name: "c3", Command: cmdStr, Timeout: 5 * time.Second},
	}

	start := time.Now()
	results := r.RunParallel(context.Background(), cmds)
	elapsed := time.Since(start)

	require.Len(t, results, len(cmds))
	for i, res := range results {
		require.NotNil(t, res, "result %d should not be nil", i)
		require.NoError(t, res.Error)
		require.Equal(t, 0, res.ExitCode)
	}

	// With concurrency=1 and 3 commands sleeping 0.3s each, total wall time should be ~0.9s.
	require.GreaterOrEqual(t, elapsed, 750*time.Millisecond)
}
