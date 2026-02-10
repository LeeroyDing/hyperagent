package executor

import (
"testing"

"github.com/stretchr/testify/assert"
)

func TestSessionManager_Basic(t *testing.T) {
m := NewSessionManager()
defer m.Cleanup()

s1, err := m.GetOrCreate("sess1")
assert.NoError(t, err)
assert.NotNil(t, s1)

s2, err := m.GetOrCreate("sess1")
assert.NoError(t, err)
assert.Equal(t, s1, s2)
}
