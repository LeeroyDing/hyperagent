package token

import (
"testing"
"github.com/stretchr/testify/assert"
)

func TestTokenManager_All(t *testing.T) {
tm, err := NewTokenManager("gemini-1.5-flash")
assert.NoError(t, err)

t.Run("CountTokens", func(t *testing.T) {
assert.Equal(t, 0, tm.CountTokens(""))
assert.Equal(t, 2, tm.CountTokens("hello world"))
})

t.Run("PruneHistory", func(t *testing.T) {
messages := []string{"msg1", "msg2", "msg3"}
pruned := tm.PruneHistory(messages, 2)
assert.Equal(t, 1, len(pruned))
assert.Equal(t, "msg3", pruned[0])
})
}
