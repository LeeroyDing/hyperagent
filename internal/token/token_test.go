package token

import (
"testing"
)

func TestTokenManager_CountTokens(t *testing.T) {
tm, err := NewTokenManager("gemini-1.5-flash")
if err != nil {
t.Fatalf("NewTokenManager() error = %v", err)
}

tests := []struct {
name string
text string
want int
}{
{"empty", "", 0},
{"simple", "hello world", 2},
{"sentence", "This is a test of the token counter.", 9},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
if got := tm.CountTokens(tt.text); got != tt.want {
t.Errorf("TokenManager.CountTokens() = %v, want %v", got, tt.want)
}
})
}
}

func TestTokenManager_PruneHistory(t *testing.T) {
tm, _ := NewTokenManager("gemini-1.5-flash")
messages := []string{
"message one",
"message two",
"message three",
}

// Each message is ~2 tokens. Total ~6.
// Prune to 4 tokens -> should keep last 2 messages.
pruned := tm.PruneHistory(messages, 4)
if len(pruned) != 2 {
t.Errorf("PruneHistory() length = %d, want 2", len(pruned))
}
if pruned[0] != "message two" || pruned[1] != "message three" {
t.Errorf("PruneHistory() content = %v, want [message two message three]", pruned)
}
}
