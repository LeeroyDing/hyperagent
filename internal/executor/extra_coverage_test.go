package executor

import (
"context"
"fmt"
"testing"
"time"

"github.com/stretchr/testify/assert"
)

func TestShellSession_ExtraCoverage(t *testing.T) {
s, err := NewShellSession("extra-cov")
if err != nil {
t.Skip("PTY not supported in this environment")
return
}
defer s.Close()

t.Run("Successful Execute", func(t *testing.T) {
got, err := s.Execute("echo hello")
assert.NoError(t, err)
assert.Contains(t, got, "hello")
})

t.Run("Execute on closed session", func(t *testing.T) {
s2, _ := NewShellSession("close-me")
s2.Close()
_, err := s2.Execute("ls")
assert.Error(t, err)
assert.Contains(t, err.Error(), "session closed")
})

t.Run("Double close", func(t *testing.T) {
s2, _ := NewShellSession("double-close")
err := s2.Close()
assert.NoError(t, err)
err = s2.Close()
assert.NoError(t, err)
})

t.Run("Context timeout in executeWithSentinel", func(t *testing.T) {
ctx, cancel := context.WithTimeout(context.Background(), 1*time.Microsecond)
defer cancel()
time.Sleep(2 * time.Millisecond) // Ensure timeout

_, err := s.executeWithSentinel(ctx, "echo slow")
assert.Error(t, err)
assert.Equal(t, context.DeadlineExceeded, err)
})

t.Run("readLoop error path", func(t *testing.T) {
s3, _ := NewShellSession("read-err")
// Closing the PTY file descriptor should trigger an error in readLoop
s3.Pty.Close()
// Wait a bit for the loop to exit
time.Sleep(50 * time.Millisecond)
})
}

func TestSessionManager_ExtraCoverage(t *testing.T) {
m := NewSessionManager()
defer m.Cleanup()

t.Run("GetOrCreate Creator Error", func(t *testing.T) {
m.Creator = func(id string) (Shell, error) {
return nil, fmt.Errorf("fail")
}
_, err := m.GetOrCreate("fail-sess")
assert.Error(t, err)
})

t.Run("GetOrCreate Race/Double Check", func(t *testing.T) {
m.Creator = func(id string) (Shell, error) {
return &MockShell{}, nil
}
s1, err := m.GetOrCreate("race-sess")
assert.NoError(t, err)
assert.NotNil(t, s1)
})
}
