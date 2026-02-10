package daemon

import (
"os"
"path/filepath"
"testing"
"github.com/stretchr/testify/assert"
)

func TestDaemonLockUnlock(t *testing.T) {
tmpDir, err := os.MkdirTemp("", "daemon_test")
if err != nil {
t.Fatal(err)
}
defer os.RemoveAll(tmpDir)

pidFile := filepath.Join(tmpDir, "test.pid")
d := NewDaemon(pidFile)

// Test Lock
err = d.Lock()
assert.NoError(t, err)

// Verify file exists and contains PID
data, err := os.ReadFile(pidFile)
assert.NoError(t, err)
assert.NotEmpty(t, string(data))

// Test double Lock (should fail)
err = d.Lock()
assert.Error(t, err)
assert.Contains(t, err.Error(), "already running")

// Test Unlock
err = d.Unlock()
assert.NoError(t, err)

// Verify file is gone
_, err = os.Stat(pidFile)
assert.True(t, os.IsNotExist(err))
}

func TestStalePIDCleanup(t *testing.T) {
tmpDir, err := os.MkdirTemp("", "daemon_stale_test")
if err != nil {
t.Fatal(err)
}
defer os.RemoveAll(tmpDir)

pidFile := filepath.Join(tmpDir, "stale.pid")

// Write a non-existent PID to the file
err = os.WriteFile(pidFile, []byte("999999"), 0644)
assert.NoError(t, err)

d := NewDaemon(pidFile)

// Lock should succeed because PID 999999 is not running
err = d.Lock()
assert.NoError(t, err)

// Verify file now contains current PID
pid, err := d.GetPID()
assert.NoError(t, err)
assert.Equal(t, os.Getpid(), pid)
}
