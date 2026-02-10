package executor

import (
"context"
"testing"
"time"
"github.com/stretchr/testify/assert"
)

func TestSession_Extra(t *testing.T) {
s, err := NewShellSession()
if err != nil {
t.Skip("PTY not available")
}
defer s.Close()

t.Run("Execute_Timeout", func(t *testing.T) {
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
defer cancel()
_, err := s.Execute(ctx, "sleep 1")
assert.Error(t, err)
})

t.Run("Execute_Empty", func(t *testing.T) {
out, err := s.Execute(context.Background(), "")
assert.NoError(t, err)
assert.Empty(t, out)
})
}
